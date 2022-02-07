package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ProductID   int                `json:"productId"`
	ProductName string             `json:"productName"`
	Price       int                `json:"price"`
	Tags        []string           `json:"tags"`
}
