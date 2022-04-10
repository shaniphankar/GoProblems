package Problem2

type Task struct {
	name string
	args []string
	start int64
	delay int64
	index int
}

func NewTask(name string, args []string, start int64, delay int64, index int) *Task {
	t := new(Task)
	t.name=name
	t.args=args
	t.start=start
	t.delay=delay
	t.index=index
	return t
}

func NewTaskWithoutIndex(name string, args []string, start int64, delay int64) *Task {
	t := new(Task)
	t.name=name
	t.args=args
	t.start=start
	t.delay=delay
	return t
}
