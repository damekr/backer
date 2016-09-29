package transfer

import (
    "os"
    log "github.com/Sirupsen/logrus"
	"path/filepath"
)


func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}


func GetAbsolutePaths(paths []string) []string{
    log.Debug("Checking absolutive paths for: ", paths)
    fileList := []string{}
    for i := range paths{
        err := filepath.Walk(paths[i], func(path string, f os.FileInfo, err error) error {
            if f.Mode().IsRegular()  {
                log.Debugf("Adding file %s to list", path)
                fileList = append(fileList, path)
            } else {
                log.Debug("Found not regular file: ", path)
            }
        return nil     
         })
        if err != nil{
            log.Error(err)
        }
    }
   return fileList
}