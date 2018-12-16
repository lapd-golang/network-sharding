package block

import (
	tx "../../../directoryService"
	s "../../../network"
)

//NodeDSBlock ...
type NodeDSBlock struct {
	Shardid     uint32
	DSblock     *DSBlock
	Sharding    *s.ShardingStructure
	Assignments *tx.TxSharingAssignments
}
