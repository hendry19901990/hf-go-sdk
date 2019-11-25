package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

)


func (setup *FabricSetup) Invoke(args []string) (string, error) {

	if setup.client == nil {
		var err error
		clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))

		setup.client, err = channel.New(clientContext)
		if err != nil {
			return "",  fmt.Errorf("%v" , err)
		}
	}

	// Add data that will be visible in the proposal, like a description of the invoke request

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: "invoke", Args: setup.parseArgs(args)})
	if err != nil {
	//	panic(fmt.Sprintf("failed to create new channel client"))
		return "", fmt.Errorf("%v", err)
	}

	// Wait for the result of the submission
	return string(response.TransactionID), nil
}

func (setup *FabricSetup) parseArgs(args []string)[][]byte{
	result := make([][]byte, 0)

	for _, arg := range args{
		result = append(result, []byte(arg))
	}

	return result
}
