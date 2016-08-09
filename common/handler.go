package common

import (
	"fmt"
)


type Backup struct {
    Paths   string
}

type Args int

type Result int

func (args *Args) ShowPaths(backup Backup, result *Result) error{
    *result = 1
    fmt.Println("Paths: ", backup.Paths)
    return nil
}