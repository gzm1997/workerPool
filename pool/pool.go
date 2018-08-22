package pool

// 用来传输任务的管道
type TaskChan chan Task

// 用来传输 传输任务的管道(就是上面那个) 的管道
type PoolChan chan TaskChan

var (
	AllTaskQueue TaskChan
	TaskPool PoolChan
)

func init() {
	AllTaskQueue = make(TaskChan)
	TaskPool = make(PoolChan)
}