package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MyContract struct {
	contractapi.Contract
}

type Product struct {
	AssetType   string `json:"AssetType"`
	AssetId     string `json:"AssetId"`
	GSTIN       string `json:"GSTIN"`
	Quantity    string `json:"Quantity"`
	InvoiceDate string `json:"InvoiceDate"`
	Weight      string `json:"Weight"`
	Status      string `json:"Status"`
	Destination string `json:"Destination"`
}

type EventData struct {
	AssetType string
	GSTIN     string
}

type PaginatedQueryResult struct {
	Records             []*Product `json:"records"`
	FetchedRecordsCount int32      `json:"fetchedRecordsCount"`
	Bookmark            string     `json:"bookmark"`
}

type HistoryQueryResult struct {
	Record    *Product `json:"record"`
	TxId      string   `json:"txId"`
	Timestamp string   `json:"timestamp"`
	IsDelete  bool     `json:"isDelete"`
}

func (c *MyContract) ProductExists(ctx contractapi.TransactionContextInterface, AssetId string) (bool, error) {
	data, err := ctx.GetStub().GetState(AssetId)
	if err != nil {
		return false, fmt.Errorf("failed to read data:%v", err)
	}
	return data != nil, nil
}

func (c *MyContract) CreateAsset(ctx contractapi.TransactionContextInterface, AssetId string, GSTIN string, Quantity string, InvoiceDate string, Weight string, Status string, Destination string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "Org1MSP" {

		exists, err := c.ProductExists(ctx, AssetId)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the product %s already exists", AssetId)
		}

		prada := Product{

			AssetType:   "Package",
			AssetId:     AssetId,
			GSTIN:       GSTIN,
			Quantity:    Quantity,
			InvoiceDate: InvoiceDate,
			Weight:      Weight,
			Status:      Status,
			Destination: Destination,
		}

		bytes, _ := json.Marshal(prada)

		err = ctx.GetStub().PutState(AssetId, bytes)
		if err != nil {
			return "", err
		} else {
			addAssetEventData := EventData{
				AssetType: "Package",
				GSTIN:     GSTIN,
			}
			eventDataByte, _ := json.Marshal(addAssetEventData)
			ctx.GetStub().SetEvent("CreateAsset", eventDataByte)
			return fmt.Sprintf("successfully added Product %v", AssetId), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}

}

func (c *MyContract) ReadAsset(ctx contractapi.TransactionContextInterface, AssetId string) (*Product, error) {
	exists, err := c.ProductExists(ctx, AssetId)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the product %s does not exist", AssetId)
	}

	bytes, _ := ctx.GetStub().GetState(AssetId)

	prada := new(Product)

	err = json.Unmarshal(bytes, &prada)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Product")
	}

	return prada, nil
}

// DeleteCar removes the instance of Asset from the world state
func (c *MyContract) DeleteAsset(ctx contractapi.TransactionContextInterface, AssetId string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "Org1MSP" {

		exists, err := c.ProductExists(ctx, AssetId)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		} else if !exists {
			return "", fmt.Errorf("The Product %s does not exist", AssetId)
		}

		err = ctx.GetStub().DelState(AssetId)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Product with id %v is deleted from the world state.", AssetId), nil
		}

	} else {
		return "", fmt.Errorf("User under following MSP:%v cannot able to perform this action", clientOrgID)
	}
}

// GetAllCars retrieves all the asset with assetype 'car'
func (c *MyContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Product, error) {

	queryString := `{
        "selector": {
            "AssetType": "Package"
        }, 
		"sort":[{ "AssetId": "desc"}]
    }`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return assetResultIteratorFunction(resultsIterator)
}

// Iterator function
func assetResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Product, error) {
	var assets []*Product
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Product
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (c *MyContract) GetAssetByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return assetResultIteratorFunction(resultsIterator)
}

// GetCarsWithPagination gives a set of car details based on a page size,  and a bookmark.
func (c *MyContract) GetAssetWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	queryString := `{"selector":{"AssetType":"Package"}}`
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets, err := assetResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetCarHistory returns the history of a car since issuance.
func (c *MyContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, AssetID string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(AssetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Product
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Product{
				AssetId: AssetID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}
