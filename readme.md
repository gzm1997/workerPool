# goroutine与工作池

标签： go 线程


----------
[github地址][1]

----------
## 背景描述 ##
我们在使用go的时候 特别是面对并发的情况下 经常需要使用多线程 goroutine可以用来解决这个问题 一个goroutine解决一个线程的问题 但是我们要知道一个系统的**最大线程数**是有限的 大到这个限制 那么线程数量就不会增加了 更重要的是 线程太多的时候 CPU需要在线程之间频繁地切换 切换过于频繁也会导致CPU的**使用率**下降的 所以我们是有必要限制在这种情况下的**线程数量**


----------
## 方案 ##
网上的一个方案就是设置工作池或者说线程池 我们从工作队列中获取工作之后只可以从这个工作池中获取一个可用的线程 然后执行工作 工作池中我们其实存的是并不是线程 而是用来传输工作的管道 我们通过限制在工作池中的管道的数量 从工作队列中获取的工作只可以被从工作池中获取出来的一条管道传输到另一边 然后在那边就行处理 这样就可以限制了线程的数量(以上思路和下面图片皆来自CSDN博主)

![微信图片_20180822185246.png-34.8kB][2]


那么有个问题是怎么限制上面所说的工作池中任务管道的数量呢？
其实我们可以将**工作池本身就定义为一个管道** 然后将限定数量的工作管道传输进这个工作池大管道中 然后所有的任务从任务队列中取出之后 都需要从工作池大管道中获取一个工作管道 然后通过这个工作管道将工作传输过去 让那边完成这个工作 完成工作之后再将这个工作管道发送回工作池管道 等待下一次被取出



----------
## 实现 ##
顶一个Task接口 里面有个完成这个工作的方法DoTask 所以的任务都需要实现这接口
```go
type Task interface {
	DoTask() error
}
```
具体定义了一种任务的结构体 实现了Task这个接口
```go
type ItTask struct {
	
}

func (t ItTask) DoTask() error {
	fmt.Println("a task is finished")
	return nil
}
```

我们需要定义一个传输工作的管道类型

    // 用来传输任务的管道
    type TaskChan chan Task

像上面所说 我们把工作池也设计为一个管道

    // 用来传输 传输任务的管道(就是上面那个) 的管道
    type PoolChan chan TaskChan

然后我们需要两个全局变量

 - 缓存我们所有有待完成的工作的一个管道
 - 工作池大管道

```go
var (
	AllTaskQueue TaskChan
	TaskPool PoolChan
)

func init() {
	AllTaskQueue = make(TaskChan)
	TaskPool = make(PoolChan)
}
```
上面开始我忘记了再init函数中对这两个公有变量进行初始化了 导致后面一开始使用这两个管道发送和接受信息一致没反应 但是也不报错 很坑爹 我觉得是因为公有变量如果没有被初始化 是不可以用来传输信息的 但是这种情况是不会报错的


然后我们需要将上面用来传输进一步封装 其实不封装也可以 主要是给给我们这个传输任务的管道加一个quit管道 以应道需要中途中断我们的任务

    type Worker struct {
    	Todo TaskChan
    	quit chan bool
    }

给Worker定义一个启动函数 这个函数将我们这个任务管道Todo发送给工作池管道TaskPool 然后等待这个工作管道中发送来新的工作任务 完成这个任务 并且重新将这个任务管道发送回工作池管道中
```go
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
//产生一个新的worker
func NewWorker() *Worker {
	w := &Worker{}
	w.Todo = make(TaskChan)
	w.quit = make(chan bool)
	return w
}
```

此外我们还需要顶一个分发器
分发器中含有两个重要字段 一个数可用的worker指针数组(其实相当于是任务管道数组) 另一个也是一个quit管道 用来接收暂停的信号

    type Dispatcher struct {
    	AvailableWorkers []*Worker
    	Quit chan bool
    }


同时我们也需要声明一个分发器运行函数 这个函数首先需要声明限定数量的worker 然后将这些worker中的任务管道发送到公有变量 那个工作池管道中 然后每个worker中的任务管道等到传送过来的消息 并且进行处理

第二部分需要做的就是监听任务队列AllTaskQueue这个管道看看有没有新的任务被发送过来 一旦检测到 从工作池管道中获取一个任务管道 并且将这个任务从这个工作管道中传输过去

main函数中要做的事情就很简单了 一个就是运行一个分发器 一个是多线程地制造多个任务 然后阻塞地等待任务呗接受和完成
```go
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
```
运行结果
![image_1clgm6h5upf4qto1d8ovfiuue15.png-23.9kB][3]





完整的源代码可以看我[github地址][1]






 


  [1]: https://github.com/gzm1997/workerPool
  [2]: http://static.zybuluo.com/gzm1997/rl41bhe4usmqigg88e7qqy5x/%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20180822185246.png
  [3]: http://static.zybuluo.com/gzm1997/82n1peoq3l5hd7h0oghzf3wi/image_1clgm6h5upf4qto1d8ovfiuue15.png