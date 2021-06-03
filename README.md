# dummy_rpc_proxy

It decode the input, encoded it again and pass downstream and return back to the response back to the requester

`/` - main endpoint we you should send `POST` requests only

`/healthz`

`/make_it_fail`

`/make_it_work`

## Compile for alpine

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
