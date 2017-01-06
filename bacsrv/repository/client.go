package repository

type ClientBucket struct {
	name string
}

func CreateClient(name string) *ClientBucket {
	repo := GetRepository()
	repo.CreateClientBucket(name)
	return &ClientBucket{name: name}
}
