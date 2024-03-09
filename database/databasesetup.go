package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func DBSet()*mongo.Client{


	client ,err := mongo.NewClient(options.Client().ApplyURI(""));

	if err != nil{
		log.Fatal(err);
	}

	ctx,cancel := context.WithTimeout(context.Background(), time.Second*10)
 
	defer cancel();
   err = client.Connect(ctx);

	if err != nil{
		log.Fatal(err);
	}

	err = client.Ping(context.TODO(), nil);

	if err != nil{
		log.Println("Failed yo connect to mongodb");
		return nil;
	}

	return client;
}

var Client *mongo.Client = DBSet();


func UserData(client *mongo.Client, collectionName string) *mongo.Collection{

	var collection *mongo.Collection = client.Database("Ecomm").Collection(collectionName);

	return collection;
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection{
	
	var collection *mongo.Collection = client.Database("Ecomm").Collection(collectionName);

	return collection;
}