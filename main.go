/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	itemContract := new(ItemContract)
	itemContract.Info.Version = "0.0.1"
	itemContract.Info.Description = "My Smart Contract"
	itemContract.Info.License = new(metadata.LicenseMetadata)
	itemContract.Info.License.Name = "Apache-2.0"
	itemContract.Info.Contact = new(metadata.ContactMetadata)
	itemContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(itemContract)
	chaincode.Info.Title = "item chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from ItemContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
