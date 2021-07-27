/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"time"
)

// Item stores a value
type Item struct {
	DocType  string    `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID       string    `json:"ID"`      //the field tags are needed to keep case from bouncing around
	Name     string    `json:"name"`
	Content  string    `json:"content"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
	Seller   string    `json:"seller"`
	Regdate  time.Time `json:"regdate"`
}

func NewItem(
	id string,
	name, content string,
	price, quantity int,
	seller string,
	regdate time.Time) *Item {

	item := new(Item)
	item.DocType = "item"
	item.ID = id
	item.Name = name
	item.Content = content
	item.Price = price
	item.Quantity = quantity
	item.Seller = seller
	return item
	// &Item{Name: name, Content: content, Price: price, Quantity: quantity, Seller: seller, Regdate: regdate}
}
