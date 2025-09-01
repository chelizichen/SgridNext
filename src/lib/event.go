package lib

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type Event struct {
	T_ID         string     // 事件ID
	T_TYPE       int        // 事件类型
	T_FUNC       func() any // 事件数据
	T_START_TIME int64      // 事件时间
	T_END_TIME   int64      // 事件时间
	T_RESPONSE   any        // 事件结果
}

type Events struct {
	eventCh            chan *Event
	responseCh         chan *Event
	handleResponseFunc func(*Event)
	sentCount          int64
	completedCount     int64
	debug              bool
}

func NewEvents(size int, handleResponseFunc func(*Event)) *Events {
	return newEvents(size, false, handleResponseFunc)
}

func newEvents(size int, debug bool, handleResponseFunc func(*Event)) *Events {
	return &Events{
		eventCh:            make(chan *Event, size),
		responseCh:         make(chan *Event, size/10),
		handleResponseFunc: handleResponseFunc,
		sentCount:          0,
		completedCount:     0,
		debug:              debug,
	}
}

func (s *Events) TransmitEvent(T_FUNC func() any) {
	// 0到100 随机数
	randNum := rand.Intn(100)
	event := &Event{
		T_TYPE: 1,
		T_FUNC: T_FUNC,
		T_ID:   fmt.Sprintf("EVENT|%v|%v", time.Now().UnixNano(), randNum),
	}
	if s.debug {
		atomic.AddInt64(&s.sentCount, 1)
	}
	s.eventCh <- event
}

func (s *Events) ExecScheduler(event *Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ExecScheduler Panic: %v \n", r)
			event.T_RESPONSE = r
			s.responseCh <- event
		}
	}()

	fmt.Printf("ExecScheduler: %v | %v \n", event.T_TYPE, event.T_ID)
	event.T_START_TIME = time.Now().Unix()
	event.T_RESPONSE = event.T_FUNC()
	s.responseCh <- event
	fmt.Printf("ExecScheduler Res: %v \n", event.T_RESPONSE)
	event.T_END_TIME = time.Now().Unix()
	fmt.Printf("ExecScheduler Time: %v \n", event.T_END_TIME-event.T_START_TIME)
}

func (s *Events) HandleResponse(response *Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("HandleResponse Panic: %v \n", r)
		}
	}()

	s.handleResponseFunc(response)
	if s.debug {
		atomic.AddInt64(&s.completedCount, 1)
		// 每1000条打印一次进度
		if atomic.LoadInt64(&s.completedCount)%1000 == 0 {
			fmt.Printf("Progress: %d completed\n", atomic.LoadInt64(&s.completedCount))
		}
	}
}

func (s *Events) Start() {
	for {
		select {
		case event := <-s.eventCh:
			fmt.Printf("Recv Event: %v | %v \n", event.T_TYPE, event.T_ID)
			go s.ExecScheduler(event)
		case response := <-s.responseCh:
			fmt.Printf("Recv Response: %v | %v \n", response.T_TYPE, response.T_ID)
			go s.HandleResponse(response)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
