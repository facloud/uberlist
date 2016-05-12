package uberlist

type TaskID uint

type Task struct {
	ID    TaskID
	Title string
}

type TaskList struct {
	Tasks []Task
}
