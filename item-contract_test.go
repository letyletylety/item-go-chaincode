// /*
//  * SPDX-License-Identifier: Apache-2.0
//  */

package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/hyperledger/fabric-chaincode-go/shim"
// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// const getStateError = "world state get error"

// type MockStub struct {
// 	shim.ChaincodeStubInterface
// 	mock.Mock
// }

// func (ms *MockStub) GetState(key string) ([]byte, error) {
// 	args := ms.Called(key)

// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (ms *MockStub) PutState(key string, value []byte) error {
// 	args := ms.Called(key, value)

// 	return args.Error(0)
// }

// func (ms *MockStub) DelState(key string) error {
// 	args := ms.Called(key)

// 	return args.Error(0)
// }

// type MockContext struct {
// 	contractapi.TransactionContextInterface
// 	mock.Mock
// }

// func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
// 	args := mc.Called()

// 	return args.Get(0).(*MockStub)
// }

// func configureStub() (*MockContext, *MockStub) {
// 	var nilBytes []byte

// 	testItem := new(Item)
// 	testItem.Value = "set value"
// 	itemBytes, _ := json.Marshal(testItem)

// 	ms := new(MockStub)
// 	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
// 	ms.On("GetState", "missingkey").Return(nilBytes, nil)
// 	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
// 	ms.On("GetState", "itemkey").Return(itemBytes, nil)
// 	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
// 	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

// 	mc := new(MockContext)
// 	mc.On("GetStub").Return(ms)

// 	return mc, ms
// }

// func TestItemExists(t *testing.T) {
// 	var exists bool
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(ItemContract)

// 	exists, err = c.ItemExists(ctx, "statebad")
// 	assert.EqualError(t, err, getStateError)
// 	assert.False(t, exists, "should return false on error")

// 	exists, err = c.ItemExists(ctx, "missingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
// 	assert.False(t, exists, "should return false when no value for key in world state")

// 	exists, err = c.ItemExists(ctx, "existingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
// 	assert.True(t, exists, "should return true when value for key in world state")
// }

// func TestCreateItem(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(ItemContract)

// 	err = c.CreateItem(ctx, "statebad", "some value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.CreateItem(ctx, "existingkey", "some value")
// 	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

// 	err = c.CreateItem(ctx, "missingkey", "some value")
// 	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
// }

// func TestReadItem(t *testing.T) {
// 	var item *Item
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(ItemContract)

// 	item, err = c.ReadItem(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
// 	assert.Nil(t, item, "should not return Item when exists errors when reading")

// 	item, err = c.ReadItem(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
// 	assert.Nil(t, item, "should not return Item when key does not exist in world state when reading")

// 	item, err = c.ReadItem(ctx, "existingkey")
// 	assert.EqualError(t, err, "Could not unmarshal world state data to type Item", "should error when data in key is not Item")
// 	assert.Nil(t, item, "should not return Item when data in key is not of type Item")

// 	item, err = c.ReadItem(ctx, "itemkey")
// 	expectedItem := new(Item)
// 	expectedItem.Value = "set value"
// 	assert.Nil(t, err, "should not return error when Item exists in world state when reading")
// 	assert.Equal(t, expectedItem, item, "should return deserialized Item from world state")
// }

// func TestUpdateItem(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(ItemContract)

// 	err = c.UpdateItem(ctx, "statebad", "new value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

// 	err = c.UpdateItem(ctx, "missingkey", "new value")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

// 	err = c.UpdateItem(ctx, "itemkey", "new value")
// 	expectedItem := new(Item)
// 	expectedItem.Value = "new value"
// 	expectedItemBytes, _ := json.Marshal(expectedItem)
// 	assert.Nil(t, err, "should not return error when Item exists in world state when updating")
// 	stub.AssertCalled(t, "PutState", "itemkey", expectedItemBytes)
// }

// func TestDeleteItem(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(ItemContract)

// 	err = c.DeleteItem(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.DeleteItem(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

// 	err = c.DeleteItem(ctx, "itemkey")
// 	assert.Nil(t, err, "should not return error when Item exists in world state when deleting")
// 	stub.AssertCalled(t, "DelState", "itemkey")
// }
