package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"encoding/json"

	c "../../crypto"
	n "../../network"
	b "../blockData/block"
	bh "../blockData/blockheader"

	ds "../../../directoryServiceProto"
	m "../../mediator"
)

//DSBlockChain ...
type DSBlockChain struct {
	DSBlocks []*b.DSBlock
}

//NewDSBlock : Creates a new DSBlock and returns it
func NewDSBlock() *b.DSBlock {
	dsBlock := &b.DSBlock{}
	return dsBlock
}

//LastDSBlock ...
func (bc *DSBlockChain) LastDSBlock() *b.DSBlock {

	noofblocks := len(bc.DSBlocks)
	lastBlock := bc.DSBlocks[noofblocks-1]
	return lastBlock
}

//AddDSBlock ...
func (bc *DSBlockChain) AddDSBlock(newDSBlock *b.DSBlock) *b.DSBlock {
	bc.DSBlocks = append(bc.DSBlocks, newDSBlock)
	return newDSBlock
}

//GenesisBlock ..
type GenesisBlock interface {
	InitializeGenesisDSBlock(dsBlockChain *DSBlockChain) *b.DSBlock
}

//GenesisDSBlock ...
type GenesisDSBlock struct {
	DSblock *b.DSBlock
}

//InitializeGenesisDSBlock implementation
func (gb *GenesisDSBlock) InitializeGenesisDSBlock(dsBlockChain *DSBlockChain) *b.DSBlock {
	gb.DSblock = constructGenesisDSBlock()
	return gb.DSblock
}

func constructGenesisDSBlock() *b.DSBlock {
	const pubKeyHex = "1c5dbfb5114647061669feecb4c20dcc8f490b263ab7c2f06ceadeb4e1f92cea"
	//const priKeyHex = "BCCDF94ACEC5B6F1A2D96BDDC6CBE22F3C6DFD89FD791F18B722080A908253CD"
	pubKey, _ := hex.DecodeString(pubKeyHex)

	powDswinners := make(map[string]n.Peer)

	v := n.Peer{IP: &[4]byte{127, 0, 0, 0}, ListenPortHost: 8080}
	powDswinners[string(pubKey[:])] = v
	//fmt.Printf("pub key in map v=%v k=%v\n", pubKey, powDswinners[string(pubKey)])
	header := bh.DSBlockHeader{
		Dsdifficulty: bh.DSPOWDIFFICULTY,
		Difficulty:   bh.POWDIFFICULTY,
		Prevhash:     []byte{0},
		Leaderpubkey: pubKey,
		Blocknum:     0,
		Epochnum:     0,
		Gasprice:     []byte{0},
		PoWDswinners: powDswinners,
		Hash: &bh.DSBlockHashSet{
			Shardinghash:  []byte{0},
			Txsharinghash: []byte{0},
			Reservedfield: []byte{0},
		},
	}

	//convert the block header to hash value
	outHeader, err := json.Marshal(header)

	if err != nil {
		log.Fatalf("Unable to convert block header struct to string: %v", err)
	}
	currentTime := nowAsUnixMilli()

	//Get SHA 256 value
	hash := c.Hash(string(outHeader))
	fmt.Printf("hash during genesis creation: %v\n", hash)
	dsBlock := b.DSBlock{
		Header: &header,
		Blockbase: &b.Base{
			Blockhash: hash,
			Cosigs:    nil,
			Timestamp: uint64(currentTime),
		},
	}
	return &dsBlock
}

