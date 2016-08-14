package main

import (
    "fmt"

)


func main(){
    fmt.Println("OK")
    archive := NewArchive([]string{"/tmp"},"Archiwum")
    existingPaths, err := archive.CheckPaths()
    fmt.Println("Val: , err: ", existingPaths, err)
    archive.MakeArchive(existingPaths, "/tmp/tmp.tar")
    // serveInterface()
}