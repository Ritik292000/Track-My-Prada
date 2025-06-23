package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	myContract := new(MyContract)
	orderContarct := new(OrderContract)
	chaincode, err := contractapi.NewChaincode(myContract, orderContarct)

	if err != nil {
		log.Panicf("Could not create chaincode." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode. " + err.Error())
	}
}
