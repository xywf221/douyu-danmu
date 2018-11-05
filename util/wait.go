package util

import "sync"

type WaitGroup struct {
	s *sync.WaitGroup
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		s: &sync.WaitGroup{},
	}
}

func (w *WaitGroup) Execute(f func()) {
	w.s.Add(1)
	go func() {
		f()
		w.s.Done()
	}()
}

func (w *WaitGroup) Wait() {
	w.s.Wait()
}
