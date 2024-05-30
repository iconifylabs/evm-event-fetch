package main

const (
	rpcURL              = "https://api.avax.network/ext/bc/C/rpc"
	contractAddr        = "0xfC83a3F252090B26f92F91DFB9dC3Eb710AdAf1b"
	blockBatchSize      = 2000
	CALL_MESSAGE_EVENT  = "CallMessage(string,string,uint256,uint256,bytes)"
	CALL_EXECUTED_EVENT = "CallExecuted(uint256,int256,string)"

	CALL_EXECUTED_ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"_reqId","type":"uint256"},{"indexed":false,"internalType":"int256","name":"_code","type":"int256"},{"indexed":false,"internalType":"string","name":"_msg","type":"string"}],"name":"CallExecuted","type":"event"}]`
	CALL_MESSAGE_ABI  = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"_from","type":"string"},{"indexed":true,"internalType":"string","name":"_to","type":"string"},{"indexed":true,"internalType":"uint256","name":"_sn","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_reqId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"_data","type":"bytes"}],"name":"CallMessage","type":"event"}]`
)