func nowAsUnixMilli() int64 {
	now := time.Now()
	log.Printf("Today Date: %v", time.Now().Format(time.ANSIC))
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func addGenesisDSBlockToBlockChain(dsBlockChain *DSBlockChain, dsBlock *b.DSBlock) *b.DSBlock {
	return dsBlockChain.AddDSBlock(dsBlock)
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *DSBlockChain {
	dsBlockChain := &DSBlockChain{nil}
	genesisBlock := GenesisDSBlock{}

	dsBlock := genesisBlock.InitializeGenesisDSBlock(dsBlockChain)
	dsBlockChain.AddDSBlock(dsBlock)

	fmt.Printf("Directory service blockchain started\n")
	return dsBlockChain
}

//ProcessNodeSBlock ...
// Returns last block in case fails to verify the nod DS Block
func ProcessNodeSBlock(n *ds.NodeDSBlock, blocks []*b.DSBlock) (b.DSBlock, bool) {

	noofblocks := len(blocks)

	if noofblocks == 0 {
		return b.DSBlock{}, false
	}

	lastBlock := blocks[noofblocks-1]

	protoDSBlock := n.GetDsblock()
	protoBlockHeader := protoDSBlock.GetHeader()

	// Checking for freshness of incoming DS Block
	b := protoBlockHeader.GetBlocknum()
	e := protoBlockHeader.GetEpochnum()

	lb := lastBlock.Header.Blocknum
	le := lastBlock.Header.Epochnum

	if !m.CheckWhetherBlockIsLatest(b, e, lb, le) {
		return *lastBlock, false
	}

	log.Println("Node DS Block is latest")

	//Verify the DSBlockHashSet member of the DSBlockHeader

	//Verify timestamp

	timestamp := protoDSBlock.GetBlockbase().GetTimestamp()
	timeoutInSec := bh.CONSENSUSOBJECTTIMEOUT + (bh.TXDISTRIBUTETIMEINMS)/1000
	if !m.VerifyTimestamp(timestamp, uint64(timeoutInSec)) {
		log.Println("Timestamp is not verified")
		return *lastBlock, false
	}

	log.Println("Timestamp is verified")

	//Check Sharding hash
	shardinghash := c.Hash(n.GetSharding().String())
	expectedshardinghash := protoBlockHeader.GetHash().GetShardinghash()

	if !bytes.Equal(shardinghash, expectedshardinghash) {
		log.Printf("Sharding structure hash in newly received DS Block doesn't match. Calculated:%v and received: %v\n", shardinghash, expectedshardinghash)
		return *lastBlock, false
	}

	log.Println("Sharding hashes are matched")

	// Check Tx sharing structure hash
	txSharingHash := c.Hash(n.GetAssignments().String())
	expectedTXSharingHash := protoBlockHeader.GetHash().GetTxsharinghash()

	if !bytes.Equal(txSharingHash, expectedTXSharingHash) {
		log.Printf("Tx sharing structure hash in newly received DS Block doesn't match. Calculated:%v and received: %v\n", shardinghash, expectedshardinghash)
		return *lastBlock, false
	}
	log.Println("Tx sharing structure hashs are matched")

	// Check the signature of this DS block

	cosigs := protoDSBlock.GetBlockbase().GetCosigs()
	if !VerifyDSBlockCoSignature(cosigs) {
		log.Printf("DSBlock cosig verification failed\n")
	}

	log.Println("Cosig verification successfully completed")

	//Map Proto DS Block to required DS Block
	dsBlock := MapProtoDSBlocktoDSBlock(n, lastBlock)
	//dtomapper.MapDSBlockToProtoBuffer(&dsBlock)
	blocks = append(blocks, &dsBlock)
	return dsBlock, true
}

//MapProtoDSBlocktoDSBlock ...
func MapProtoDSBlocktoDSBlock(nb *ds.NodeDSBlock, lastBlock *b.DSBlock) b.DSBlock {

	protoDSBlock := nb.GetDsblock()

	powDSWinner := *protoDSBlock.GetHeader().GetDswinners()[0]

	powDSWinnerBytes := powDSWinner.GetKey().GetData()

	peerInfo := powDSWinner.GetVal().GetData()

	//Decode Peer from Proto buffer
	decoder := bytes.NewBuffer(peerInfo) // Stand-in for a network connection
	dec := gob.NewDecoder(decoder)       // Will read from network.

	var peer ds.ProtoPeer

	err := dec.Decode(&peer)
	if err != nil {
		log.Fatal("peer info decode error:", err)
	}

	sign := c.DecodePubKey(powDSWinnerBytes)
	winPubkey := c.PublicKey("Secret message", *sign)
	PeerInfo := n.Peer{
		IP:             &[4]byte{127, 0, 0, 0},
		ListenPortHost: peer.Listenporthost,
	}
	PeerInfoMap := make(map[string]n.Peer)

	PeerInfoMap[winPubkey.String()] = PeerInfo
	dsBlockHeader := bh.DSBlockHeader{
		Dsdifficulty: protoDSBlock.GetHeader().GetDsdifficulty(),
		Difficulty:   protoDSBlock.GetHeader().GetDifficulty(),
		Prevhash:     lastBlock.Header.Prevhash,
		Leaderpubkey: c.EncodePubKey(c.CreateSignature()),
		Blocknum:     protoDSBlock.GetHeader().GetBlocknum(),
		Epochnum:     protoDSBlock.GetHeader().GetEpochnum(),
		Gasprice:     []byte{5},
		Swinfo:       []byte{0, 0, 0, 0},
		PoWDswinners: PeerInfoMap,
		Hash: &bh.DSBlockHashSet{
			Shardinghash:  protoDSBlock.GetHeader().GetHash().GetShardinghash(),
			Txsharinghash: protoDSBlock.GetHeader().GetHash().GetShardinghash(),
			Reservedfield: protoDSBlock.GetHeader().GetHash().GetReservedfield(),
		},
		Committeehash: protoDSBlock.GetHeader().GetCommitteehash(),
	}
	cosigs := *(protoDSBlock.GetBlockbase().GetCosigs())

	cs1p := *(cosigs.GetCs1())
	cs1d := cs1p.Data
	cs1 := c.DecodePubKey(cs1d)

	cs2p := *(cosigs.GetCs2())
	cs2d := cs2p.Data
	cs2 := c.DecodePubKey(cs2d)

	coSignatures := b.CoSignatures{
		Cs1: *cs1,
		B1:  cosigs.GetB1(),
		Cs2: *cs2,
		B2:  cosigs.GetB2(),
	}
	base := b.Base{
		Blockhash: protoDSBlock.GetBlockbase().GetBlockhash(),
		Cosigs:    &coSignatures,
		Timestamp: protoDSBlock.GetBlockbase().GetTimestamp(),
	}
	dsBlock := b.DSBlock{
		Header:    &dsBlockHeader,
		Blockbase: &base,
	}
	fmt.Printf("dsBlockHeader size: %v\n", dsBlock)

	return dsBlock
}

//VerifyDSBlockCoSignature ...
// cosigs: ProtoBlockBase_CoSignatures
func VerifyDSBlockCoSignature(cosigs *ds.ProtoBlockBase_CoSignatures) bool {
	B1 := cosigs.GetB1()

	log.Printf("DS Committee size is: %v\n", len(B1))

	cs1 := *cosigs.GetCs1()
	sign1 := c.DecodePubKey(cs1.Data)

	B2 := cosigs.GetB2()

	cs2 := *cosigs.GetCs2()
	sign2 := c.DecodePubKey(cs2.Data)

	message := "Secret message"

	isSign1Verified := c.Verify(message, *sign1, (*sign1).R)
	isSign2Verified := c.Verify(message, *sign2, (*sign2).R)

	//DS Committee size should be 3

	if len(B1) != 3 {
		log.Printf("Mismatch: DS committee size = 3, , co-sig bitmap size = %v\n", len(B1))
		return false
	}
	if len(B2) != 3 {
		log.Printf("Mismatch: DS committee size = 3, , co-sig bitmap size = %v\n", len(B2))
		return false
	}
	if !isSign2Verified && isSign1Verified {
		log.Printf("DSBlock cosig verification failed for Signature 1\n")
		return false
	}

	return true
}
