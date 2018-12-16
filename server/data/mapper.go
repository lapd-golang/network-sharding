package dtomapper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unsafe"

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

		for pubkey, peer := range pwoDswinners {

			b := []byte(pubkey)
			bs := fmt.Sprintf("%v", b)
			var pubkeyBytes []byte

			for _, ps := range strings.Split(strings.Trim(bs, "[]"), " ") {
				pi, _ := strconv.Atoi(ps)
				pubkeyBytes = append(pubkeyBytes, byte(pi))
			}

			fmt.Printf("pubkeyBytes=%v \n", pubkeyBytes)

			buf := &bytes.Buffer{}
			err := binary.Write(buf, binary.LittleEndian, *peer.IP)

			if err != nil {
				log.Fatalf("Failed to serialize the Peer information %v", err)
			}

			protoPowDSWinner := proto.ProtoDSBlock_DSBlockHeader_PowDSWinners{
				Key: &proto.ByteArray{Data: pubkeyBytes},
				Val: &proto.ByteArray{Data: buf.Bytes()},
			}

			protoPowDSWinners = append(protoPowDSWinners, &protoPowDSWinner)
		}

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

		if dscosigs == nil {

			//Serialize CS1
			Cs1 := c.Signature{}
			Cs2 := c.Signature{}
			var byteCs1Slice []byte = *(*[]byte)(unsafe.Pointer(&Cs1))
			var byteCs2Slice []byte = *(*[]byte)(unsafe.Pointer(&Cs2))

			protoCosigs = proto.ProtoBlockBase_CoSignatures{
				Cs1: &proto.ByteArray{Data: byteCs1Slice},
				B1:  []bool{},
				Cs2: &proto.ByteArray{Data: byteCs2Slice},
				B2:  []bool{},
			}

		} else {

			cosigs := blockbase.Cosigs
			//Serialize CS1
			Cs1 := cosigs.Cs1
			Cs2 := cosigs.Cs2
			var byteCs1Slice []byte = *(*[]byte)(unsafe.Pointer(&Cs1))
			var byteCs2Slice []byte = *(*[]byte)(unsafe.Pointer(&Cs2))

			protoCosigs = proto.ProtoBlockBase_CoSignatures{
				Cs1: &proto.ByteArray{Data: byteCs1Slice},
				B1:  cosigs.B1,
				Cs2: &proto.ByteArray{Data: byteCs2Slice},
				B2:  cosigs.B2,
			}
		}
		protoBlockBase := proto.ProtoBlockBase{
			Blockhash: blockhash,
			Cosigs:    &protoCosigs,
			Timestamp: blockbase.Timestamp,
		}
		//fmt.Printf("Cosigs: %v\n", *(protoBlockBase.Cosigs))

		protoDSBlock := proto.ProtoDSBlock{
			Header:    &protoHeader,
			Blockbase: &protoBlockBase,
		}
		blocks = append(blocks, &protoDSBlock)
	}

	fmt.Printf("blocks size in mapper: %v\n", len(blocks))

	return blocks
}
