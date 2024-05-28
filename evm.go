package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL         = "https://api.avax.network/ext/bc/C/rpc"
	contractAddr   = "0xfC83a3F252090B26f92F91DFB9dC3Eb710AdAf1b"
	blockBatchSize = 2000
	eventSig       = "CallMessage(string,string,uint256,uint256,bytes)"
)

func evm_main() {
	sequence_numbers := []*big.Int{}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get latest block header: %v", err)
	}

	latestBlock := header.Number.Uint64()
	startBlock := uint64(40428048) // Xcall Deploy Height

	contractAddress := common.HexToAddress(contractAddr)
	eventSigHash := crypto.Keccak256Hash([]byte(eventSig))

	contractAbi, err := abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"_from","type":"string"},{"indexed":true,"internalType":"string","name":"_to","type":"string"},{"indexed":true,"internalType":"uint256","name":"_sn","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_reqId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"_data","type":"bytes"}],"name":"CallMessage","type":"event"}]`))
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

				fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
				// fmt.Printf("From: %s\n", event.From)
				// fmt.Printf("To: %s\n", event.To)
				fmt.Printf("SN: %s\n", event.SN.String())
				// fmt.Printf("ReqId: %s\n", event.ReqId.String())
				// fmt.Printf("Data: %s\n", event.Data)
				fmt.Println()
				sequence_numbers = append(sequence_numbers, event.SN)
			}
		}
	}

	fmt.Printf("%+v\n\n", sequence_numbers)
}
