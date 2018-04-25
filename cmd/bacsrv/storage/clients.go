package storage

//func checkIfClientBucketExists(name string) bool {
//	repolocation := Storage.AssetsLocation
//	log.Debugf("Checking if %s  bucket exists, under mainrepository: %s", name, repolocation)
//	bucketFolder := filepath.Join(repolocation, bucketsLocation, name)
//	log.Debug("Checking file bucket as foleder: ", bucketFolder)
//	repo, err := os.Stat(bucketFolder)
//	if err == nil && repo.IsDir() {
//		// TODO make more validations
//		log.Infof("clientDefinition %s bucket exists", name)
//		return true
//	}
//	return false
//}
//
//func InitClientsBuckets() error {
//	repo := GetRepository()
//	allClients := config.GetAllClients()
//	log.Debug("Number of integrated clients: ", len(allClients))
//	for _, v := range allClients {
//		log.Printf("clientDefinition info: %s", v.NameWithPath)
//		if !checkIfClientBucketExists(v.NameWithPath) {
//			log.Infof("clientDefinition %s bucket does not exist, creating...", v.NameWithPath)
//			err := repo.CreateClientBucket(v.NameWithPath)
//			if err != nil {
//				log.Errorf("Could not create client %s bucket", v.NameWithPath)
//			}
//		}
//	}
//	return nil
//}
//
