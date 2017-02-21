# Prerequisite: GOROOT and GOARCH must be defined

TARG=baclnt

# build main executable

main:
	go build -ldflags "-X main.commit=$(shell git rev-parse HEAD)" -o $(TARG) 
