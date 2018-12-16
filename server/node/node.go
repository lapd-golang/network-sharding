package node

import (
	c "../crypto"
	p "../network"
)

//ProcessBlockAction ...
const ProcessBlockAction = "PROCESS_DSBLOCK"

//Node ...
type Node struct {
	PubKey c.PubKey
	Peer   p.Peer
}
