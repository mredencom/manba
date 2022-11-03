package proxy

import (
	"math"

	"github.com/fagongzi/gateway/pkg/pb/metapb"
	"github.com/fagongzi/gateway/pkg/store"
	"github.com/fagongzi/log"
	"github.com/fagongzi/util/format"
)

var (
	eventTypeStatusChanged = store.EvtType(math.MaxInt32)
	eventSrcStatusChanged  = store.EvtSrc(math.MaxInt32)
)

type statusChanged struct {
	meta   metapb.Server
	status metapb.Status
}

func (r *dispatcher) watch() {
	log.Info("router start watch meta data")

	go r.readyToReceiveWatchEvent()
	err := r.store.Watch(r.watchEventC, r.watchStopC)
	log.Errorf("router watch failed, errors:\n%+v",
		err)
}

func (r *dispatcher) readyToReceiveWatchEvent() {
	for {
		evt := <-r.watchEventC
		switch evt.Src {
		case store.EventSrcCluster:
			r.doClusterEvent(evt)
		case store.EventSrcServer:
			r.doServerEvent(evt)
		case store.EventSrcBind:
			r.doBindEvent(evt)
		case store.EventSrcAPI:
			r.doAPIEvent(evt)
		case store.EventSrcRouting:
			r.doRoutingEvent(evt)
		case store.EventSrcProxy:
			r.doProxyEvent(evt)
		case store.EventSrcPlugin:
			r.doPluginEvent(evt)
		case store.EventSrcApplyPlugin:
			r.doApplyPluginEvent(evt)
		case eventSrcStatusChanged:
			r.doStatusChangedEvent(evt)
		default:
			log.Warnf("unknown event <%+v>", evt)
		}
	}
}

func (r *dispatcher) doRoutingEvent(evt *store.Evt) {
	routing, _ := evt.Value.(*metapb.Routing)
	switch evt.Type {
	case store.EventTypeNew:
		r.addRouting(routing)
	case store.EventTypeDelete:
		r.removeRouting(format.MustParseStrUInt64(evt.Key))
	case store.EventTypeUpdate:
		r.updateRouting(routing)
	default:
		log.Warnf("unknown routing event <%+v>", evt)
	}
}

func (r *dispatcher) doProxyEvent(evt *store.Evt) {
	proxy, _ := evt.Value.(*metapb.Proxy)
	switch evt.Type {
	case store.EventTypeNew:
		r.addProxy(proxy)
	case store.EventTypeDelete:
		r.removeProxy(evt.Key)
	default:
		log.Warnf("unknown proxy event <%+v>", evt)

	}
}

func (r *dispatcher) doAPIEvent(evt *store.Evt) {
	api, _ := evt.Value.(*metapb.API)

	switch evt.Type {
	case store.EventTypeNew:
		r.addAPI(api)
	case store.EventTypeDelete:
		r.removeAPI(format.MustParseStrUInt64(evt.Key))
	case store.EventTypeUpdate:
		r.updateAPI(api)
	default:
		log.Warnf("unknown API event <%+v>", evt)
	}
}

func (r *dispatcher) doClusterEvent(evt *store.Evt) {
	cluster, _ := evt.Value.(*metapb.Cluster)

	switch evt.Type {
	case store.EventTypeNew:
		r.addCluster(cluster)
	case store.EventTypeDelete:
		r.removeCluster(format.MustParseStrUInt64(evt.Key))
	case store.EventTypeUpdate:
		r.updateCluster(cluster)
	default:
		log.Warnf("unknown cluster event <%+v>", evt)
	}
}

func (r *dispatcher) doServerEvent(evt *store.Evt) {
	svr, _ := evt.Value.(*metapb.Server)

	switch evt.Type {
	case store.EventTypeNew:
		r.addServer(svr)
	case store.EventTypeDelete:
		r.removeServer(format.MustParseStrUInt64(evt.Key))
	case store.EventTypeUpdate:
		r.updateServer(svr)
	default:
		log.Warnf("unknown server event <%+v>", evt)
	}
}

func (r *dispatcher) doBindEvent(evt *store.Evt) {
	bind, _ := evt.Value.(*metapb.Bind)

	switch evt.Type {
	case store.EventTypeNew:
		r.addBind(bind)
	case store.EventTypeDelete:
		r.removeBind(bind)
	default:
		log.Warnf("unknown bind event <%+v>", evt)
	}
}

func (r *dispatcher) doPluginEvent(evt *store.Evt) {
	value, _ := evt.Value.(*metapb.Plugin)

	switch evt.Type {
	case store.EventTypeNew:
		r.addPlugin(value)
	case store.EventTypeDelete:
		r.removePlugin(format.MustParseStrUInt64(evt.Key))
	case store.EventTypeUpdate:
		r.updatePlugin(value)
	default:
		log.Warnf("unknown plugin event <%+v>", evt)
	}
}

func (r *dispatcher) doApplyPluginEvent(evt *store.Evt) {
	value, _ := evt.Value.(*metapb.AppliedPlugins)

	switch evt.Type {
	case store.EventTypeNew:
		r.updateAppliedPlugin(value)
	case store.EventTypeDelete:
		r.removeAppliedPlugin()
	case store.EventTypeUpdate:
		r.updateAppliedPlugin(value)
	default:
		log.Warnf("unknown applyPlugin event <%+v>", evt)
	}
}

func (r *dispatcher) doStatusChangedEvent(evt *store.Evt) {
	value := evt.Value.(statusChanged)
	oldStatus := r.getServerStatus(value.meta.ID)

	if oldStatus == value.status {
		return
	}

	newValues := r.copyBinds(metapb.Bind{})
	for _, binds := range newValues {
		hasServer := false
		for _, bind := range binds.servers {
			if bind.svrID == value.meta.ID {
				hasServer = true
				bind.status = value.status
			}
		}

		if hasServer {
			binds.actives = append(binds.actives, value.meta)
			newActives := make([]metapb.Server, 0, len(binds.actives))
			for _, active := range binds.actives {
				if active.ID != value.meta.ID || value.status == metapb.Up {
					newActives = append(newActives, active)
				}
			}

			binds.actives = newActives
		}
	}

	r.binds = newValues
	log.Infof("server <%d> changed to %s", value.meta.ID, value.status.String())
}
