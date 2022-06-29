package pool

type Task struct {
	f func() // 一个task里面应该有一个具体的业务， 业务名就叫f
}

// NewTask 创建一个任务
func NewTask(handleFunc func()) *Task {
	return &Task{
		f: handleFunc,
	}
}

//执行业务的方法
func (t *Task) Execute() {
	t.f() // 调用任务中已经绑定好的业务方法
}
