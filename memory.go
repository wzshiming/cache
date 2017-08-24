package cache

import (
	"sync"
	"time"

	"github.com/wzshiming/task"
)

type Memory struct {
	m sync.Map
	t *task.Task
}

func NewMemory() *Memory {
	return &Memory{
		m: sync.Map{},
		t: task.NewTask(1),
	}
}

func (m *Memory) Get(key string) interface{} {
	i, _ := m.m.Load(key)
	return i
}

func (m *Memory) Put(key string, val interface{}, timeout time.Duration) error {
	m.m.Store(key, val)
	m.t.Add(time.Now().Add(timeout), func() {
		m.m.Delete(key)
	})
	return nil
}

func (m *Memory) Delete(key string) error {
	m.m.Delete(key)
	return nil
}

func (m *Memory) IsExist(key string) bool {
	_, ok := m.m.Load(key)
	return ok
}
