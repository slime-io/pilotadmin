package controllers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	k8ssource "sigs.k8s.io/controller-runtime/pkg/source"
	"slime.io/slime/framework/bootstrap"
	"slime.io/slime/framework/util"
	controller2 "slime.io/slime/modules/pilotadmin/controller"
	"slime.io/slime/modules/pilotadmin/source"
	"slime.io/slime/modules/pilotadmin/source/aggregate"
	"slime.io/slime/modules/pilotadmin/source/pilot"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	api "slime.io/slime/modules/pilotadmin/api/v1alpha1"
)

const (
	RequestLimitAdmin    = ":8080/admin/limit/request"
	ConnectionLimitAdmin = ":8080/admin/limit/connection"
	LoadBalanceAdmin     = ":8080/admin/loadbalance"
	HttpSchema           = "http://"
	RequestLimit         = "/request"
	ConnectionLimit      = "/connection"
	ClusterLimit         = "/cluster"
)

var (
	log             = logf.Log.WithName("controller_pilotadmin")
	DebounceAfter   = 35 * time.Second // 最小静默时间，可以设置为小于pa的刷新频率的时间，如25s
	DebounceMax     = 5 * time.Minute  // 最大延迟时间，可以设置为静默时间的几倍
	LBWeightRange   = 0.5
	LBKeepLimitTime = int64(20) // 触发LB之后的限流保持时间,单位s
)

