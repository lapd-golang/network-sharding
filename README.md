## Block validation processing before being included in the Blockchain
After having cloned the project, go to project location and install the below mentioned dependencies.

- Install protoc compiler - [protoc](https://github.com/protocolbuffers/protobuf/releases)

- Install grpc Go package - `go get -u google.golang.org/grpc`.

- Install context package - `go get -u golang.org/x/net/context`.

- Install protobuf package - `go get -u github.com/golang/protobuf/proto`

### Steps to start the blockchain server

Open terminal and go to project location

To start the blockchain server execute the following command
```
go run server/main.go
```

### Steps to use RPC client to access blockchain server and it's services

Once the server is started, open another terminal for RPC calls to blockchain server.

Go to project folder and execute the below command to see details about a block with block number 3.
```
go run rpcClient/main.go -block 3
```
To add a new block to blockchain execute the below command
```
go run rpcClient/main.go -opt addBlock
```





