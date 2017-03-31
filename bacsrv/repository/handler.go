package repository

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

type Repository struct {
	Location string
	// Size uint64
}

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func GetRepository() *Repository {
	log.Debugf("Getting a repository under: ", MainRepository.Location)
	return MainRepository
}

func (r *Repository) GetCapacityStatus() (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(r.Location, &fs)
	if err != nil {
		log.Error("Cannot check file system capacity")
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func (r *Repository) CreateClientBucket(name string) error {
	const bucketsLocation string = "/data/"
	clientBucketLocation := r.Location + bucketsLocation + name
	log.Debugf("Creating client bucket under: ", clientBucketLocation)
	err := os.MkdirAll(clientBucketLocation, 0700)
	if err != nil {
		log.Debugf("Client bucket %s exists, skipping", name)
	}
	return nil
}

func InitRepository() error {
	_, err := CreateRepository()
	if err != nil {
		log.Println("Cannot create repository, error: ", err.Error())
		return err
	}
	return nil
}
