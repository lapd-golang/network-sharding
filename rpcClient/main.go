package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"time"

	c "../server/crypto"

	dsProto "../directoryServiceProto"
	"google.golang.org/grpc"
)

var client dsProto.DSBlockchainClient

func main() {
	var operation string
	var blocknum string
	flag.StringVar(&operation, "opt", "", "Usage")
	flag.StringVar(&blocknum, "block", "", "Usage")
	flag.Parse()
	getBlockchainFlag := flag.Bool("list", false, "List all blocks")
	flag.Parse()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	client = dsProto.NewDSBlockchainClient(conn)

	if operation == "addBlock" {
		putDSBlock()
	}

	if len(blocknum) > 0 {
		getDSBlock()
	}

	if *getBlockchainFlag {
		getBlockchain()
	}
}

func putDSBlock() {
	//DSBlockChain
	blockchain, getErr := client.GetBlockchain(context.Background(), &dsProto.GetDSBlockchainRequest{})
	if getErr != nil {
		log.Fatalf("Unable to get directory service blockchain: %v", getErr)
	}

	noofblocks := len(blockchain.Blocks)
	if noofblocks == 0 {
		log.Printf("Genesis block is not yet added to blockchain..")
		return
	}
	lastBlock := *blockchain.Blocks[noofblocks-1]
	//Get PoW1 Winner Node
	powDSwinners := getPowDSWinners()

	//get Sharding Structure with 3 members in it
	protoShardStruct := getProtoShardingStructure()

	//DS Committee will have 3 DS Nodes
	dsCommittee := getDSCommittee()

	// Not setting any value for transaction sharing
	txSharingAssign := dsProto.ProtoTxSharingAssignments{
		Dsnodes:    nil,
		Shardnodes: nil,
	}
	//Set DS Block Hashset

	shardinghash := c.Hash(protoShardStruct.String())

	fmt.Printf("1st sharding struc hash: %v \n", shardinghash)
	txsharinghash := c.Hash(txSharingAssign.String())
	dsBlockHashSet := dsProto.ProtoDSBlock_DSBlockHashSet{
		Shardinghash:  shardinghash,
		Txsharinghash: txsharinghash,
		Reservedfield: []byte{0},
	}

	latestblockno := lastBlock.GetHeader().GetBlocknum() + 1
	latestepochno := lastBlock.GetHeader().GetEpochnum() + 1

	header := dsProto.ProtoDSBlock_DSBlockHeader{
		Dsdifficulty:  5,
		Difficulty:    2,
		Prevhash:      lastBlock.GetBlockbase().GetBlockhash(),
		Leaderpubkey:  &dsProto.ByteArray{Data: []byte{1, 2, 3}},
		Blocknum:      latestblockno,
		Epochnum:      latestepochno,
		Gasprice:      &dsProto.ByteArray{Data: []byte{0}},
		Swinfo:        &dsProto.ByteArray{Data: []byte{0}},
		Dswinners:     powDSwinners,
		Hash:          &dsBlockHashSet,
		Committeehash: c.Hash(dsCommittee.String()),
	}

	currentTime := nowAsUnixMilli()

	//Get SHA 256 value of the Header
	hash := c.Hash(header.String())

	//Get co-signatures
	cosigs := getCoSignatures(&dsCommittee)

	protoBlockBase := dsProto.ProtoBlockBase{
		Blockhash: hash,
		Cosigs:    cosigs,
		Timestamp: uint64(currentTime),
	}

	dsblock := dsProto.ProtoDSBlock{
		Header:    &header,
		Blockbase: &protoBlockBase,
	}

	nodeDSBlock := dsProto.NodeDSBlock{
		Shardid:     1,
		Dsblock:     &dsblock,
		Sharding:    protoShardStruct,
		Assignments: &txSharingAssign,
	}

	dsBlock, putErr := client.PutDSBlock(context.Background(), &nodeDSBlock)

	if putErr != nil {
		log.Fatalf("Unable to Put DS Block from the Node: %v", putErr)
	}
	//fmt.Printf("%x \n", *dsBlock)
	log.Printf("New Block Added: %s\n", *dsBlock)
}

func getDSBlock() {
	getDSBlockRequest := dsProto.GetDSBlockRequest{
		Blocknum: 1,
	}
	dsBlock, getErr := client.GetDSBlock(context.Background(), &getDSBlockRequest)

	if getErr != nil {
		log.Fatalf("unable to get directory service blockchain: %v", getErr)
	}

	log.Println("DS Block Details")
	log.Println("--------------------------------")
	log.Printf("%v \n", dsBlock)
	log.Println("--------------------------------")
}

func getBlockchain() {
	blockchain, getErr := client.GetBlockchain(context.Background(), &dsProto.GetDSBlockchainRequest{})

	if getErr != nil {
		log.Fatalf("unable to get directory service blockchain: %v", getErr)
	}

	log.Println("blocks:")

	for _, b := range blockchain.Blocks {
		log.Printf("%v \n", b)
	}
}

