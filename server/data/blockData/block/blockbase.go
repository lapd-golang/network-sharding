package block

import (
	c "../../../crypto"
)

//Base (Block Base) of the Block ...
type Base struct {
	Blockhash []byte
	Cosigs    *CoSignatures
	Timestamp uint64
}

//CoSignatures ...
type CoSignatures struct {
	Cs1 c.Signature
	B1  []bool
	Cs2 c.Signature
	B2  []bool
}
