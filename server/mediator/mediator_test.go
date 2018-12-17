package mediator

import (
	"testing"
)

func TestCheckWhetherBlockIsLatest(t *testing.T) {
	var lb uint64 = 1 //last block number
	var cb uint64 = 2 //current block number

	var le uint64 = 1 //last epoch number
	var ce uint64 = 2 //current epoch number

	if !CheckWhetherBlockIsLatest(cb, ce, lb, le) {
		t.Errorf("Block is not latest")
	}

}

func TestCheckWhetherBlockIsNotLatest(t *testing.T) {
	var lb uint64 = 2 //last block number
	var cb uint64 = 2 //current block number

	var le uint64 = 2 //last epoch number
	var ce uint64 = 2 //current epoch number

	if CheckWhetherBlockIsLatest(cb, ce, lb, le) {
		t.Errorf("Block is not latest")
	}

}
