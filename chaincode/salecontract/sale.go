/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("salecontract")

// SaleContract example simple Chaincode implementation
type SaleContract struct {
	Buyer           string
	Seller          string
	DataHash        string
	SignatureBuyer  string
	SignatureSeller string
	Status          int
}

const (
	PROPOSED = iota
	ACCEPTED
	REJECTED
)

func (t *SaleContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### sale contract  Init ###########")

	var err error
	_, args := stub.GetFunctionAndParameters()

	var contract SaleContract
	var jsonContract = args[0]
	err = json.Unmarshal([]byte(jsonContract), &contract)
	if err != nil {
		logger.Error("Could not fetch sale contract from ledger", err)
		return shim.Error("Cannot unmarshal contract values")
	}

	// Initialize the chaincode
	logger.Info(args[0])
	if contract.Buyer == "" {
		return shim.Error("Expecting buyer for a sale contract")
	}
	if contract.Seller == "" {
		return shim.Error("Expecting seller for a sale contract")
	}

	contract.Status = PROPOSED

	logger.Info("buyer = %d, seller = %d,dataHash = %d, status = %d\n", contract.Buyer, contract.Seller, contract.DataHash, contract.Status)

	//// Write the state to the ledger
	var contractJson, marsErr = json.Marshal(contract)
	if marsErr != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contract.Buyer, contractJson)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contract.Seller, contractJson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *SaleContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "accept" {
		// Deletes an entity from its state
		return t.accept(stub, args)
	}
	//
	//if function == "reject" {
	//	// queries an entity state
	//	return t.reject(stub, args)
	//}

	logger.Errorf("Unknown action, check the first argument, must be one of 'accept', 'reject'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'accept', 'delete'. But got: %v", args[0]))
}

func (t *SaleContract) accept(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var buyer, seller string // Entities
	var Aval, Bval int       // Asset holdings
	var X int                // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	buyer = args[0]
	seller = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(buyer)
	if err != nil {
		return shim.Error("Failed to get state of contract")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(seller)
	if err != nil {
		return shim.Error("Failed to get state of contract")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}

	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	logger.Infof("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(buyer, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(seller, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {

	err := shim.Start(new(SaleContract))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
