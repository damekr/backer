# Prerequisite: GOROOT and GOARCH must be defined

TARG=bacsrv

# build main executable

main:
	go build -ldflags "-X main.commit=$(shell git rev-parse HEAD)" -o $(TARG)

race:
	go build -race -ldflags "-X main.commit=$(shell git rev-parse HEAD)" -o $(TARG)

