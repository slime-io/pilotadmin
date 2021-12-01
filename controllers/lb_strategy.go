package controllers

import (
	api "slime.io/slime/modules/pilotadmin/api/v1alpha1"
)

type loadBalanceStrategy interface {
	CalLoadBalance(admin *api.PilotAdmin)
	SetLB2Pilot(lb *LoadBalanceAttr) error
}

type LoadBalanceAttr struct {
	ip            string // pilot实例ip
	disconnection int64  // 需要pilot断开的连接数
	keepAliveCon  int64  // 断开后pilot持有的连接数
	keepLimitTime int64  // 触发LB之后的限流保持时间,单位s
}
