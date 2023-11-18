package erlang

import (
	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/node"
)

type Node struct {
	srv node.Node
}

type NodeConfig struct {
	NodeName string   `yaml:"nodeName" json:"nodeName"`
	Cookie   string   `yaml:"cookie" json:"cookie"`
	Nodes    []string `yaml:"nodes" json:"nodes"`
}

func NewErlangNode(c *NodeConfig) (*Node, error) {
	n, err := ergo.StartNode(c.NodeName, c.Cookie, node.Options{})
	if err != nil {
		return nil, err
	}
	for _, s := range c.Nodes {
		err = n.Connect(s)
		if err != nil {
			return nil, err
		}
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
