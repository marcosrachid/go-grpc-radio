# go-grpc-stream
====

GRPC server and client with golang based on grpc lib to stream audio from server to client

- - - - 

## Necessary Technology Versions

Technology  | Version
------------- | -------------
Go | go1.14.3 linux/amd64
Docker | 18.09.6
docker-compose | 1.24.1
libprotoc | 3.13.0

## Pre-install

Pre-install portaudio and libmpg123 to run

    $ apt-get update && apt-get install alsa-utils avahi-utils portaudio19-dev libmpg123-dev -y

## Pre-running

If you're going to start the application manually

    $ cd ssl && ./generate-keys.sh

## Regenerate protobuf files

If any changes were made to proto files 

    $ cd internal/pb && protoc --go_out=. --go-grpc_out=. stream.proto

## Running

To run the stream server we create a docker container for it

    $ docker-compose up -d

## Configurations

### Server Environment Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| PORT | Server Port | 4000 |

### Client Environment Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| ADDR | Server IP | localhost |
| PORT | Server Port | 4000 |