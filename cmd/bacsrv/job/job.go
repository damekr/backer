package job

import (
	"github.com/damekr/backer/cmd/bacsrv/task"
)

type Job struct {
	Tasks []task.Task
	ID    int
	Name  string
}

var id = 0

var Jobs []*Job

func Create(name string) *Job {
	id++
	newJob := &Job{
		ID:   id,
		Name: name,
	}
	Jobs = append(Jobs, newJob)
	return newJob
}

func (j *Job) AddTask(task task.Task) error {
	j.Tasks = append(j.Tasks, task)
	return nil
}

func (j *Job) Start() error {
	for _, t := range j.Tasks {
		t.Run()
	}
	return nil

}
