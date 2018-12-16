package directoryservice

import (
	p "../network"
	n "../node"
)

//TxSharingAssignments ...
type TxSharingAssignments struct {
	DSnodes    []*n.Node
	Shardnodes []*AssignedNodes
}

//AssignedNodes ...
type AssignedNodes struct {
	Receivers []p.Peer
	Senders   []p.Peer
}
