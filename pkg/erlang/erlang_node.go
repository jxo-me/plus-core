package erlang

import (
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	"github.com/gogf/gf/v2/container/gmap"
)

const (
	DefaultProcessName = "ErlangRPC"
	DefaultCallTimeout = 5
)

type Node struct {
	srv     node.Node
	process *gmap.StrAnyMap
}

type NodeConfig struct {
	NodeName   string   `yaml:"nodeName" json:"nodeName"`
	Cookie     string   `yaml:"cookie" json:"cookie"`
	ServerName string   `yaml:"serverName" json:"serverName"`
	Nodes      []string `yaml:"nodes" json:"nodes"`
}

func NewErlangNode(c *NodeConfig) (*Node, error) {
	n, err := ergo.StartNode(c.NodeName, c.Cookie, node.Options{})
	if err != nil {
		return nil, err
	}

	return &Node{srv: n}, nil
}

func (s *Node) Start() error {
	s.srv.Wait()
	return nil
}

func (s *Node) Get() node.Node {
	return s.srv
}

func (s *Node) Run() {
	s.srv.Wait()
}

func (s *Node) Stop() error {
	s.srv.Stop()
	return nil
}

func (s *Node) Spawn(name string, opts gen.ProcessOptions, object gen.ProcessBehavior, args ...etf.Term) (gen.Process, error) {
	var err error
	var p gen.Process
	v := s.process.GetOrSetFuncLock(name, func() interface{} {
		p, err = s.srv.Spawn(name, opts, object, args...)
		if err != nil {
			return nil
		}
		return p
	})
	if v != nil {
		if p, ok := v.(gen.Process); ok {
			if !p.IsAlive() {
				p, err = s.srv.Spawn(name, opts, object, args...)
				if err != nil {
					return nil, err
				}
				s.process.Set(name, p)
			}
			return p, nil
		}
	}

	return nil, err
}

func (s *Node) Call(message *RpcMessage) (etf.Term, error) {
	return s.CallWithTimeout(message, DefaultCallTimeout)
}

func (s *Node) CallWithTimeout(req *RpcMessage, timeout int) (etf.Term, error) {
	reqSrv := &RpcServer{
		res: make(chan interface{}, 2),
	}
	p1, err := s.Spawn(DefaultProcessName, gen.ProcessOptions{}, reqSrv)
	if err != nil {
		return nil, err
	}

	to := gen.ProcessID{Name: "rex", Node: req.Node}
	ref := p1.MakeRef()
	from := etf.Tuple{p1.Self(), ref}
	message := etf.Tuple{
		etf.Atom("call"),
		etf.Atom(req.Mod),
		etf.Atom(req.Fun),
		etf.List(req.Args),
		p1.Self(),
	}
	msg := etf.Term(etf.Tuple{etf.Atom("$gen_call"), from, message})
	err = p1.PutSyncRequest(ref)
	if err != nil {
		return nil, err
	}
	if err = p1.Send(to, msg); err != nil {
		return nil, err
	}
	value, err := p1.WaitSyncReply(ref, timeout)

	return value, err
}
