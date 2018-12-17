package dtomapper

import (
	"fmt"
	"log"

	proto "../../directoryServiceProto"
	c "../crypto"
	d "../data/blockData/block"
)

//MapToProtoBuffer ...
func MapToProtoBuffer(DSBlocks []*d.DSBlock, blocks []*proto.ProtoDSBlock) []*proto.ProtoDSBlock {
	log.Printf("DSBlocks at MapToProtoBuffer =%v \n", len(DSBlocks))
	for _, b := range DSBlocks {
		Header := b.Header
		protoPowDSWinners := []*proto.ProtoDSBlock_DSBlockHeader_PowDSWinners{}
		pwoDswinners := Header.PoWDswinners

		leaderSign := Header.Leaderpubkey
		leaderpubkey := c.DecodePubKey(leaderSign)
		peerInfo := pwoDswinners[leaderpubkey.String()]

		log.Printf("ListenPortHost =%v \n", peerInfo.ListenPortHost)
		log.Printf("IP =%v \n", peerInfo.IP)

		var arr []byte
		var vvv []byte
		if peerInfo.IP != nil {
			vvv := *peerInfo.IP
			copy(arr[:], string(vvv))
			vvv = *peerInfo.IP
			copy(arr[:], string(vvv))
		} else {
			vvv = []byte{127, 0, 0, 0}
		}

		protoPowDSWinner := proto.ProtoDSBlock_DSBlockHeader_PowDSWinners{
			Key: &proto.ByteArray{Data: vvv},
		}

		protoPowDSWinners = append(protoPowDSWinners, &protoPowDSWinner)

		protoHeader := proto.ProtoDSBlock_DSBlockHeader{
			Dsdifficulty: Header.Dsdifficulty,
			Difficulty:   Header.Difficulty,
			Prevhash:     Header.Prevhash,
			Leaderpubkey: &proto.ByteArray{Data: Header.Leaderpubkey},
			Blocknum:     Header.Blocknum,
			Epochnum:     Header.Epochnum,
			Gasprice:     &proto.ByteArray{Data: Header.Gasprice},
			Swinfo:       &proto.ByteArray{Data: Header.Swinfo},
			Dswinners:    protoPowDSWinners,
		}

		// Get Blockbase details from DTO
		blockbase := b.Blockbase
		blockhash := blockbase.Blockhash
		dscosigs := blockbase.Cosigs

		var protoCosigs = proto.ProtoBlockBase_CoSignatures{}

		if dscosigs != nil {
			Cs1 := dscosigs.Cs1
			Cs2 := dscosigs.Cs2
			protoCosigs = proto.ProtoBlockBase_CoSignatures{
				Cs1: &proto.ByteArray{Data: c.EncodePubKey(&Cs1)},
				B1:  dscosigs.B1,
				Cs2: &proto.ByteArray{Data: c.EncodePubKey(&Cs2)},
				B2:  dscosigs.B2,
			}
		}

		protoBlockBase := proto.ProtoBlockBase{
			Blockhash: blockhash,
			Cosigs:    &protoCosigs,
			Timestamp: blockbase.Timestamp,
		}

		protoDSBlock := proto.ProtoDSBlock{
			Header:    &protoHeader,
			Blockbase: &protoBlockBase,
		}
		blocks = append(blocks, &protoDSBlock)
	}

	fmt.Printf("blocks size in mapper: %v\n", len(blocks))

	return blocks
}
