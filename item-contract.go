/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 여러 아이템 query에 필요한 구조체
type QueryResult struct {
	Key    string `json:"key"`
	Record *Item  `json:"value"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*Item `json:"items"`
	FetchedRecordsCount int32   `json:"fetchedRecordsCount"`
	Bookmark            string  `json:"bookmark"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Item     `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// ItemContract contract for managing CRUD for Item
type ItemContract struct {
	contractapi.Contract
}

// ItemExists returns true when asset with given ID exists in world state
func (c *ItemContract) ItemExists(ctx contractapi.TransactionContextInterface, itemID string) (bool, error) {
	data, err := ctx.GetStub().GetState(itemID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateItem creates a new instance of Item
func (c *ItemContract) CreateItem(ctx contractapi.TransactionContextInterface,
	id, name, content string, price, quantity int, seller string,
) error {
	exists, err := c.ItemExists(ctx, id)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	// 시간
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	timestamp, err := ptypes.Timestamp(txTimestamp)

	if err != nil {
		return err
	}

	item := NewItem(id, name, content, price, quantity, seller, timestamp)

	bytes, _ := json.Marshal(item)

	return ctx.GetStub().PutState(id, bytes)
}

func (c *ItemContract) ChangeItem(ctx contractapi.TransactionContextInterface,
	id, name, content string, price, quantity int, seller string,
) error {
	exists, err := c.ItemExists(ctx, id)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	item := new(Item)
	item.Name = name
	item.Price = price
	item.Quantity = quantity
	item.Seller = seller

	bytes, _ := json.Marshal(item)

	return ctx.GetStub().PutState(id, bytes)
}

// DeleteItem deletes an instance of Item from the world state
func (c *ItemContract) DeleteItem(ctx contractapi.TransactionContextInterface, itemID string) error {
	exists, err := c.ItemExists(ctx, itemID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("the asset %s does not exist", itemID)
	}

	return ctx.GetStub().DelState(itemID)
}

// 물품1개 조회
// ReadItem retrieves an instance of Item from the world state
func (c *ItemContract) ReadItem(ctx contractapi.TransactionContextInterface, itemID string) (*Item, error) {
	exists, err := c.ItemExists(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", itemID)
	}

	bytes, _ := ctx.GetStub().GetState(itemID)

	item := new(Item)

	err = json.Unmarshal(bytes, item)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Item")
	}

	return item, nil
}

// 물품전체 조회
func (c *ItemContract) ReadAllItem(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]QueryResult, error) {
	results := []QueryResult{}

	// The keys are returned by the iterator in lexical order. Note
	// that startKey and endKey can be empty string, which implies unbounded range
	// query on start or end.
	resultIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()

		if err != nil {
			return nil, err
		}

		item := new(Item)
		_ = json.Unmarshal(queryResponse.Value, item)

		queryResult := QueryResult{Key: queryResponse.Key, Record: item}
		results = append(results, queryResult)
	}
	return results, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Item, error) {
	var items []*Item
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var item Item
		err = json.Unmarshal(queryResult.Value, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

// 물품조건 조회
func (c *ItemContract) QueryItems(ctx contractapi.TransactionContextInterface, queryString string) ([]*Item, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Item, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// 물품조건조회 개수 제한
func (c *ItemContract) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*Item, error) {

	resultsIterator, _, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// 변경 이력 조회
// GetAssetHistory returns the chain of custody for an asset since issuance.
func (c *ItemContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var item Item
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &item)
			if err != nil {
				return nil, err
			}
		} else {
			item = Item{
				ID: assetID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &item,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
