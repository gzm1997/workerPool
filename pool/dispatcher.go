package pool

type Dispatcher struct {
	AvailableWorkers []*Worker
	Quit chan bool
}

const MaxNumOfWorkers = 5

func (dis *Dispatcher) Run() {
	// 将限定数量的worker添加到数组中 然后通过start发送到pool里面 然后等待下面从pool中取出worker运行任务
	for i := 0; i < MaxNumOfWorkers; i++ {
		w := NewWorker()
		dis.AvailableWorkers = append(dis.AvailableWorkers, w)
		// 添加到pool中 并且并且等到有没有任务 如果有那么久执行
		w.Start()
	}

	for {
		select {
		// 检测到有任务从任务队列中传送过来
		case task := <- AllTaskQueue:
			// 从pool中获取可用的任务传输通道
			w := <- TaskPool
			//把这个任务通过这个通道发送过去 让另一边接受处理
			w <- task
		case <-dis.Quit:
			return
		}
	}
}


