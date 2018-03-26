package task

// Task defines operations between client and server
type Task interface {
	Run()
	Stop()
}
