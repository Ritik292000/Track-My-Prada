package main

import "fmt"

func main() {

	submitTxnFn(
		"org1",
		"pradachannel",
		"PRADA_CHAINCODE",
		"MyContract",
		"invoke",
		make(map[string][]byte),
		"CreateAsset",
		"ASSET-05",
		"IN-98765HJ",
		"4",
		"34-11-2023",
		"21KG",
		"deliverd",
		"In warehouse",
	)

	// privateData := map[string][]byte{
	//  "Quantity":       []byte("10Kg"),
	//  "Status":      []byte("Out of Shop"),
	//  "DateOfOrder":      []byte("22-09-2000"),
	//  "DealerName": []byte("Ravi"),
	// }

	submitTxnFn("org2", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "private", make(map[string][]byte), "CreateOrder", "ORD-03")
	// result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "ReadAsset", "ASSET-04")
	result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetAllAssets")
	// submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "query", make(map[string][]byte),"GetAllOrders")
	// submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte),"GetMatchingOrders", "Car-06")
	// submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "invoke", make(map[string][]byte),"MatchOrder", "Ord-06", "ORD-01")
	// submitTxnFn("org3", "pradachannel", "PRADA_CHAINCODE", "MyContract", "invoke", make(map[string][]byte),"RegisterAsset", "Car-06", "Dani", "KL-01-CD-01")
	fmt.Println(result)
}
