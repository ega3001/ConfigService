package node

import (
	"path"

	"github.com/lrita/cmap"
)

// concrete types
type Node struct {
	c      CfgClient
	key    string
	value  cmap.Map[string, any]
	childs map[string]*Node
}

func NewNode(c CfgClient, key string) *Node {
	return &Node{c, key, cmap.Map[string, any]{}, map[string]*Node{}}
}

func InitNodes(c CfgClient) (*Node, error) {
	root := NewNode(c, "/vbox")

	// System
	media := root.CreateChild("system").CreateChild("media")

	// Modules
	lprModule := root.CreateChild("modules").CreateChild("LPR")
	lprModule.CreateChild("lists")

	// Groups of cameras
	root.CreateChild("groups")

	if err := root.load(); err != nil {
		return nil, err
	}
	if err := lprModule.Put(map[string]any{"name": "LPR"}); err != nil {
		return nil, err
	}
	if err := media.Put(map[string]any{"recordTTL": 1, "eventTTL": 1}); err != nil {
		return nil, err
	}
	return root, nil
}

func (n *Node) SetKey(key string) {
	n.key = key
}

func (n *Node) GetKey() string {
	return n.key
}

func (n *Node) ChildsAmount() int {
	return len(n.childs)
}

func (n *Node) ListChildKeys() []string {
	keys := make([]string, 0, len(n.childs))
	for k := range n.childs {
		keys = append(keys, k)
	}

	return keys
}

func (n *Node) CreateChild(key string) *Node {
	child := NewNode(n.c, path.Join(n.key, key))
	n.childs[key] = child
	return child
}

func (n *Node) RemoveChild(key string) error {
	if chN, ok := n.childs[key]; ok {
		if err := chN.c.Delete(chN.key); err != nil {
			return err
		}
		delete(n.childs, key)
		return nil
	}
	return ErrNodeNotExists
}

func (n *Node) GetChild(key string) (*Node, error) {
	if n.childs[key] != nil {
		return n.childs[key], nil
	}
	return nil, ErrNodeNotExists
}

func (n *Node) GetChilds() []*Node {
	result := make([]*Node, 0, len(n.childs))
	for _, child := range n.childs {
		result = append(result, child)
	}
	return result
}

func (n *Node) Get() map[string]any {
	m := make(map[string]any)
	n.value.Range(func(key string, value any) bool {
		m[key] = value
		return true
	})
	return m
}

func (n *Node) Put(data map[string]any) error {
	if err := n.c.Set(n.key, data); err != nil {
		return err
	}
	for key, value := range data {
		n.value.Store(key, value)
	}
	return nil
}

func (n *Node) load() error {
	if err := n.Init(); err != nil {
		return err
	}

	childKeys, err := n.c.Children(n.key)
	if err != nil {
		return err
	}
	for _, chKey := range childKeys {
		if _, ok := n.childs[chKey]; !ok {
			n.CreateChild(chKey)
		}
	}
	for _, chN := range n.childs {
		if err = chN.load(); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) Init() error {
	exists, err := n.c.Exists(n.key)
	if err != nil {
		return err
	}
	if !exists {
		return n.c.Create(n.key, n.Get())
	}
	go func() {
		n.c.WatchNodeChanges(n.key, func(m map[string]any, err error) {
			if err != nil {
				return
			}
			for key, value := range m {
				n.value.Store(key, value)
			}
		})
	}()
	return nil
}
