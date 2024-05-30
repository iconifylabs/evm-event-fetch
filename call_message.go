package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func call_message(startBlock uint64, detail_logs bool) (map[*big.Int]*big.Int, map[*big.Int]string, []*big.Int, error) {
	sequence_numbers := []*big.Int{}
	request_ids := []*big.Int{}
	txHashes := make(map[*big.Int]string)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get latest block header: %v", err)
	}

	latestBlock := header.Number.Uint64()

	contractAddress := common.HexToAddress(contractAddr)
	eventSigHash := crypto.Keccak256Hash([]byte(CALL_MESSAGE_EVENT))

	contractAbi, err := abi.JSON(strings.NewReader(CALL_MESSAGE_ABI))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse contract ABI: %v", err)
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
			return nil, nil, nil, fmt.Errorf("failed to filter logs: %v", err)
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
					return nil, nil, nil, fmt.Errorf("failed to unpack log data: %v", err)
				}

				event.From = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
				event.To = common.BytesToAddress(vLog.Topics[2].Bytes()).Hex()
				event.SN = new(big.Int).SetBytes(vLog.Topics[3].Bytes())

				if detail_logs {
					fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
					fmt.Printf("Txn Hash: %s\n", vLog.TxHash.Hex())
					fmt.Printf("From: %s\n", event.From)
					fmt.Printf("To: %s\n", event.To)
					fmt.Printf("SN: %s\n", event.SN.String())
					fmt.Printf("ReqId: %s\n", event.ReqId.String())
					fmt.Printf("Data: %x\n", event.Data)
					fmt.Println()
				}
				sequence_numbers = append(sequence_numbers, event.SN)
				request_ids = append(request_ids, event.ReqId)
				txHashes[event.ReqId] = vLog.TxHash.Hex()
			}
		}
	}

	mappings := make(map[*big.Int]*big.Int)
	for i := 0; i < len(sequence_numbers); i++ {
		mappings[request_ids[i]] = sequence_numbers[i]
	}

	return mappings, txHashes, request_ids, nil
}
