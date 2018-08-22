package pool

import "fmt"

type Task interface {
	DoTask() error
}


type ItTask struct {
	
}

func (t ItTask) DoTask() error {
	fmt.Println("a task is finished")
	return nil
}