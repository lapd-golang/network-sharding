package block

import bh "../blockheader"

//DSBlock ...
type DSBlock struct {
	Header    *bh.DSBlockHeader
	Blockbase *Base
}
