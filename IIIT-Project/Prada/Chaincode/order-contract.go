package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// OrderContract contract for managing CRUD for Order
type OrderContract struct {
	contractapi.Contract
}

type Order struct {
	AssetType   string `json:"AssetType"`
	Quantity    string `json:"Quantity"`
	Status      string `json:"Status"`
	DateOfOrder string `json:"DateOfOrder"`
	DealerName  string `json:"DealerName"`
	OrderID     string `json:"OrderID"`
}

func getCollectionName() string {
	collectionName := "OrderCollection"
	return collectionName
}

// OrderExists returns true when asset with given ID exists in private data collection
func (o *OrderContract) OrderExists(ctx contractapi.TransactionContextInterface, orderID string) (bool, error) {
	collectionName := getCollectionName()

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, orderID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateOrder creates a new instance of Order
func (o *OrderContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "Org2MSP" {
		exists, err := o.OrderExists(ctx, orderID)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("The product %s already exists", orderID)
		}

		order := new(Order)

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", err
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("Please provide the private data ")
		}

		quantity, exists := transientData["Quantity"]
		if !exists {
			return "", fmt.Errorf("The make was not specified in transient data. Please try again")
		}
		order.Quantity = string(quantity)

		status, exists := transientData["Status"]
		if !exists {
			return "", fmt.Errorf("The model was not specified in transient data. Please try again")
		}
		order.Status = string(status)

		dateOfOrder, exists := transientData["DateOfOrder"]
		if !exists {
			return "", fmt.Errorf("The color was not specified in transient data. Please try again")
		}
		order.DateOfOrder = string(dateOfOrder)

		dealerName, exists := transientData["DealerName"]
		if !exists {
			return "", fmt.Errorf("The dealer was not specified in transient data. Please try again")
		}
		order.DealerName = string(dealerName)

		order.AssetType = "Order"
		order.OrderID = orderID

		bytes, _ := json.Marshal(order)

		collectionName := getCollectionName()

		return fmt.Sprintf("Order with id %v added successfully", orderID), ctx.GetStub().PutPrivateData(collectionName, orderID, bytes)
	} else {
		return fmt.Sprintf("Order cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

// ReadOrder retrieves an instance of Order from the private data collection
func (o *OrderContract) ReadOrder(ctx contractapi.TransactionContextInterface, orderID string) (*Order, error) {
	exists, err := o.OrderExists(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", orderID)
	}

	return ReadPrivateState(ctx, orderID)

}

func ReadPrivateState(ctx contractapi.TransactionContextInterface, orderID string) (*Order, error) {
	collectionName := getCollectionName()

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, orderID)
	if err != nil {
		return nil, err
	}
	order := new(Order)

	err = json.Unmarshal(bytes, order)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal private data collection data to type Order")
	}

	return order, nil
}

// DeleteOrder deletes an instance of Order from the private data collection
func (o *OrderContract) DeleteOrder(ctx contractapi.TransactionContextInterface, orderID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientOrgID == "Org1MSP" || clientOrgID == "Org2MSP" {
		exists, err := o.OrderExists(ctx, orderID)

		if err != nil {
			return fmt.Errorf("Could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("The asset %s does not exist", orderID)
		}

		collectionName := getCollectionName()

		return ctx.GetStub().DelPrivateData(collectionName, orderID)
	} else {
		return fmt.Errorf("Organisation with %v cannot delete the order", clientOrgID)
	}
}

// GetAllOrders retrieves all the asset with assetype 'Order'

//	func (o *OrderContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
//		collectionName := getCollectionName()
//		queryString := `{"selector":{"AssetType":"Order"}}`
//		resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
//		if err != nil {
//			return nil, err
//		}
//		defer resultsIterator.Close()
//		return orderResultIteratorFunction(resultsIterator)
//	}
func (o *OrderContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
	collectionName := getCollectionName()
	queryString := `{"selector":{"AssetType":"Order"}}`

	// Log collection name and query string
	fmt.Println("Collection Name:", collectionName)
	fmt.Println("Query String:", queryString)

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		fmt.Println("Error getting query result:", err)
		return nil, err
	}
	defer resultsIterator.Close()

	// Log that the query result iterator was obtained successfully
	fmt.Println(resultsIterator)

	return orderResultIteratorFunction(resultsIterator)
}

// GetOrdersByRange gives a range of order details based on a start key and an end key
func (o *OrderContract) GetOrdersByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*Order, error) {
	collectionName := getCollectionName()
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return orderResultIteratorFunction(resultsIterator)

}

// iterator function

func orderResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Order, error) {
	var orders []*Order
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var order Order
		err = json.Unmarshal(queryResult.Value, &order)
		if err != nil {
			return nil, err

		}
		orders = append(orders, &order)
	}

	return orders, nil
}
