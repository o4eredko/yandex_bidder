package entities

type Task struct {
	Func interface{}
	Args interface{}
}

func NewTask(f interface{}, args interface{}) *Task {
	return &Task{
		Func: f,
		Args: args,
	}
}
