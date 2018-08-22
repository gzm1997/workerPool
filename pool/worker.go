package pool

import (
	"fmt"
)

type Worker struct {
	Todo TaskChan
	quit chan bool
}

func (w Worker) Start() {
	go func() {
		for {
			TaskPool <- w.Todo
			select {
			case task := <-w.Todo:
				if err := task.DoTask(); err != nil {
					fmt.Println("task fail")
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func NewWorker() *Worker {
	w := &Worker{}
	w.Todo = make(TaskChan)
	w.quit = make(chan bool)
	return w
}