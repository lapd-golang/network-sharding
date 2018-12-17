package mediator

import "log"

//CheckWhetherBlockIsLatest ...
// b: Block Number
// e: Epoch Number
// lb: Last block Epoch Number
// le: Last block Epoch Number
func CheckWhetherBlockIsLatest(b, e, lb, le uint64) bool {
	if lb != (b - 1) {
		log.Printf("Block numbers did not match. Expected block no:%v, Actual block no:%v\n", b, lb)
		return false
	}

	if le != (e - 1) {
		log.Printf("Epoch numbers did not match. Expected block no:%v, Actual block no:%v\n", b, lb)
		return false
	}

	return true
}

//VerifyTimestamp ...
func VerifyTimestamp(timestampInMicrosec, timeoutInSec uint64) bool {
	// systemTimestampVarianceInSecond := 3600
	// var lowBound int64
	// var hiBound int64
	return true
}
