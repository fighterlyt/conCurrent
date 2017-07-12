package conCurrent

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Manager struct {
	count        int64  //并发数量
	totalCount   uint64 //历史数量
	panicHandler func()
}

func (m *Manager) AddWork(work Work) {
	m.start()
	go func() {
		defer func() {
			if rcv := recover(); rcv != nil {
				m.panicHandler()
			}
		}()
		work(context.TODO())
		m.finish()
	}()
}

func (m *Manager) start() {
	atomic.AddInt64(&m.count, 1)
	atomic.AddUint64(&m.totalCount, 1)
}

func (m *Manager) finish() {
	atomic.AddInt64(&m.count, -1)
}

func NewManger() *Manager {
	return &Manager{}
}

func (m Manager) show() {
	fmt.Printf("%d/%d\n", m.count, m.totalCount)
}
