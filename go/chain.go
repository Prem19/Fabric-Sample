
/*
  Sample Chaincode to record and query ledger for POC
 */

package main

 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}


type Data struct {
	Name string `json:"name"`
	Date string `json:"date"`
	Enabled  string `json:"enabled"`
	Status  string `json:"status"`
	Reporter  string `json:"reporter"`
}

/*
 * The Init method *
 Called when the Smart Contract is instantiated by the network
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 Called when an application requests to run the Smart Contract 
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	if function == "record" {
		return s.record(APIstub, args)
	} else if function == "queryAll" {
		return s.queryAll(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}


/*
 * The record method *
 */
func (s *SmartContract) record(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments")
	}

	var data = Data{ Name: args[1], Date: args[2], Enabled: args[3], Status: args[4], Reporter: args[5] }

	dataAsBytes, _ := json.Marshal(data)
	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record : %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAll data method *
 */
func (s *SmartContract) queryAll(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}