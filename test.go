package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// get the args from the transaction proposal
	args := stub.GetStringArgs()

	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Excepeting a key and a value")
	}
	// Set up any variables or assets here by claling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(arg[0], []byte(args[1]))

	// if error occurs
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}

	// success
	return shim.Success(nil)
}

// First, let’s add the Invoke function’s signature.

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Extract the function and args from the transaction proposal

	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	if fn == "set" {
		result, err = set(stub, args)
	} else {
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments, Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))

	if err != nil {
		return "", fmt.Errorf("Failed to set asset : %s", args[0])
	}
	return args[1], nil
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments, Expecting a key and a value")
	}

	value, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error : %s ", args[0], err)
	}

	if value == nil {
		return "", fmt.Errorf("Asset not founc : %s", args[0])
	}

	return string(value), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(SimpleAsset))
	if err != nil {
		fmt.Printf("Error Starting SimpleAsset chaincode : %s", err)
	} else {
		fmt.Printf("Starting chaincode function main() executed successfully")
	}
}
