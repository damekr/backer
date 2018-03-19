// +build linux darwin

package storage

const bucketsLocation string = "/data/"

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

//func GetRepository() *Storage {
//	log.Debug("Getting a storage under: ", Storage.Location)
//	return Storage
//}
//
//func (r *Storage) GetCapacityStatus() (disk DiskStatus) {
//	fs := syscall.Statfs_t{}
//	err := syscall.Statfs(r.Location, &fs)
//	if err != nil {
//		log.Error("Cannot check file system capacity")
//	}
//	disk.All = fs.Blocks * uint64(fs.Bsize)
//	disk.Free = fs.Bfree * uint64(fs.Bsize)
//	disk.Used = disk.All - disk.Free
//	return
//}
//
//func (r *Storage) CreateClientBucket(name string) error {
//	clientBucketLocation := r.Location + bucketsLocation + name
//	log.Debugf("Creating client bucket under: ", clientBucketLocation)
//	err := os.MkdirAll(clientBucketLocation, 0700)
//	if err != nil {
//		log.Debugf("clientDefinition bucket %s exists, skipping", name)
//	}
//	return nil
//}
//
//func (r *Storage) GetClientBucket(name string) (*ClientBucket, error) {
//	clientLocation := filepath.Join(r.Location, bucketsLocation, name)
//	if !checkIfClientBucketExists(name) {
//		log.Errorf("clientDefinition %s bucket does not exist", name)
//		return nil, errors.New("clientDefinition bucket does not exists")
//	}
//	return &ClientBucket{
//		Location: clientLocation,
//	}, nil
//}
//
//func (r *Storage) GetMetadataPath() string {
//	return filepath.Join(r.Location, ".meta/db")
//}
//
//func InitRepository() error {
//	_, err := CreateRepository()
//	if err != nil {
//		log.Println("Cannot create storage, error: ", err.Error())
//		return err
//	}
//	return nil
//}
