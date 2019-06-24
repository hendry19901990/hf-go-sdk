package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

)


func (setup *FabricSetup) Invoke(from, to, value  string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "move")
	args = append(args, from)
	args = append(args, to)
	args = append(args, value)


	if setup.client == nil {
		var err error
		clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))

		setup.client, err = channel.New(clientContext)
		if err != nil {
			return "",  fmt.Errorf("failed to create new channel client %v" , err)
		}
	}

	// Add data that will be visible in the proposal, like a description of the invoke request

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: "invoke", Args: [][]byte{[]byte(args[0]),[]byte(args[1]),[]byte(args[2]),[]byte(args[3])}})
	if err != nil {
	//	panic(fmt.Sprintf("failed to create new channel client"))
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission


	return string(response.TransactionID), nil
}
