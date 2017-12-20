package cache

import (
	"sync"
	"time"

	"github.com/wzshiming/task"
)

type Node struct {
	n *task.Node
	d interface{}
	t time.Time
}

type Memory struct {
	m sync.Map
	t *task.Task
}

var _ Cache = (*Memory)(nil)

func NewMemory() *Memory {
	return &Memory{
		m: sync.Map{},
		t: task.NewTask(1),
	}
}

func (m *Memory) node(key string, i interface{}) *Node {
	n, ok := i.(*Node)
	if !ok {
		m.m.Delete(key)
		return nil
	}
	return n
}

func (m *Memory) load(key string) *Node {
	i, ok := m.m.Load(key)
	if !ok {
		return nil
	}
	return m.node(key, i)
}

func (m *Memory) GetOrPut(key string, val interface{}, timeout time.Duration) (interface{}, bool) {
	i, ok := m.m.LoadOrStore(key, &Node{
		t: time.Now(),
		d: val,
	})
	n := m.node(key, i)
	if ok {
		return n.d, true
	}
	m.SetTimeout(key, timeout)
	return n.d, false
}

func (m *Memory) Get(key string) interface{} {
	n := m.load(key)
	if n == nil {
		return nil
	}

	return n.d
}

func (m *Memory) SetTimeout(key string, timeout time.Duration) error {
	if timeout <= 0 {
		return nil
	}
	n := m.load(key)
	if n == nil {
		return nil
	}

	if n.n != nil {
		defer m.t.Cancel(n.n)
	}

	n.n = m.t.Add(n.t.Add(timeout), func() {
		m.Delete(key)
	})
	return nil
}

func (m *Memory) Put(key string, val interface{}, timeout time.Duration) error {
	// 关闭原有值得过期时间
	n := m.load(key)
	if n != nil && n.n != nil {
		defer m.t.Cancel(n.n)
	}

	m.m.Store(key, &Node{
		t: time.Now(),
		d: val,
	})

	return m.SetTimeout(key, timeout)
}

func (m *Memory) Delete(key string) error {
	n := m.load(key)
	if n == nil {
		return nil
	}

	// 关闭过期清除
	if n.n != nil {
		m.t.Cancel(n.n)
	}

	m.m.Delete(key)
	return nil
}

func (m *Memory) IsExist(key string) bool {
	_, ok := m.m.Load(key)
	return ok
}
