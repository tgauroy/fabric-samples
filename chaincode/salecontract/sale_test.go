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
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInitFailed(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status == shim.OK {
		fmt.Println("Init sucess but failed expected", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkStateNotExist(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes != nil {
		fmt.Println("State", name, "have value")
		t.FailNow()
	}

}

func checkAcceptFailed(t *testing.T, stub *shim.MockStub, name string, actor string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("accept"), []byte(name), []byte(actor)})
	if res.Status == shim.OK {
		fmt.Println("Accept", name, "error accept should failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		fmt.Println("Accept", name, "Payload is not null")
		t.FailNow()
	}
}

func checkAccept(t *testing.T, stub *shim.MockStub, name string, actor string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("accept"), []byte(name), []byte(actor)})
	if res.Status != shim.OK {
		fmt.Println("Accept", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Accept", name, "failed to get contract")
		t.FailNow()
	}

	var contract SaleContract
	var err = json.Unmarshal([]byte(res.Payload), &contract)
	if err != nil {
		logger.Error("Could not fetch sale contract from payload", err)
		t.FailNow()
	}

	if contract.Status != ACCEPTED {
		fmt.Println("Contract status value", name, "was not ACCEPTED as expected")
		t.FailNow()
	}
}

func checkRejectFailed(t *testing.T, stub *shim.MockStub, name string, actor string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("reject"), []byte(name), []byte(actor)})
	if res.Status == shim.OK {
		fmt.Println("Reject", name, "error accept should failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		fmt.Println("Reject", name, "Payload is not null")
		t.FailNow()
	}
}

func checkReject(t *testing.T, stub *shim.MockStub, name string, actor string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("reject"), []byte(name), []byte(actor)})
	if res.Status != shim.OK {
		fmt.Println("Reject", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Reject", name, "failed to get contract")
		t.FailNow()
	}

	var contract SaleContract
	var err = json.Unmarshal([]byte(res.Payload), &contract)
	if err != nil {
		logger.Error("Could not fetch sale contract from payload", err)
		t.FailNow()
	}

	if contract.Status != REJECTED {
		fmt.Println("Contract status value", name, "was not ACCEPTED as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestSaleContract_Init(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-001",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          PROPOSED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(string(totoStr))

	checkInit(t, stub, [][]byte{[]byte("init"), totoStr})
	checkState(t, stub, "SALE-001", string(totoStr))
}

func Test_SaleContract_Init_Proposed_Status(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-002",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          ACCEPTED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(string(totoStr))

	checkInitFailed(t, stub, [][]byte{[]byte("init"), totoStr})
	checkStateNotExist(t, stub, "SALE-002", string(totoStr))
}

func Test_Buyer_accept_contract_and_status_accepted(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-003",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          PROPOSED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(string(totoStr))

	checkInit(t, stub, [][]byte{[]byte("init"), totoStr})
	checkState(t, stub, "SALE-003", string(totoStr))
	checkAccept(t, stub, "SALE-003", "Acheteur")

	toto.Status = ACCEPTED

	totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	checkState(t, stub, "SALE-003", string(totoStr))

}

func Test_Seller_cannot_accept_contract(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-004",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          PROPOSED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(string(totoStr))

	checkInit(t, stub, [][]byte{[]byte("init"), totoStr})
	checkAcceptFailed(t, stub, "SALE-004", "Vendeur")
	checkState(t, stub, "SALE-004", string(totoStr))

}

func Test_Buyer_reject_contract_and_status_accepted(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-005",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          PROPOSED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(string(totoStr))

	checkInit(t, stub, [][]byte{[]byte("init"), totoStr})
	checkState(t, stub, "SALE-005", string(totoStr))
	checkReject(t, stub, "SALE-005", "Acheteur")

	toto.Status = REJECTED

	totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	checkState(t, stub, "SALE-005", string(totoStr))

}

func Test_Seller_cannot_reject_contract(t *testing.T) {
	scc := new(SaleContract)
	stub := shim.NewMockStub("ex02", scc)
	toto := &SaleContract{
		Contract:        "SALE-006",
		Buyer:           "Acheteur",
		Seller:          "Vendeur",
		DataHash:        "Hash",
		SignatureBuyer:  "sgn1",
		SignatureSeller: "sgn2",
		Status:          PROPOSED,
	}
	var totoStr, err = json.Marshal(toto)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(string(totoStr))

	checkInit(t, stub, [][]byte{[]byte("init"), totoStr})
	checkRejectFailed(t, stub, "SALE-006", "Vendeur")
	checkState(t, stub, "SALE-006", string(totoStr))

}



