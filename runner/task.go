package runner

type Task interface {
	Run(r Runner[Task]) error
	Id() string
}
