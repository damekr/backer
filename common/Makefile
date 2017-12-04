# Make file to build protobuffer files

all: baclnt bacsrv


baclnt:
	protoc --go_out=plugins=grpc:. protoclnt/baclnt.proto

bacsrv:
	protoc --go_out=plugins=grpc:. protosrv/bacsrv.proto

