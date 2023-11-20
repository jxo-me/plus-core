package erlang

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/lib"
)

type RpcServer struct {
	gen.Server
	res chan interface{}
}

type RpcMessage struct {
	Node string
	Mod  string
	Fun  string
	Args []etf.Term
}

func (tgs *RpcServer) Init(process *gen.ServerProcess, args ...etf.Term) error {
	tgs.res <- nil
	return nil
}

func (tgs *RpcServer) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	tgs.res <- message
	return gen.ServerStatusOK
}

func (tgs *RpcServer) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	return message, gen.ServerStatusOK
}

func (tgs *RpcServer) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch m := message.(type) {
	case RpcMessage:
		lib.Log("[%s] RPC calling: %s:%s:%s", process.NodeName(), m.Node, m.Mod, m.Fun)
		msg := etf.Tuple{
			etf.Atom("call"),
			etf.Atom(m.Mod),
			etf.Atom(m.Fun),
			etf.List(m.Args),
			process.Self(),
		}
		to := gen.ProcessID{Name: "rex", Node: m.Node}
		if v, e := process.Call(to, msg); e != nil {
			lib.Log("[%s] RPC calling: %s:%s:%s,  error:%#v", process.NodeName(), m.Node, m.Mod, m.Fun, e)
			tgs.res <- e
		} else {
			tgs.res <- v
		}
		return gen.ServerStatusOK
	}

	return gen.ServerStatusOK
}

type makeCall struct {
	to      interface{}
	message interface{}
}

type makeCast struct {
	to      interface{}
	message interface{}
}

func (tgs *RpcServer) HandleDirect(process *gen.ServerProcess, ref etf.Ref, message interface{}) (interface{}, gen.DirectStatus) {
	switch m := message.(type) {
	case makeCall:
		return process.Call(m.to, m.message)
	case makeCast:
		return nil, process.Cast(m.to, m.message)
	}
	return nil, lib.ErrUnsupportedRequest
}

func (tgs *RpcServer) Terminate(process *gen.ServerProcess, reason string) {
	tgs.res <- reason
}
