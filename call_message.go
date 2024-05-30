package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func call_message(startBlock uint64, detail_logs bool) {
	sequence_numbers := []*big.Int{}
	request_ids := []*big.Int{}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get latest block header: %v", err)
	}

	latestBlock := header.Number.Uint64()

	contractAddress := common.HexToAddress(contractAddr)
	eventSigHash := crypto.Keccak256Hash([]byte(CALL_MESSAGE_EVENT))

	contractAbi, err := abi.JSON(strings.NewReader(CALL_MESSAGE_ABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	for start := startBlock; start < latestBlock; start += blockBatchSize {
		end := start + blockBatchSize - 1
		if end > latestBlock {
			end = latestBlock
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(start)),
			ToBlock:   big.NewInt(int64(end)),
			Addresses: []common.Address{contractAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to filter logs: %v", err)
		}

		for _, vLog := range logs {
			if vLog.Topics[0] == eventSigHash {
				var event struct {
					From  string
					To    string
					SN    *big.Int
					ReqId *big.Int
					Data  []byte
				}

				err := contractAbi.UnpackIntoInterface(&event, "CallMessage", vLog.Data)
				if err != nil {
					log.Fatalf("Failed to unpack log data: %v", err)
				}

				event.From = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
				event.To = common.BytesToAddress(vLog.Topics[2].Bytes()).Hex()
				event.SN = new(big.Int).SetBytes(vLog.Topics[3].Bytes())

				if detail_logs {
					fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
					fmt.Printf("Txn Hash: %s\n", vLog.TxHash)
					// fmt.Printf("From: %s\n", event.From)
					// fmt.Printf("To: %s\n", event.To)
					fmt.Printf("SN: %s\n", event.SN.String())
					fmt.Printf("ReqId: %s\n", event.ReqId.String())
					// fmt.Printf("Data: %s\n", event.Data)
					fmt.Println()
				}
				sequence_numbers = append(sequence_numbers, event.SN)
				request_ids = append(request_ids, event.ReqId)
			}
		}
	}
	mappings := make(map[*big.Int]*big.Int)

	for i := 0; i < len(sequence_numbers); i++ {
		mappings[request_ids[i]] = sequence_numbers[i]
	}

	jsonData, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	file, err := os.Create("mappings.json")
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to JSON file: %v", err)
	}

	fmt.Println("Mappings saved to mappings.json")
}
