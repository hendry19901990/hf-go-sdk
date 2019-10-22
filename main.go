package main

import (
	"fmt"
	"github.com/hendry19901990/hf-go-sdk/blockchain"
	"os"
)

const INIT_ = false

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer.hf.abl.io",

		// Channel parameters
		ChannelID:     "abl",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hendry19901990/hf-go-sdk/configuration/channel-artifacts/channel.tx",

		// Chaincode parameters
		ChainCodeID:     "mycc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hendry19901990/hf-go-sdk/blockchain/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}



	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize(INIT_)
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC(INIT_)
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

  var result string

  result, err =  fSetup.Invoke("a","b","5")
  if err != nil {
		fmt.Printf("Unable to invoke the chaincode: %v\n", err)
		return
	}
  fmt.Println(result)

  result, err =  fSetup.Query("a")
  if err != nil {
		fmt.Printf("Unable to invoke the chaincode: %v\n", err)
		return
	}
  fmt.Println(result)


}
