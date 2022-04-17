package runner

import (
	"errors"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Status int

func (s Status) String() string {
	switch s {
	case Timeout:
		return "timeout"
	case Interrupt:
		return "interrupt"
	default:
		return "normal"
	}
}

var (
	ErrTimeout   = errors.New("timeout")
	ErrInterrupt = errors.New("interrupt")
)

func Status2Err(s Status) error {
	switch s {
	case Timeout:
		return ErrTimeout
	case Interrupt:
		return ErrInterrupt
	default:
		return nil
	}
}

const (
	Normal Status = iota
	Timeout
	Interrupt
)

type Runner[T any] interface {
	launch(opts ...Option) error
	start() error
	do(task Task) error
	Add(tasks ...Task) error
	Status() Status
	Wait()
}

type RunnerImpl[T any] struct {
	fo Option
	s  Status
	wg sync.WaitGroup
}

func (r *RunnerImpl[T]) Wait() {
	r.wg.Wait()
}

func (r *RunnerImpl[T]) Status() Status {
	return r.s
}

func (r *RunnerImpl[T]) launch(opts ...Option) error {
	r.fo = parseOptions(opts...)
	signal.Notify(r.fo.signal(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		r.start()
	}()
	go func() {
		select {
		case <-r.fo.timeouts():
			r.s = Timeout
			log.Println("time out")
		case <-r.fo.signal():
			r.s = Interrupt
			log.Println("signal out")
		}
	}()
	return nil
}

func (r *RunnerImpl[T]) start() error {
	go func(r *RunnerImpl[T]) {
		for task := range r.fo.task() {
			go r.do(task)
		}
	}(r)
	return nil
}

func (r *RunnerImpl[T]) Add(tasks ...Task) error {
	r.wg.Add(len(tasks))
	go func() {
		for _, task := range tasks {
			r.fo.task() <- task
		}
	}()
	return nil
}

func (r *RunnerImpl[T]) do(task Task) error {
	defer r.wg.Done()
	log.Printf("Run %s in\n", task.Id())
	st := time.Now()
	err := task.Run(r)
	log.Printf("Run %s out, spend time: %s, err: %v\n", task.Id(), time.Since(st), err)
	return err
}

func New[T Task](opts ...Option) Runner[T] {
	r := new(RunnerImpl[T])
	r.launch(opts...)
	return r
}
