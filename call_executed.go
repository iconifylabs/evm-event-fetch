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

func call_executed(startBlock uint64, detail_logs bool) ([]*big.Int, error) {
	request_ids := []*big.Int{}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Avalanche client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block header: %v", err)
	}

	latestBlock := header.Number.Uint64()
	contractAddress := common.HexToAddress(contractAddr)
	eventSigHash := crypto.Keccak256Hash([]byte(CALL_EXECUTED_EVENT))

	contractAbi, err := abi.JSON(strings.NewReader(CALL_EXECUTED_ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
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
			return nil, fmt.Errorf("failed to filter logs: %v", err)
		}

		for _, vLog := range logs {
			if vLog.Topics[0] == eventSigHash {
				var event struct {
					ReqId *big.Int
					Code  *big.Int
					Msg   string
				}

				err := contractAbi.UnpackIntoInterface(&event, "CallExecuted", vLog.Data)
				if err != nil {
					return nil, fmt.Errorf("failed to unpack CallExecuted log data: %v", err)
				}

				event.ReqId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())

				if detail_logs {
					fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
					fmt.Printf("Transaction Hash: %s\n", vLog.TxHash.Hex())
					fmt.Printf("ReqId: %s\n", event.ReqId.String())
					fmt.Println()
				}

				request_ids = append(request_ids, event.ReqId)
			}
		}
	}

	return request_ids, nil
}
