package main

import (
    "fmt"

)


func main(){
    fmt.Println("OK")
    archive := NewArchive([]string{"/tmp"},"Archiwum")
    archive.MakeArchive("/home/damian/tmp.tar")
    // serveInterface()
}