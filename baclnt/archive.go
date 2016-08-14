package main


import (
    "fmt"
	"os"
	"log"
    "archive/tar"
    "io"
	"path/filepath"
)


type Archive struct {
    Name    string
    Paths   []string
    Size    int64
}

func NewArchive(Paths []string, Name string) *Archive{
    return &Archive{
        Name: Name,
        Paths: Paths,
        Size: 0,
    }
}



func (a *Archive) CheckPaths() ([]string, error){
    existingPaths := []string{}
    for _,path := range a.Paths{
        log.Printf("Checking path %s", path)
        _, err := os.Stat(path)
        if err != nil{
            log.Printf("Path %s does not exist\n", path)
        } else {
            existingPaths = append(existingPaths, path)
        }

    }
    return existingPaths, nil
}

func addFile(tw * tar.Writer, location string) error {
    file, err := os.Open(location)
    if err != nil {
        return err
    }
    defer file.Close()
    if stat, err := file.Stat(); err == nil {
        header := new(tar.Header)
        header.Name = location
        header.Size = stat.Size()
        header.Mode = int64(stat.Mode())
        header.ModTime = stat.ModTime()
        if err := tw.WriteHeader(header); err != nil {
            return err
        }
        if _, err := io.Copy(tw, file); err != nil {
            return err
        }
    }
    return nil
}

func getFilesAbsPaths(paths []string) []string{
    fileList := []string{}
    for i := range paths{
        err := filepath.Walk(paths[i], func(path string, f os.FileInfo, err error) error {
            if f.Mode().IsRegular()  {
                 fileList = append(fileList, path)
            }
       
        return nil     
         })
        if err != nil{
            log.Fatalln(err)
        }
    }
    for _, file := range fileList {
        fmt.Println("File: ", file)
    }
   return fileList
}


func (a *Archive) MakeArchive(existingPaths []string, location string){
    log.Printf("Making tar package from %s \n", existingPaths)
    absFilePaths := getFilesAbsPaths(existingPaths)
    tarfile, err := os.Create(location)
    if err != nil{
        log.Panic("The location does not exists")
    }
    defer tarfile.Close() 
    tw := tar.NewWriter(tarfile)
    defer tw.Close()
    for i := range absFilePaths {
        if err := addFile(tw, absFilePaths[i]); err != nil {
            log.Fatalln(err)
        }
    }


}