func getDSCommittee() dsProto.ProtoDSCommittee {
	sign1 := c.CreateSignature()
	sign2 := c.CreateSignature()
	sign3 := c.CreateSignature()

	//fmt.Printf("DS Comm Pub key: %v\n", sign1.R)

	p1 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 3}},
		Listenporthost: 8080,
	}

	p2 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 4}},
		Listenporthost: 8080,
	}

	p3 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 5}},
		Listenporthost: 8080,
	}

	node1 := dsProto.ProtoDSNode{
		Pubkey: &dsProto.ByteArray{Data: c.EncodePubKey(sign1)},
		Peer:   &dsProto.ByteArray{Data: []byte(p1.String())},
	}

	node2 := dsProto.ProtoDSNode{
		Pubkey: &dsProto.ByteArray{Data: c.EncodePubKey(sign2)},
		Peer:   &dsProto.ByteArray{Data: []byte(p2.String())},
	}

	node3 := dsProto.ProtoDSNode{
		Pubkey: &dsProto.ByteArray{Data: c.EncodePubKey(sign3)},
		Peer:   &dsProto.ByteArray{Data: []byte(p3.String())},
	}

	dsnodes := []*dsProto.ProtoDSNode{}

	dsnodes = append(dsnodes, &node1)
	dsnodes = append(dsnodes, &node2)
	dsnodes = append(dsnodes, &node3)

	dsCommittee := dsProto.ProtoDSCommittee{
		Dsnodes: dsnodes,
	}

	return dsCommittee
}

func getProtoShardingStructure() *dsProto.ProtoShardingStructure {
	//Number of members in a shard is 3
	sign1 := c.CreateSignature()
	sign2 := c.CreateSignature()
	sign3 := c.CreateSignature()

	p1 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 3}},
		Listenporthost: 8080,
	}

	p2 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 4}},
		Listenporthost: 8080,
	}

	p3 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 5}},
		Listenporthost: 8080,
	}

	member1 := dsProto.ProtoShardingStructure_Member{
		Pubkey:   &dsProto.ByteArray{Data: c.EncodePubKey(sign1)},
		Peerinfo: &dsProto.ByteArray{Data: []byte(p1.String())},
	}
	member2 := dsProto.ProtoShardingStructure_Member{
		Pubkey:   &dsProto.ByteArray{Data: c.EncodePubKey(sign2)},
		Peerinfo: &dsProto.ByteArray{Data: []byte(p2.String())},
	}
	member3 := dsProto.ProtoShardingStructure_Member{
		Pubkey:   &dsProto.ByteArray{Data: c.EncodePubKey(sign3)},
		Peerinfo: &dsProto.ByteArray{Data: []byte(p3.String())},
	}
	members := []*dsProto.ProtoShardingStructure_Member{}
	members = append(members, &member1)
	members = append(members, &member2)
	members = append(members, &member3)

	shard := dsProto.ProtoShardingStructure_Shard{
		Members: members,
	}

	shards := []*dsProto.ProtoShardingStructure_Shard{}
	shards = append(shards, &shard)
	shardingStuct := dsProto.ProtoShardingStructure{
		Shards: shards,
	}
	return &shardingStuct
}

func nowAsUnixMilli() int64 {
	now := time.Now()
	log.Printf("Today Date: %v", time.Now().Format(time.ANSIC))
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func getPowDSWinners() []*dsProto.ProtoDSBlock_DSBlockHeader_PowDSWinners {
	sign := c.CreateSignature()

	powDSWinners := []*dsProto.ProtoDSBlock_DSBlockHeader_PowDSWinners{}

	p1 := dsProto.ProtoPeer{
		Ipaddress:      &dsProto.ByteArray{Data: []byte{127, 0, 0, 1}},
		Listenporthost: 8080,
	}

	//Encode Peer

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.

	err := enc.Encode(p1)

	if err != nil {
		log.Fatal("encode error for peer:", err)
	}

	powDSWinner1 := dsProto.ProtoDSBlock_DSBlockHeader_PowDSWinners{
		Key: &dsProto.ByteArray{Data: c.EncodePubKey(sign)},
		Val: &dsProto.ByteArray{Data: network.Bytes()},
	}

	powDSWinners = append(powDSWinners, &powDSWinner1)

	return powDSWinners
}

func getCoSignatures(d *dsProto.ProtoDSCommittee) *dsProto.ProtoBlockBase_CoSignatures {
	dsNodes := d.GetDsnodes()
	ba1 := dsNodes[0].GetPubkey()
	ba2 := dsNodes[1].GetPubkey()
	ba3 := dsNodes[2].GetPubkey()

	data1 := ba1.Data
	data2 := ba2.Data
	data3 := ba3.Data

	sign1 := c.DecodePubKey(data1)
	sign2 := c.DecodePubKey(data2)
	sign3 := c.DecodePubKey(data3)

	res := c.Verify("Secret message", *sign1, (*sign1).R)

	fmt.Printf("res is: %v\n", res)

	//log.Printf("sign pubKey: %v \n", sign1.S)
	var network1 bytes.Buffer         // Stand-in for a network connection
	enc1 := gob.NewEncoder(&network1) // Will write to network.

	err1 := enc1.Encode(sign1)
	if err1 != nil {
		log.Fatal("encode error1:", err1)
	}

	var network2 bytes.Buffer         // Stand-in for a network connection
	enc2 := gob.NewEncoder(&network2) // Will write to network.

	err2 := enc2.Encode(sign2)
	if err2 != nil {
		log.Fatal("encode error2:", err2)
	}

	var network3 bytes.Buffer         // Stand-in for a network connection
	enc3 := gob.NewEncoder(&network3) // Will write to network.

	err3 := enc3.Encode(sign3)
	if err3 != nil {
		log.Fatal("encode error3:", err3)
	}

	B1 := []bool{true, true, true}

	B2 := []bool{true, true, true}

	cosigs := dsProto.ProtoBlockBase_CoSignatures{
		Cs1: &dsProto.ByteArray{Data: network1.Bytes()},
		B1:  B1,
		Cs2: &dsProto.ByteArray{Data: network2.Bytes()},
		B2:  B2,
	}

	return &cosigs
}
