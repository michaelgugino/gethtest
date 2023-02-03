# geth node + racecouse on kind

## setup kind
On Ubuntu Jammy, run `install-all.sh` to install needed dependencies

Make sure you add necessary dirs to your PATH.

## geth setup

### cli Options
```
podman run -d -p 30303:30303 -p 8545:8545 ethereum/client-go:v1.10.25 --http --http.addr "0.0.0.0" --http.corsdomain="*"  --http.port 8545 --http.vhosts "*" --http.api 'admin,personal,eth,net,web3,txpool,miner,clique,debug' --password /data/password --allow-insecure-unlock --nodiscover --verbosity 4 --miner.gastarget 804247552 --miner.gasprice 0 --mine --datadir /data --networkid 1999 --miner.etherbase=0x00000000000000000000000000000000000000af
```

--syncmode 'full' --port 30311 --http --http.addr "0.0.0.0" --http.corsdomain="*"  -http.port 8545 --http.vhosts "*" --http.api 'admin,personal,eth,net,web3,txpool,miner,clique,debug' --networkid %d --miner.gasprice 0 --password /data/password --mine --allow-insecure-unlock --nodiscover --verbosity 4 --miner.gaslimit 16777215

### Add user:
```
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"personal_newAccount","params":["password"],"id":1}' localhost:8545

curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method": "personal_unlockAccount", "params": ["0xa1a58b295f10f4e52480a89340c0444a5e49cbca", "password", 0]}' localhost:8545
```

## build and install operator

## Expose service
```
kubectl port-forward svc/racecourse-sample --address 0.0.0.0 3000:3000 &
```
