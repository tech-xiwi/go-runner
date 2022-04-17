package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tech-xiwi/go-runner/runner"
)

type SimpleTask struct {
	id    int
	steps int
}

func (s SimpleTask) Id() string {
	return fmt.Sprintf("taskId: %d", s.id)
}

func (s SimpleTask) canRun(r runner.Runner[runner.Task]) bool {
	if r.Status() != runner.Normal {
		return false
	}
	return true
}

func (s SimpleTask) Run(r runner.Runner[runner.Task]) error {
	defer func() {
		log.Printf("%s run out with steps: %d\n", s.Id(), s.steps)
	}()
	if s.canRun(r) {
		time.Sleep(time.Duration(s.id) * time.Second) // 此处模拟用户的业务处理逻辑 1
		s.steps++
	} else {
		log.Printf("%s run 1 meet status:%s quit\n", s.Id(), r.Status())
		return runner.Status2Err(r.Status())
	}

	if s.canRun(r) {
		time.Sleep(time.Duration(s.id) * time.Second) // 此处模拟用户的业务处理逻辑 2
		s.steps++
	} else {
		log.Printf("%s run 2 meet status:%s quit\n", s.Id(), r.Status())
		return runner.Status2Err(r.Status())
	}
	return nil
}

var count int

func createTask() runner.Task {
	count++
	return &SimpleTask{
		id: count,
	}
}

func main() {
	sr := runner.New[SimpleTask](runner.WithTimeout(time.After(5*time.Second)), runner.WithSingle(make(chan os.Signal, 1)), runner.WithTask(make(chan runner.Task, 1)))
	sr.Add(createTask(), createTask(), createTask())
	go func() {
		time.Sleep(2 * time.Second)
		sr.Add(createTask(), createTask(), createTask())
	}()
	sr.Wait()
}
