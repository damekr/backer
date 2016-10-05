package repository

import (
	"os"
    log "github.com/Sirupsen/logrus"
	"syscall"
)


func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

type Repository struct {
    Location string
    // Size uint64

}

type DiskStatus struct {
	All  uint64 
	Used uint64 
	Free uint64 
}


func (r *Repository) GetCapacityStatus() (disk DiskStatus) {
    fs := syscall.Statfs_t{}
    err := syscall.Statfs(r.Location, &fs)
    if err != nil {
        log.Error("Cannot check file system capacity")
    }
    disk.All  = fs.Blocks * uint64(fs.Bsize)
    disk.Free = fs.Bfree * uint64(fs.Bsize)
    disk.Used  = disk.All - disk.Free
    return
}

func (r *Repository) CreateClientBucket(name string) error {
    const bucketsLocation string = "/data/"
    err := os.Mkdir(r.Location + bucketsLocation + name, 0700)
    if err != nil {
        log.Error("Cannot create bucket: ", name)
        return err
    }
    return nil
}