# Indexer proxy

## eth_getLogs

Returns an array of all logs matching a given filter object.

##### Parameters

1. `Object` - The filter options:
  - `fromBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `toBlock`: `QUANTITY|TAG` - (optional, default: `"latest"`) Integer block number, or `"latest"` for the last mined block or `"pending"`, `"earliest"` for not yet mined transactions.
  - `address`: `DATA|Array`, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.
  - `topics`: `Array of DATA`,  - (optional) Array of 32 Bytes `DATA` topics. Each topic can also be an array of DATA with "or" options.

```js
params: [{
  "fromBlock": "0x1",
  "toBlock": "0x2",
  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
  "topics": ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b", null, [0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b, 0x000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc]]
}]

```

##### Request

```js
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"blockHash": "0x7c5a35e9cb3e8ae0e221ab470abae9d446c3a5626ce6689fc777dcffcab52c70", "topics":["0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80"]}],"id":1}'
```

##### Returns

```json
{
   "jsonrpc":"2.0",
   "id":1,
   "result":[
      {
         "address":"0xee6e1bd3b1aa5c0d6edf2809d40ac345c309b228",
         "topics":[
            "0x699d391aca8d4d4d43e61892e60e3baac3d506825f8f3de8acfbc998d949a308"
         ],
         "data":"0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000157b2275736572223a224a6f686e20536d697468227d0000000000000000000000",
         "blockNumber":"0xe825bb",
         "transactionHash":"0x883753df7bc7eb5a7b3e02ee3c41efb4ca4a9dea6c59d040a00b7b15425b347e",
         "transactionIndex":"0x0",
         "blockHash":"0x9a09fc11a1157055062d2ae1a1742a91ae57191ec5ee206806b2314f18ba0501",
         "logIndex":"0x0",
         "removed":false
      }
   ]
}
```