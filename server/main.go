package main

import (
	"log"
	"net"

	ds "../directoryServiceProto"
	m "./data"
	blockchain "./data/blockChainData"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to listen port 8080: %v", err)
	}
	srv := grpc.NewServer()
	ds.RegisterDSBlockchainServer(srv, &Server{
		DSBlockChain: blockchain.NewBlockchain(),
	})

	//Signs := []*c.Signature{}
	//fmt.Printf("size of pointers: %v\n", c.Signatures.)
	srv.Serve(listener)
}

//Server PutDSBlock(context.Context, *NodeDSBlock) (*ProtoDSBlock, error)
type Server struct {
	DSBlockChain *blockchain.DSBlockChain
}

//PutDSBlock ...
func (s *Server) PutDSBlock(ctx context.Context, in *ds.NodeDSBlock) (*ds.ProtoDSBlock, error) {
	bc := s.DSBlockChain
	blocksptr := &bc.DSBlocks
	blocks := *blocksptr

	dsBlock, ok := blockchain.ProcessNodeSBlock(in, blocks)
	log.Printf("DS Block Detail: %v\n", dsBlock)
	if !ok {
		return in.GetDsblock(), nil
	}
	blocks = append(blocks, &dsBlock)
	log.Printf("size of blockchain in put is now::: %v\n\n", len(blocks))
	blocksptr = &blocks

	bc = s.DSBlockChain
	remainingBlocks := bc.DSBlocks
	remainingBlocks = append(remainingBlocks, &dsBlock)
	bc.DSBlocks = remainingBlocks
	log.Printf("Again size of blockchain in put is now::: %v\n\n", len(bc.DSBlocks))
	return in.GetDsblock(), nil
}

//GetDSBlock ...
func (s *Server) GetDSBlock(ctx context.Context, in *ds.GetDSBlockRequest) (*ds.ProtoDSBlock, error) {
	blockNum := in.GetBlocknum()
	log.Printf("Block number enquired is: %v\n", blockNum)
	blocks := []*ds.ProtoDSBlock{}

	blocks = m.MapToProtoBuffer(s.DSBlockChain.DSBlocks, blocks)
	log.Printf("size of blockchain in get is now::: %v\n\n", len(blocks))
	for _, b := range blocks {
		log.Println("here it is -------------")

		expectedBlockNum := b.Header.Blocknum
		log.Println("then here it is -------------")
		if expectedBlockNum == blockNum {
			return b, nil
		}
	}
	return new(ds.ProtoDSBlock), nil
}

//GetBlockchain ...
func (s *Server) GetBlockchain(context.Context, *ds.GetDSBlockchainRequest) (*ds.GetDSBlockchainResponse, error) {

	blocks := []*ds.ProtoDSBlock{}

	blocks = m.MapToProtoBuffer(s.DSBlockChain.DSBlocks, blocks)

	resp := ds.GetDSBlockchainResponse{
		Blocks: blocks,
	}

	return &resp, nil
}
