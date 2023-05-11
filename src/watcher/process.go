package watcher

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func process(logEntry types.Log) {
	//log.Println(logEntry.TxHash.Hex())
	// Write event to database
	processWh(logEntry)
	processDb(logEntry)
}
