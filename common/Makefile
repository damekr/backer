# Make file to build protobuffer files


baclnt:
	protoc --go_out=plugins=grpc:. protoclnt/baclnt.proto

backer:
	protoc --go_out=plugins=grpc:. proto/backer.proto

all: backer
