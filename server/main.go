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

	if !ok {
		return in.GetDsblock(), nil
	}
	blocks = append(blocks, &dsBlock)

	blocksptr = &blocks

	bc = s.DSBlockChain
	remainingBlocks := bc.DSBlocks
	remainingBlocks = append(remainingBlocks, &dsBlock)
	bc.DSBlocks = remainingBlocks

	return in.GetDsblock(), nil
}

//GetDSBlock ...
func (s *Server) GetDSBlock(ctx context.Context, in *ds.GetDSBlockRequest) (*ds.ProtoDSBlock, error) {
	blockNum := in.GetBlocknum()

	bcServer := s.DSBlockChain
	blocks := bcServer.DSBlocks

	for _, bptr := range blocks {
		b := *bptr
		expectedBlockNum := b.Header.Blocknum
		if expectedBlockNum == blockNum {
			return m.MapToProtoDSBlock(&b), nil
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
