package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Query query the chaincode to get the state of hello
func (setup *FabricSetup) Query(id string) (string, error) {

	if setup.client == nil {
		var err error
		clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
		setup.client, err = channel.New(clientContext)
		if err != nil {
			return "",  fmt.Errorf("failed to create new channel client %v" , err)
		}
	}

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: "invoke", Args: [][]byte{[]byte("query"), []byte(id)}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
