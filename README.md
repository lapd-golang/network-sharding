# Block processing before being included in the Blockchain

### Steps to start the blocchain server

After having cloned to project.

Open terminal and go to project location

To start the blockchain server execute the following command

go run server/main.go


### Steps to start RPC client to access blockchain server

Once the server is started, open another terminal for RPC calls to blockchain server.

Go to project folder and execute the below command see details about a block with a block number 3.

go run rpcClient/main.go -block 3

To add a new block to blockchain execute the below command

go run rpcClient/main.go -opt addBlock






