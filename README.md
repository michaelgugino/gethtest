# geth node + racecouse on kind

## setup kind
On Ubuntu Jammy, run `install-all.sh` to install needed dependencies

Make sure you add necessary go dirs to your PATH.

## geth setup

### cli Options
```
podman run -d -p 30303:30303 -p 8545:8545 ethereum/client-go:v1.10.25 --http --http.addr "0.0.0.0" --http.corsdomain="*"  --http.port 8545 --http.vhosts "*" --http.api 'admin,personal,eth,net,web3,txpool,miner,clique,debug' --password /data/password --allow-insecure-unlock --nodiscover --verbosity 4 --miner.gastarget 804247552 --miner.gasprice 0 --mine --datadir /data --networkid 1999 --miner.etherbase=0x00000000000000000000000000000000000000af
```

### Add user:
```
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"personal_newAccount","params":["password"],"id":1}' localhost:8545

curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method": "personal_unlockAccount", "params": ["changeme_aaaaa", "password", 0]}' localhost:8545
```

If the node is deployed on k8s, be sure to port-forward so you can run these commands.

## build and install operator
Use the Makefile.

Everything is already built and published.  Should just need to run
```
make install
make deploy
```

## Deploy && Expose service

```
kubectl apply -f operator/config/samples/gethtest_v1_racecourse.yaml
kubectl port-forward svc/racecourse-sample --address 0.0.0.0 3000:3000 &
```
