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
    Name        string
    Paths       []string
    Size        int64
    Location    string
}

func NewArchive(Paths []string, Name string) *Archive{
    return &Archive{
        Name: Name,
        Paths: Paths,
        Size: 0,
        Location: "",
    }
}



func (a *Archive) validatePaths(){
    validatedPaths := []string{}
    for _,path := range a.Paths{
        log.Printf("Checking path %s", path)
        _, err := os.Stat(path)
        if err != nil{
            log.Printf("Path %s does not exist\n", path)
        } else {
            validatedPaths = append(validatedPaths, path)
        }
    }
    a.Paths = validatedPaths
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

// getFilesAbsPaths checks if given file exists and sets abs path to be included in tar archive
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

// MakeArchive is able to create an archive from given paths
func (a *Archive) MakeArchive(location string){
    a.validatePaths()
    log.Printf("Making tar package from %s \n", a.Paths)
    a.Location = location
    absFilePaths := getFilesAbsPaths(a.Paths)
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
    fi, err := tarfile.Stat()
    if err != nil {
        log.Fatalln(err)       
    }
    a.Size = fi.Size()
    log.Printf("File size: %d Bytes", fi.Size())

}