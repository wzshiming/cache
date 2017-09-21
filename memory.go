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

func (m *Memory) Get(key string) interface{} {
	i, ok := m.m.Load(key)
	if !ok {
		return nil
	}
	n, ok := i.(*Node)
	if !ok {
		return nil
	}
	return n.d
}

func (m *Memory) Put(key string, val interface{}, timeout time.Duration) error {
	err := m.Delete(key)
	if err != nil {
		return err
	}

	if timeout < 0 {
		timeout = 0
	}

	now := time.Now()
	m.m.Store(key, &Node{
		t: now,
		d: val,
		n: m.t.Add(now.Add(timeout), func() {
			m.m.Delete(key)
		}),
	})
	return nil
}

func (m *Memory) Delete(key string) error {
	i, ok := m.m.Load(key)
	if !ok {
		return nil
	}
	n, ok := i.(*Node)
	if !ok {
		m.m.Delete(key)
		return nil
	}
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