type averageConLoadBalance struct{}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PilotAdmin Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, env *bootstrap.Environment) error {
	return add(mgr, NewReconciler(mgr, env))
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(mgr manager.Manager, env *bootstrap.Environment) *ReconcilePilotAdmin {
	eventChan := make(chan source.Event)
	src := &aggregate.Source{}
	ms, err := pilot.NewPilotEndpointSource(env.K8SClient, eventChan)
	if err != nil {
		return nil
	}
	src.AppendSource(ms)
	podSource := pilot.NewPodSource([]*kubernetes.Clientset{env.K8SClient})
	lbs := &averageConLoadBalance{}
	r := &ReconcilePilotAdmin{
		client:           mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		eventChan:        eventChan,
		source:           src,
		metricInfo:       cmap.New(),
		podSource:        podSource,
		pilotLBEventChan: cmap.New(),
		lbStrategy:       lbs,
		stop:             env.Stop,
	}
	r.source.Start(env.Stop)
	r.WatchSource(env.Stop)
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pilotadmin-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PilotAdmin
	err = c.Watch(&k8ssource.Kind{Type: &api.PilotAdmin{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePilotAdmin implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePilotAdmin{}

// ReconcilePilotAdmin reconciles a PilotAdmin object
type ReconcilePilotAdmin struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme

	metricInfo   cmap.ConcurrentMap
	MetricSource source.Source

	metricInfoLock sync.RWMutex
	eventChan      chan source.Event
	source         source.Source
	podSource      *pilot.PodSource

	pilotLBEventChan cmap.ConcurrentMap
	lbStrategy       loadBalanceStrategy

	stop <-chan struct{}
}

func (r *ReconcilePilotAdmin) Refresh(request reconcile.Request, args map[string]string) (reconcile.Result, error) {
	status := make(map[string]*api.PilotAdminStatus_EndpointStatus)
	for k, v := range args {
		s := r.podSource.Get(v)
		if s != nil {
			s["ip"] = v
			ps := &api.PilotAdminStatus_EndpointStatus{
				Status: s,
			}
			status[k] = ps
		} else {
			log.Error(nil, "status is nil", "pod", k)
		}
	}
	pilotStatus := api.PilotAdminStatus{
		Replicas:  int64(len(status)),
		Endpoints: status,
	}

	// update status
	instance := &api.PilotAdmin{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	oldStatus := instance.Status
	instance.Status = pilotStatus
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		log.Error(err, "更新pilotAdmin status时失败, name:"+request.Name+",namespace:"+request.Namespace)
		return reconcile.Result{}, err
	}
	for k, e := range pilotStatus.Endpoints {
		if _, ok := oldStatus.Endpoints[k]; !ok {
			r.process(instance, k, e)
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcilePilotAdmin) process(instance *api.PilotAdmin, pod string, ep *api.PilotAdminStatus_EndpointStatus) {
	spec := instance.Spec.Limit
	if ep.Status == nil {
		log.Error(fmt.Errorf("status未设置"), "status未设置", "pod", pod)
	}
	if ip, ok := ep.Status["ip"]; !ok {
		log.Error(fmt.Errorf("ip未设置"), "ip未设置", "pod", pod)
	} else {
		qps, err := util.CalculateTemplateString(spec.Qps, ep.Status)
		if err != nil {
			log.Error(err, "计算限流额度时出错", "condition", spec.Qps, "status", ep.Status)
		}
		if i, ok := r.metricInfo.Get(ip + RequestLimit); ok {
			if i != qps {
				requestLimit(ip, qps)
				r.metricInfo.Set(ip, qps)
			}
		} else {
			requestLimit(ip, qps)
			r.metricInfo.Set(ip+RequestLimit, qps)
		}

		connections, err := util.CalculateTemplateString(spec.Connections, ep.Status)
		if err != nil {
			log.Error(err, "计算限流额度时出错", "condition", spec.Connections, "status", ep.Status)
		}
		clusters, err := util.CalculateTemplateString(spec.Clusters, ep.Status)
		if err != nil {
			log.Error(err, "计算限流额度时出错", "condition", spec.Clusters, "status", ep.Status)
		}

		oldConnection, _ := r.metricInfo.Get(ip + ConnectionLimit)
		oldCluster, _ := r.metricInfo.Get(ip + ClusterLimit)
		if oldConnection != connections || oldCluster != clusters {
			connectionLimit(ip, connections, clusters)
			r.metricInfo.Set(ip+ConnectionLimit, connections)
			r.metricInfo.Set(ip+ClusterLimit, clusters)
		}
	}
}

func requestLimit(ip string, request int) {
	log.Info("开始设置xds请求限流", "pod", ip)
	url := fmt.Sprintf(ip+RequestLimitAdmin+"?request=%d", request)
	res, err := http.Get(HttpSchema + url)
	if err != nil {
		log.Error(err, "设置失败", "pod", ip)
	}
	if res == nil || res.StatusCode != 200 {
		log.Error(fmt.Errorf("pilot端错误"), "设置失败", "pod", ip)
	} else {
		log.Info("设置成功")
	}
}

func connectionLimit(ip string, connection int, cluster int) {
	log.Info("开始设置xds链接数限流", "pod", ip)
	url := fmt.Sprintf(ip+ConnectionLimitAdmin+"?maxConnection=%d&maxCluster=%d", connection, cluster)
	res, err := http.Get(HttpSchema + url)
	if err != nil {
		log.Error(err, "设置失败", "pod", ip)
	}
	if res == nil || res.StatusCode != 200 {
		log.Error(fmt.Errorf("pilot端错误"), "设置失败", "pod", ip)
	} else {
		log.Info("设置成功")
	}
}

func (r *ReconcilePilotAdmin) processLoadBalance(ns, name string) {

	var eventChan chan *api.PilotAdmin
	paAddr := ns + "/" + name
	if inter, ok := r.pilotLBEventChan.Get(paAddr); ok {
		eventChan = inter.(chan *api.PilotAdmin)
	} else {
		log.Error(fmt.Errorf("LoadBalance: goroutine 处理LB失败，无可用channel"), "pilot地址", paAddr)
		return
	}

	var timeChan <-chan time.Time
	var lastReceiveTime time.Time
	var firstReceiveTime time.Time
	debouncedEvents := 0
	var admin *api.PilotAdmin
	// 去抖 并 触发负载均衡
	for {
		select {
		case admin = <-eventChan:
			lastReceiveTime = time.Now()
			if debouncedEvents == 0 { // 这个周期第一个事件到达
				firstReceiveTime = time.Now()
				timeChan = time.After(DebounceAfter) // 等待静默时间
			}
			debouncedEvents++
		case <-timeChan:
			lastDuration := time.Since(lastReceiveTime)
			firstDuration := time.Since(firstReceiveTime)
			if lastDuration > DebounceAfter || firstDuration > DebounceMax { // 最近一次事件发生时间距今是否超过了静默时间 或者超过了最大等待时间
				point := admin
				go r.lbStrategy.CalLoadBalance(point)
				debouncedEvents = 0
			} else {
				timeChan = time.After(DebounceAfter - lastDuration)
			}
		case <-r.stop:
			return
		}
	}
}

func (lb *averageConLoadBalance) CalLoadBalance(admin *api.PilotAdmin) {

	paAddr := admin.Namespace + "/" + admin.Name
	log.Info("LoadBalance: LB计算", "PilotAdmin", paAddr)
	var sumCon int64
	var minCon int64 = math.MaxInt64
	for key, val := range admin.Status.Endpoints {
		v, ok := val.Status["connections"]
		if !ok {
			log.Error(fmt.Errorf("LoadBalance: 获取connections不存在"), "podName", key)
			return
		}
		con, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		if con < minCon {
			minCon = con
		}
		sumCon = sumCon + con
	}
	averageCon := sumCon / admin.Status.Replicas
	if averageCon < 2 {
		log.Info("LoadBalance: 平均连接数小于2，不进行LB处理", "AverageConnection", averageCon, "PilotAdmin", paAddr)
		return
	}

	weightPercent := admin.Spec.Loadbalance.Weight
	if (float32(averageCon) * (1 - weightPercent)) > (float32(minCon)) { // averageCon*(1-weightPercent) > minCon
		// 有可能是新扩容实例。做全局lb。只要大于平均连接，就触发lb
		lb.doLoadBalance(admin, averageCon, 0)
	} else {
		lb.doLoadBalance(admin, averageCon, weightPercent)
	}
}

func (lb *averageConLoadBalance) doLoadBalance(admin *api.PilotAdmin, averageCon int64, weightPercent float32) {

	paAddr := admin.Namespace + "/" + admin.Name
	log.Info("LoadBalance: 开始处理负载均衡", "PilotAdmin", paAddr, "AvgCon", averageCon, "weightPercent", weightPercent)
	weightAvg := float32(averageCon) * (1 + weightPercent)
	for key, val := range admin.Status.Endpoints {
		v := val.Status["connections"]
		con, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}

		if (float32(con)) > (weightAvg) { // 对该实例需要触发lb,  con > averageCon*(1+weightPercent)
			disconn := con - averageCon
			if ip, ok := val.Status["ip"]; !ok {
				log.Error(fmt.Errorf("LoadBalance: 获取podIp不存在"), "podName", key)
				return
			} else {
				if disconn > 0 {
					log.Info("LoadBalance: 触发pilot的负载均衡", "podName", key, "podIP", ip, "disconnection", disconn)
					lbAttr := &LoadBalanceAttr{
						ip:            ip,
						disconnection: disconn,
						keepAliveCon:  averageCon,
						keepLimitTime: LBKeepLimitTime,
					}
					if erro := lb.SetLB2Pilot(lbAttr); erro != nil {
						return
					}
				}
			}
		}
	}
	log.Info("LoadBalance: 处理负载均衡完毕", "PilotAdmin", paAddr, "AvgCon", averageCon, "weightPercent", weightPercent)
}

func (lb *averageConLoadBalance) SetLB2Pilot(lbattr *LoadBalanceAttr) error {

	url := fmt.Sprintf(lbattr.ip+LoadBalanceAdmin+"?disconnect=%d&keepcon=%d&keepcontime=%d",
		lbattr.disconnection, lbattr.keepAliveCon, lbattr.keepLimitTime)
	res, err := http.Get(HttpSchema + url)
	if err != nil {
		log.Error(err, "LoadBalance: 负载均衡设置失败", "podIP", lbattr.ip, "disconnection", lbattr.disconnection, "keepCon", lbattr.keepAliveCon, "keepLimitTime", lbattr.keepLimitTime)
		return fmt.Errorf("LoadBalance: 负载均衡设置失败")
	}
	if res.StatusCode == 200 {
		log.Info("LoadBalance: 负载均衡设置成功",
			"podIP", lbattr.ip, "disconnection", lbattr.disconnection,
			"keepCon", lbattr.keepAliveCon, "keepLimitTime", lbattr.keepLimitTime)
		return nil
	} else {
		log.Error(fmt.Errorf("LoadBalance: 负载均衡pilot端错误"), "负载均衡设置失败", "podIP", lbattr.ip, "disconnection", lbattr.disconnection, "keepCon", lbattr.keepAliveCon, "keepLimitTime", lbattr.keepLimitTime)
		return fmt.Errorf("LoadBalance: 负载均衡pilot端错误")
	}
}

func (r *ReconcilePilotAdmin) WatchSource(stop <-chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			case e := <-r.eventChan:
				switch e.EventType {
				case source.Update, source.Add:
					if _, err := r.Refresh(reconcile.Request{NamespacedName: e.Loc}, e.Info); err != nil {
					}
				}
			}
		}
	}()
}

func (r *ReconcilePilotAdmin) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PilotAdmin")

	// Fetch the PilotAdmin instance
	instance := &api.PilotAdmin{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	// 异常分支
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	// 资源删除
	if err != nil && errors.IsNotFound(err) {
		for _, f := range controller2.DeleteHook[controller2.PilotAdmin] {
			if err := f(request, r); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// 资源更新
	for _, f := range controller2.UpdateHook[controller2.PilotAdmin] {
		if err := f(instance, r); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcilePilotAdmin) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.PilotAdmin{}).
		Complete(r)
}

func DoUpdate(i metav1.Object, args ...interface{}) error {
	if len(args) == 0 {
		log.Error(nil, "pilotAdmin doUpdate方法参数不足")
		return nil
	}
	if this, ok := args[0].(*ReconcilePilotAdmin); !ok {
		log.Error(nil, "pilotAdmin doUpdate方法参数不足")
	} else {
		if instance, ok := i.(*api.PilotAdmin); !ok {
			log.Error(nil, "pilotAdmin doUpdate方法第一参数需为自身")
		} else {
			eps := instance.Status.Endpoints
			// 如果状态为空，需获取状态后再进行限流管理

			this.source.WatchAdd(types.NamespacedName{
				Namespace: instance.Namespace,
				Name:      instance.Name,
			})

			// 状态不为空，直接进行限流管理
			for k, ep := range eps {
				this.process(instance, k, ep)
			}

			paAddr := instance.Namespace + "/" + instance.Name
			if ok := this.pilotLBEventChan.SetIfAbsent(paAddr, make(chan *api.PilotAdmin)); ok {
				log.Info("LoadBalance: 启动goroutine处理lb事件", "PilotAdmin instance", paAddr)
				go this.processLoadBalance(instance.Namespace, instance.Name)
			}

			// 判断是否需要lb
			if instance.Status.Replicas < 2 {
				log.Info("LoadBalance: Pilot Replicas数量小于2, 不进行负载均衡处理", "Replicas", instance.Status.Replicas)
			} else if len(instance.Status.Endpoints) < 2 {
				log.Info("LoadBalance: Pilot Endpoints数量小于2, 不进行负载均衡处理", "Endpoints", len(instance.Status.Endpoints))
			} else {
				if inter, ok := this.pilotLBEventChan.Get(paAddr); ok {
					channel := inter.(chan *api.PilotAdmin)
					channel <- instance
				}
			}
		}
	}
	return nil
}

func DoRemove(request reconcile.Request, args ...interface{}) error {
	if len(args) == 0 {
		log.Error(nil, "pilotAdmin DoRemove方法参数不足")
		return nil
	}
	if this, ok := args[0].(*ReconcilePilotAdmin); !ok {
		log.Error(nil, "pilotAdmin DoRemove方法第一参数需为自身")
	} else {
		paAddr := request.Namespace + "/" + request.Name
		this.source.WatchRemove(request.NamespacedName)
		if chann, ok := this.pilotLBEventChan.Get(paAddr); ok {
			close(chann.(chan *api.PilotAdmin))
		}
		this.pilotLBEventChan.Remove(paAddr)
		log.Info("LoadBalance: PilotAdmin资源移除", "PilotAdmin", paAddr)
	}
	return nil
}
