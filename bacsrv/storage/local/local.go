package local



type Local struct {
	Location string
}


func Create(location string) (*Local, error) {
	//repolocation := config.GetMainRepositoryLocation()
	//if checkIfRepoExists(repolocation) {
	//	log.Infof("Storage %s exists, skipping creating", repolocation)
	//	Storage = &Storage{Location: repolocation}
	//	return Storage, nil
	//}
	//err := os.MkdirAll(repolocation+"/.meta/init", 0700)
	//if err != nil {
	//	log.Errorf("Cannot create storage %s...", repolocation)
	//	return nil, err
	//}
	//errd := os.MkdirAll(repolocation+"/data", 0700)
	//if errd != nil {
	//	log.Error("Cannot create data directory inside storage")
	//	return nil, errd
	//}
	//
	//erri := os.MkdirAll(repolocation+"/locks", 0700)
	//if erri != nil {
	//	log.Error("Cannot create locks directory inside storage")
	//	return nil, erri
	//}
	//errdb := os.MkdirAll(repolocation+"/.meta/db", 0700)
	//if erri != nil {
	//	log.Error("Cannot create dbs directory inside storage")
	//	return nil, errdb
	//}
	//log.Infof("Storage %s has been created successfully", repolocation)
	//Storage = &Storage{Location: repolocation}
	local := &Local{}
	return local, nil
}

func (l Local) SaveFile(){

}

func (l Local) RemoveFile(){

}
