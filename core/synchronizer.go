package core

import (
	"math/rand"
	"sync"
)

type Synchronizer struct {
	out   *Port
	in    *Port
	queue chan int64
	tasks map[int64]chan pullFunc
	done  map[int64]chan bool
	mutex *sync.Mutex
}

type pushFunc func(port *Port)
type pullFunc func(port *Port)

func (s *Synchronizer) Init(in, out *Port) {
	s.in = in
	s.out = out
	s.queue = make(chan int64)
	s.tasks = make(map[int64]chan pullFunc)
	s.done = make(map[int64]chan bool)
	s.mutex = &sync.Mutex{}
}

func (s *Synchronizer) Push(push pushFunc) int64 {
	var token int64
	token = rand.Int63()

	s.mutex.Lock()
	push(s.out)
	s.tasks[token] = make(chan pullFunc)
	s.done[token] = make(chan bool)
	s.queue <- token // order is important! worker expects to have s.tasks[token] once it can pull the token from queue
	s.mutex.Unlock()

	return token
}

func (s *Synchronizer) Pull(token int64, pull pullFunc) {
	s.tasks[token] <- pull
	<-s.done[token]
	delete(s.tasks, token)
}

func (s *Synchronizer) Worker() {
	for true {
		token := <-s.queue
		pull := <-s.tasks[token]
		pull(s.in)
		s.done[token] <- true
	}
}
