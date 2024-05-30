package main

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

func main() {
	startBlock := uint64(46048251) // Starting block
	detail_logs := false           // Set to true if you want detailed logs

	mappings, txHashes, _, err := call_message(startBlock, detail_logs)
	if err != nil {
		fmt.Printf("Error in call_message: %v\n", err)
		return
	}

	executedReqIds, err := call_executed(startBlock, detail_logs)
	if err != nil {
		fmt.Printf("Error in call_executed: %v\n", err)
		return
	}

	executedReqIdsSet := make(map[string]bool)
	for _, reqId := range executedReqIds {
		executedReqIdsSet[reqId.String()] = true
	}

	// Extract keys from mappings and sort them
	var snKeys []*big.Int
	for _, sn := range mappings {
		snKeys = append(snKeys, sn)
	}
	sort.Slice(snKeys, func(i, j int) bool { return snKeys[i].Cmp(snKeys[j]) < 0 })

	fmt.Printf("%-15s %-15s %-15s %-40s\n", "SN", "Req_ID", "Delivered", "Txn Hash (if not delivered)")
	fmt.Println(strings.Repeat("-", 90))

	for _, sn := range snKeys {
		var reqId *big.Int
		for k, v := range mappings {
			if v.Cmp(sn) == 0 {
				reqId = k
				break
			}
		}
		delivered := "No"
		txnHash := ""
		if executedReqIdsSet[reqId.String()] {
			delivered = "Yes"
		} else {
			txnHash = txHashes[reqId]
		}
		fmt.Printf("%-15s %-15s %-15s %-40s\n", sn.String(), reqId.String(), delivered, txnHash)
	}
}
