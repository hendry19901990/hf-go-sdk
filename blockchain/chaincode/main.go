package main

import (
	"log"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"encoding/base64"
)

// ABLChaincode implementation of Chaincode
type ABLChaincode struct {
}

type RequestAccess struct {
  ReqType              string       `json:"req_type"` // "recruiter" || "employer" || "landlord" || other
  ReqName              string       `json:"req_name"`
  ReqCompany           string       `json:"req_company"`  //optional
  ReqJobPosition       string       `json:"req_job_position"`   //optional
  Accepted             bool         `json:"accepted"`
	City                 string       `json:"city"  gorm:"column:city"`  //optional
	State                string       `json:"state"  gorm:"column:state"`  //optional
	CreatedDate          time.Time    `json:"created_date"`
	AcceptedDate         time.Time    `json:"accepted_date"`
  DNATypes             []int64      `json:"dna_types"`
}

const (
	PERMISSIONS    = "permissions"
	REQUEST_ACCESS = "request_access"
)

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ABLChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("########### ABL Chaincode Init REQUEST ###########")

  permissions  := make(map[string][]string) // user_temp_id + user_id => [dna_types]
  requests     := make(map[string]map[string]RequestAccess) //  user_id => map(req_email => RequestAccess)
	// Write the state to the ledger
	if err := stub.PutState(PERMISSIONS, t.toJsonPermission(permissions)); err != nil {
		return shim.Error(err.Error())
	}
	if err := stub.PutState(REQUEST_ACCESS, t.toJsonRequests(requests)); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *ABLChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("########### ABLChaincode Invoke ###########")

	// Get the function and arguments from the request
	_, args := stub.GetFunctionAndParameters()
	function := args[0]

	if function == "request_accept" {
		// a user accepts a request receives arg[1] user_temp_id, arg[2]  user_id, arg[3] a base64 of the string array [dna_types]
		return t.request_accept(stub, args)
	}else if function   == "request_access" {
		  // an user make a request access receives arg[1] user_id,  arg[2] req_email, arg[3] RequestAccess object in base64
			return t.request_access(stub, args)
	} else if function == "query" {
		// the  "Query" is arg[1] (user_temp_id + user_id)
		return t.query(stub, args)
	} else if function == "permissions" {
		return t.permissions(stub, args)
	} else if function == "get_all_request_access" {
		return t.getAllRequestAccess(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// query
// Every readonly functions in the ledger will be here
func (t *ABLChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("########### ABLChaincode query ###########")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting user_temp_key of the user to query")
	}

	user_temp_key        := args[1]
	permissions_bts, err := stub.GetState(PERMISSIONS)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	permissions := t.fromJsonPermission(permissions_bts)

	dnas, ok    := permissions[user_temp_key]
	if !ok {
		dnas = make([]string, 0)
	}

  dnas_bytes, _ := json.Marshal(dnas)
	jsonResp := "{\"dnas\":"  + string(dnas_bytes) + "}"
	log.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(dnas_bytes)
}

func (t *ABLChaincode) permissions(stub shim.ChaincodeStubInterface, args []string) pb.Response {
   log.Println("########### ABLChaincode permissions ###########")

	 permissions_bts, err := stub.GetState(PERMISSIONS)
 	 if err != nil {
 		  return shim.Error("Failed to get state")
 	 }

	 return shim.Success(permissions_bts)
}

func (t *ABLChaincode) getAllRequestAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
   log.Println("########### ABLChaincode getAllRequestAccess ###########")

	 requests_bts, err := stub.GetState(REQUEST_ACCESS)
 	 if err != nil {
 		  return shim.Error("Failed to get state")
 	 }

	 return shim.Success(requests_bts)
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *ABLChaincode) request_accept(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	  log.Println("########### ABLChaincode request_accept ###########")
		if len(args) != 5 {
			return shim.Error("Incorrect number of arguments. Expecting 5")
		}

		userTempID := args[1]
		userID     := args[2]

		user_temp_key := userTempID + "_" + userID
		permissions_bts, err := stub.GetState(PERMISSIONS)
		if err != nil {
			return shim.Error("Failed to get state permissions")
		}
		permissions := t.fromJsonPermission(permissions_bts)
		dnas          := t.fromBase64ToArray(args[3])
		permissions[user_temp_key] = dnas
		// Write the state back to the ledger
		err = stub.PutState(PERMISSIONS, t.toJsonPermission(permissions))
		if err != nil {
			return shim.Error(err.Error())
		}

		//upgrade state of request_access
		requests_bts, err := stub.GetState(REQUEST_ACCESS)
		req_email := args[4]
		if err != nil {
			return shim.Error("Failed to get state request_access")
		}
		requests := t.fromJsonRequests(requests_bts)
		user_map, ok_user := requests[userID]
		if ok_user {
			request_obj, _          := user_map[req_email]
			request_obj.Accepted     = true
			request_obj.AcceptedDate = time.Now()
			user_map[req_email]      = request_obj
			requests[userID]         = user_map
			user_map_bytes, _ := json.Marshal(user_map)
			log.Printf("Invoke REQUEST_ACCESS UserMap: %s\n", string(user_map_bytes))
		}
		err = stub.PutState(REQUEST_ACCESS, t.toJsonRequests(requests))
		if err != nil {
			return shim.Error(err.Error())
		}
		//upgrade state of request_access

    return shim.Success(nil)
}

