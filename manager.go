package conCurrent

import (
	"context"
	"fmt"
	"log"
	"sync"
)

const (
	defaultUpper        = 10000
	WaitPolicy   Policy = iota
	DiscardPolicy
)

var (
	defaultPanicHandler = func(data interface{}) {
		log.Println("捕捉到panic", data)
	}
)

type Policy int

type Manager struct {
	count        int64  //并发数量
	totalCount   uint64 //历史数量
	upperLimit   int64  //并发上限
	upperPolicy  Policy //数量触发上限后的策略
	panicHandler func(interface{})
	waitingList  List

	*sync.Mutex
}

func (m *Manager) AddWork(work Work) {
	m.Lock()

	if m.count >= m.upperLimit { //如果数量已经超过上限
		m.processUpper(work)
		m.Unlock()

	} else {
		m.start()
		m.Unlock()
		log.Println("释放锁")

		go func(m *Manager, work Work) {
			defer func() {
				if rcv := recover(); rcv != nil {
					m.panicHandler(rcv)
					m.finish()
				}

			}()

			work(context.TODO())

			m.finish()

		}(m, work)
	}
}

func (m *Manager) start() {
	m.count++
	m.totalCount++
	m.show()

}

func (m *Manager) finish() {
	defer log.Println("完成")

	m.Lock()
	if work, ok := m.waitingList.Pop().(Work); ok {
		log.Println("从队列中取", m.waitingList.Len())
		m.totalCount++
		m.Unlock()

		defer func() {
			if rcv := recover(); rcv != nil {
				m.panicHandler(rcv)
				m.finish()
			}
		}()
		work(context.TODO())
	} else {
		m.Lock()
		m.count--
		m.show()

		m.Unlock()

	}

}

func NewManger() *Manager {
	return &Manager{
		upperLimit:   defaultUpper,
		upperPolicy:  WaitPolicy,
		panicHandler: defaultPanicHandler,
		waitingList:  &PriorityQueue{},
		Mutex:        &sync.Mutex{},
	}
}

func (m Manager) show() {
	fmt.Printf("%d/%d/%d\n", m.count, m.totalCount, m.waitingList.Len())
}

//processUpper 处理超过上限的并发
func (m *Manager) processUpper(work Work) {
	switch m.upperPolicy {
	case WaitPolicy:
		log.Println("放入队列")
		m.waitingList.Push(&Item{
			value:    work,
			priority: 1,
		})
	case DiscardPolicy:
		log.Println("达到上限")
	}
}
