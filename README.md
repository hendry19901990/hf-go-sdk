# ABL HFL

This tutorial has been made on **Ubuntu 16.04** .

```bash
cd configuration
```

## Fabric key material
Certificates x509

```bash
./bin/cryptogen generate --config=./crypto-config.yaml
```

The orderer genesis block, channel tx and anchor tx

```bash
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ABL -outputBlock ./artifacts/orderer.genesis.block
```

```bash
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ABL -outputCreateChannelTx ./artifacts/abl.channel.tx -channelID abl
```

```bash
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ABL -outputAnchorPeersUpdate ./artifacts/org1.abl.anchors.tx -channelID abl -asOrg ABLOrganization1
```

###Now that your Hyperledger-Fabric Network is perfectly set up you can launch it :
```bash
cd configuration
docker-compose up -d
```

###Build and run the backend
```bash
dep ensure
go build
./hf-go-sdk
```
