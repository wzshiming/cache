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
		t: task.NewTask(8),
	}
}

func (m *Memory) load(key string) *Node {
	i, ok := m.m.Load(key)
	if !ok {
		return nil
	}
	n, ok := i.(*Node)
	if !ok {
		m.m.Delete(key)
		return nil
	}
	return n
}

func (m *Memory) Get(key string) interface{} {
	n := m.load(key)
	if n == nil {
		return nil
	}

	return n.d
}

func (m *Memory) Put(key string, val interface{}, timeout time.Duration) error {
	// 关闭原有值得过期时间
	n := m.load(key)
	if n != nil && n.n != nil {
		m.t.Cancel(n.n)
	}

	now := time.Now()
	n = &Node{
		t: now,
		d: val,
	}

	// 设置值过期时间
	if timeout > 0 {
		n.n = m.t.Add(now.Add(timeout), func() {
			m.Delete(key)
		})
	}

	m.m.Store(key, n)
	return nil
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
