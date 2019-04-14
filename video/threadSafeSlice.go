package video

import "sync"

type threadSafeSlice struct {
	sync.Mutex
	workers []*worker
}

func (s *threadSafeSlice) len() int {
	s.Lock()
	defer s.Unlock()

	return len(s.workers)
}

func (s *threadSafeSlice) push(w *worker) {
	s.Lock()
	defer s.Unlock()

	s.workers = append(s.workers, w)
}

func (s *threadSafeSlice) iter(routine func(*worker) bool) {
	s.Lock()
	defer s.Unlock()

	for i := len(s.workers) - 1; i >= 0; i-- {
		remove := routine(s.workers[i])
		if remove {
			s.workers[i] = nil
			s.workers = append(s.workers[:i], s.workers[i+1:]...)
		}
	}
}
