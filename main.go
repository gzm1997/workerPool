package main

import (
	"workerPool/pool"
)

func main() {
	dis := pool.Dispatcher{}
	forever := make(chan bool)
	go func() {
		dis.Run()
	}()
	for i := 0; i < 10; i++ {
		go func() {
			t := pool.ItTask{}
			pool.AllTaskQueue <- t
		}()
	}
	<- forever
}
