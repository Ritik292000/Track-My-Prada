package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Product struct {
	AssetId     string `json:"AssetId"`
	GSTIN       string `json:"GSTIN"`
	Quantity    string `json:"Quantity"`
	InvoiceDate string `json:"InvoiceDate"`
	Weight      string `json:"Weight"`
	Status      string `json:"Status"`
	Destination string `json:"Destination"`
}

type Order struct {
	Quantity    string `json:"Quantity"`
	Status      string `json:"Status"`
	DateOfOrder string `json:"DateOfOrder"`
	DealerName  string `json:"DealerName"`
	OrderID     string `json:"OrderID"`
}

type ProductData struct {
	AssetType   string `json:"AssetType"`
	AssetId     string `json:"AssetId"`
	GSTIN       string `json:"GSTIN"`
	Quantity    string `json:"Quantity"`
	InvoiceDate string `json:"InvoiceDate"`
	Weight      string `json:"Weight"`
	Status      string `json:"Status"`
	Destination string `json:"Destination"`
}

type OrderData struct {
	AssetType   string `json:"AssetType"`
	Quantity    string `json:"Quantity"`
	Status      string `json:"Status"`
	DateOfOrder string `json:"DateOfOrder"`
	DealerName  string `json:"DealerName"`
	OrderID     string `json:"OrderID"`
}

type Match struct {
	OrderId string `json:"orderId"`
	AssetId string `json:"AssetId"`
}

type ProdHistory struct {
	Record    *ProductData `json:"record"`
	TxId      string       `json:"txId"`
	Timestamp string       `json:"timestamp"`
	IsDelete  bool         `json:"isDelete"`
}

type Register struct {
	AssetId    string `json:"AssetId"`
	AssetOwner string `json:"AssetOwner"`
	RegNumber  string `json:"regNumber"`
}

func main() {
	router := gin.Default()

	var wg sync.WaitGroup
	wg.Add(1)
	go ChaincodeEventListener("org1", "pradachannel", "PRADA_CHAINCODE", &wg)

	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetAllAssets")

		var Assets []ProductData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the Assets slice
			if err := json.Unmarshal([]byte(result), &Assets); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		// ctx.JSON(http.StatusOK,Assets)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Prada Application", "AssetList": Assets,
		})
	})

	router.GET("/org1", func(ctx *gin.Context) {
		result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetAllAssets")

		var Assets []ProductData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the Assets slice
			if err := json.Unmarshal([]byte(result), &Assets); err != nil {
				fmt.Println("Error", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "org1.html", gin.H{
			"title": "Org1 Dashboard", "AssetList": Assets,
		})
	})

	router.POST("/api/prada", func(ctx *gin.Context) {
		var req Product
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request prabhat"})
			return
		}

		fmt.Printf("Asset response %s", req)
		submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "invoke", make(map[string][]byte), "CreateAsset", req.AssetId, req.GSTIN, req.Quantity, req.InvoiceDate, req.Weight, req.Status, req.Destination)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/prada/:id", func(ctx *gin.Context) {
		AssetId := ctx.Param("id")

		result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "ReadAsset", AssetId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/order/match-Asset", func(ctx *gin.Context) {
		AssetID := ctx.Query("AssetId")
		result := submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetMatchingOrders", AssetID)

		// fmt.Printf("result %s", result)

		var orders []OrderData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &orders); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "matchOrder.html", gin.H{
			"title": "Matching Orders", "orderList": orders, "AssetId": AssetID,
		})
	})

	router.POST("/api/prada/match-order", func(ctx *gin.Context) {
		var req Match
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("match  %s", req)
		submitTxnFn("org1", "pradachannel", "PRADA_CHAINCODE", "MyContract", "invoke", make(map[string][]byte), "MatchOrder", req.AssetId, req.OrderId)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/event", func(ctx *gin.Context) {
		result := getEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"AssetEvent": result})

	})

	router.GET("/org2", func(ctx *gin.Context) {
		result := submitTxnFn("org2", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")

		var orders []OrderData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &orders); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		// ctx.JSON(http.StatusOK,orders)
		ctx.HTML(http.StatusOK, "org2.html", gin.H{
			"title": "Org2 Dashboard", "orderList": orders,
		})
	})

	//Get all orders
	router.GET("/api/order/all", func(ctx *gin.Context) {

		result := submitTxnFn("org2", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")

		var orders []OrderData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &orders); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		// ctx.JSON(http.StatusOK, gin.H{"orderList": orders})
		ctx.HTML(http.StatusOK, "orders.html", gin.H{
			"title": "All Orders", "orderList": orders,
		})
	})

	router.POST("/api/order", func(ctx *gin.Context) {
		var req Order
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("order  %s", req)

		privateData := map[string][]byte{
			"Quantity":    []byte(req.Quantity),
			"Status":      []byte(req.Status),
			"DateOfOrder": []byte(req.DateOfOrder),
			"DealerName":  []byte(req.DealerName),
		}

		submitTxnFn("org2", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "private", privateData, "CreateOrder", req.OrderID)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/order/:id", func(ctx *gin.Context) {
		orderId := ctx.Param("id")

		result := submitTxnFn("org2", "pradachannel", "PRADA_CHAINCODE", "OrderContract", "query", make(map[string][]byte), "ReadOrder", orderId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/org3", func(ctx *gin.Context) {
		result := submitTxnFn("org3", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetAllAssets")

		var Assets []ProductData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the Assets slice
			if err := json.Unmarshal([]byte(result), &Assets); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "org3.html", gin.H{
			"title": "Org3 Dashboard", "AssetList": Assets,
		})
	})

	router.GET("/api/prada/history", func(ctx *gin.Context) {
		AssetID := ctx.Query("AssetId")
		result := submitTxnFn("org3", "pradachannel", "PRADA_CHAINCODE", "MyContract", "query", make(map[string][]byte), "GetAssetHistory", AssetID)

		// fmt.Printf("result %s", result)

		var Assets []ProdHistory

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &Assets); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "history.html", gin.H{
			"title": "Asset History", "itemList": Assets,
		})
	})

	router.POST("/api/prada/register", func(ctx *gin.Context) {
		var req Register
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("Asset response %s", req)
		submitTxnFn("org3", "pradachannel", "PRADA_CHAINCODE", "MyContract", "invoke", make(map[string][]byte), "RegisterAsset", req.AssetId, req.AssetOwner, req.RegNumber)

		ctx.JSON(http.StatusOK, req)
	})

	router.Run("localhost:8080")
}
