# Indexer proxy

## eth_getLogs

FILTER OBJECT

```
address [optional] - a string representing the address (20 bytes) to check for balance

fromBlock [optional, default is "latest"] - an integer block number, or the string "latest", "earliest" or "pending"

toBlock [optional, default is "latest"] - an integer block number, or the string "latest", "earliest" or "pending"

topics[optional] - Array of 32 Bytes DATA topics. Topics are order-dependent.

blockhash:[optional] With the addition of EIP-234, blockHash restricts the logs returned to the single block with the 32-byte hash 

blockHash. Using blockHash is equivalent to fromBlock = toBlock = the block number with hash blockHash. If blockHash is present in in the filter criteria, then neither fromBlock nor toBlock are allowed.
```

// Request

```

curl https://mainnet.infura.io/v3/YOUR-PROJECT-ID \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"blockHash": "0x7c5a35e9cb3e8ae0e221ab470abae9d446c3a5626ce6689fc777dcffcab52c70", "topics":["0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80"]}],"id":1}'

 curl -XPOST -H "content-type: application/json" --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}],"id":74}' http://localhost:5000

```



```
 curl -XPOST -H "content-type: application/json" --data '{"jsonrpc":"2.0","method":"eth_getLogs","params": [],"id":1}' http://localhost:5000
 ```