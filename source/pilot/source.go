/*
* @Author: yangdihang
* @Date: 2020/10/13
 */

package pilot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cmap "github.com/orcaman/concurrent-map"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"slime.io/slime/modules/pilotadmin/source"
	"slime.io/slime/modules/pilotadmin/source/k8s"
)

const PilotStatusUrl = ":8080/admin/status"

var log = logf.Log.WithName("source_pilot_source")

type PodStatus struct {
	Qps         float64 `json:"qps"`
	Connections int64   `json:"connections"`
	Cluster     int64   `json:"cluster"`
}

type PodSource struct {
	k8sClient []*kubernetes.Clientset
}

func NewPodSource(k8sClient []*kubernetes.Clientset) *PodSource {
	return &PodSource{k8sClient: k8sClient}
}

func (s *PodSource) Get(podIP string) map[string]string {
	ret := make(map[string]string)
	if podIP != "" {
		res, err := http.Get("http://" + podIP + PilotStatusUrl)
		if err != nil || res.StatusCode != 200 {
			log.Error(err, "获取pilot status信息失败", "pod", podIP)
			return nil
		}
		body, err := ioutil.ReadAll(res.Body)
		ps := &PodStatus{}
		err = json.Unmarshal(body, ps)
		if err != nil {
			log.Error(err, "解析pilot status信息失败", "pod", podIP)
			return nil
		}
		ret["qps"] = fmt.Sprintf("%d", int(ps.Qps))
		ret["connections"] = fmt.Sprintf("%d", ps.Connections)
		ret["clusters"] = fmt.Sprintf("%d", ps.Cluster)
		return ret
	}
	return nil
}

func NewPilotEndpointSource(c *kubernetes.Clientset, eventChan chan source.Event) (*k8s.Source, error) {
	watcher, err := c.CoreV1().Endpoints("").Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	es := &k8s.Source{
		EventChan:  eventChan,
		Watcher:    watcher,
		K8sClient:  []*kubernetes.Clientset{c},
		UpdateChan: make(chan types.NamespacedName),
		Interest:   cmap.New(),
	}
	es.SetHandler(epGetHandler, epWatcherHandler, epTimerHandler, updateHandler)
	return es, nil
}

// ep handler
func updateHandler(endpointSource *k8s.Source, loc types.NamespacedName) {
	material := endpointSource.Get(loc)
	if len(material) == 0 {
		return
	}
	endpointSource.EventChan <- source.Event{
		EventType: source.Add,
		Loc:       loc,
		Info:      material,
	}
}

func epWatcherHandler(m *k8s.Source, e watch.Event) {
	if e.Object == nil {
		return
	}
	ep, ok := e.Object.(*v1.Endpoints)
	if ok {
		loc := types.NamespacedName{
			Namespace: ep.Namespace,
			Name:      ep.Name,
		}
		if _, exist := m.Interest.Get(loc.Namespace + "/" + loc.Name); !exist {
			return
		}
		updateHandler(m, loc)
	}
}

func epTimerHandler(m *k8s.Source) {
	m.RLock()
	for k := range m.Interest.Items() {
		if index := strings.Index(k, "/"); index == -1 || index == len(k)-1 {
			continue
		} else {
			ns := k[:index]
			name := k[index+1:]
			updateHandler(m, types.NamespacedName{
				Namespace: ns,
				Name:      name,
			})
		}
	}
	m.RUnlock()
}

func epGetHandler(m *k8s.Source, meta types.NamespacedName) map[string]string {
	ret := make(map[string]string)
	for _, client := range m.K8sClient {
		eps, err := client.CoreV1().Endpoints(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if err != nil {
			log.Error(err, "获取endpoints时出错")
			return nil
		}
		for _, subset := range eps.Subsets {
			for _, ep := range subset.Addresses {
				if ep.TargetRef.Kind == "Pod" {
					ret[ep.TargetRef.Name] = ep.IP
				}
			}
		}
	}
	return ret
}
