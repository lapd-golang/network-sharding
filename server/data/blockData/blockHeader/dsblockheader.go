package blockheader

import (
	p "../../../network"
)

//DSBlockHeader ...
type DSBlockHeader struct {
	Dsdifficulty  uint32
	Difficulty    uint32
	Prevhash      []byte
	Leaderpubkey  []byte
	Blocknum      uint64
	Epochnum      uint64
	Gasprice      []byte
	Swinfo        []byte
	PoWDswinners  map[string]p.Peer // Key : Pubkey and Value : Peer Detail
	Hash          *DSBlockHashSet
	Committeehash []byte
}

//DSBlockHashSet ...
type DSBlockHashSet struct {
	Shardinghash  []byte
	Txsharinghash []byte
	Reservedfield []byte
}
