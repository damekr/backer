package task



type Task interface {
	Start(chan int)
	Stop()
}

