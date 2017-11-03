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
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("salecontract")

// SaleContract example simple Chaincode implementation
type SaleContract struct {
	Contract        string
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

	if contract.Status != PROPOSED {
		return shim.Error("Only status proposed to init new contract")
	}

	logger.Info("buyer = %d, seller = %d,dataHash = %d, status = %d\n", contract.Buyer, contract.Seller, contract.DataHash, contract.Status)

	//// Write the state to the ledger
	var contractJson, marsErr = json.Marshal(contract)
	if marsErr != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contract.Contract, contractJson)
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
		logger.Info("Accept invoked")
		// Deletes an entity from its state
		return t.accept(stub, args)
	}

	if function == "reject" {
		// queries an entity state
		return t.reject(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'accept', 'reject'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'accept', 'delete'. But got: %v", args[0]))
}

func (t *SaleContract) accept(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var contractId = args[0]
	var validator = args[1]

	// Get the state from the ledger
	contractbytes, err := stub.GetState(contractId)
	if err != nil {
		return shim.Error("Failed to get state of contract")
	}
	if contractbytes == nil {
		return shim.Error("Contract not found")
	}

	var contract SaleContract
	err = json.Unmarshal([]byte(contractbytes), &contract)
	if err != nil {
		logger.Error("Could not fetch sale contract from ledger", err)
		return shim.Error("Cannot unmarshal contract values")
	}

	if contract.Status != PROPOSED {
		logger.Error("Could accept a contract with a status different than PROPOSED")
		return shim.Error("Could accept a contract with a status different thant PROPOSED")
	}

	if (validator != contract.Buyer) {
		logger.Error("Only Buyer can accept contract")
		return shim.Error("Only Buyer can accept contract")
	}

	contract.Status = ACCEPTED

	// Write the state back to the ledger
	var contractToSave, marsErr = json.Marshal(contract)
	if marsErr != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contract.Contract, []byte(contractToSave))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(contractToSave)
}


func (t *SaleContract) reject(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var contractId = args[0]
	var validator = args[1]

	// Get the state from the ledger
	contractbytes, err := stub.GetState(contractId)
	if err != nil {
		return shim.Error("Failed to get state of contract")
	}
	if contractbytes == nil {
		return shim.Error("Contract not found")
	}

	var contract SaleContract
	err = json.Unmarshal([]byte(contractbytes), &contract)
	if err != nil {
		logger.Error("Could not fetch sale contract from ledger", err)
		return shim.Error("Cannot unmarshal contract values")
	}

	if contract.Status != PROPOSED {
		logger.Error("Could accept a contract with a status different than PROPOSED")
		return shim.Error("Could accept a contract with a status different thant PROPOSED")
	}

	if (validator != contract.Buyer) {
		logger.Error("Only Buyer can reject contract")
		return shim.Error("Only Buyer can reject contract")
	}

	contract.Status = REJECTED

	// Write the state back to the ledger
	var contractToSave, marsErr = json.Marshal(contract)
	if marsErr != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contract.Contract, []byte(contractToSave))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(contractToSave)
}

func main() {

	err := shim.Start(new(SaleContract))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
