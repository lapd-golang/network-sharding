package network

import (
	c "../crypto"
)

//ShardingStructure ...
type ShardingStructure struct {
	Shards []*Shard
}

//Member ...
type Member struct {
	Pubkey     c.PubKey
	Peerinfo   Peer
	Reputation uint32
}

//Shard ...
type Shard struct {
	Members []*Member
}