func (t *ABLChaincode) request_access(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Println("########### ABLChaincode request_access ###########")
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	requests_bts, err := stub.GetState(REQUEST_ACCESS)
	if err != nil {
		return shim.Error("Failed to get state request_access")
	}
	requests := t.fromJsonRequests(requests_bts)

	userID       := args[1]
	req_email    := args[2]
	request_obj  := t.fromBase64ToRequestObj(args[3])

	user_map, ok_user := requests[userID]
	if !ok_user {
		user_map           = make(map[string]RequestAccess)
	}

	user_map[req_email]  = request_obj
	requests[userID]     = user_map

	err = stub.PutState(REQUEST_ACCESS, t.toJsonRequests(requests))
	if err != nil {
		return shim.Error(err.Error())
	}
	//upgrade state of request_access

	return shim.Success(nil)

}

func (t *ABLChaincode) toJsonPermission(permissions map[string][]string) []byte {
	 permissions_bts, _ := json.Marshal(permissions)
	 return permissions_bts
}

func (t *ABLChaincode) toJsonRequests(requests map[string]map[string]RequestAccess) []byte {
	 requests_bts, _ := json.Marshal(requests)
	 return requests_bts
}

func (t *ABLChaincode) fromJsonPermission(permissions_bts []byte) map[string][]string {
	var permissions map[string][]string
	json.Unmarshal(permissions_bts, &permissions)
	return permissions
}

func (t *ABLChaincode) fromJsonRequests(requests_bts []byte) map[string]map[string]RequestAccess {
	var requests map[string]map[string]RequestAccess
	json.Unmarshal(requests_bts, &requests)
	return requests
}

func (t *ABLChaincode) fromBase64ToArray(dnas_base64 string) []string {
   dnas_bts, _ := base64.URLEncoding.DecodeString(dnas_base64)
	 var dnas []string
	 json.Unmarshal(dnas_bts, &dnas)
 	 return dnas
}

func (t *ABLChaincode) fromBase64ToRequestObj(requests_base64 string)  RequestAccess {
   requests_bts, _ := base64.URLEncoding.DecodeString(requests_base64)
	 var requests RequestAccess
	 json.Unmarshal(requests_bts, &requests)
 	 return requests
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ABLChaincode))
	if err != nil {
		log.Printf("Error starting ABL Service chaincode: %s", err)
	}
}
