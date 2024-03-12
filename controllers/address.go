package controllers

import (
	"context"
	"ecomm/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			log.Println("User ID is required")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(userID)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		err = c.BindJSON(&addresses)

		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}

		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}

		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error:)")
		}

		var addressInfo []bson.M

		err = pointcursor.All(ctx, &addressInfo)

		if err != nil {
			panic(err)
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
		} else {
			c.IndentedJSON(400, "Not Allowed")
		}

		defer cancel()
		ctx.Done()

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			log.Println("userID is empty!")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error!")
		}
		var editaddress models.Address
		err = c.BindJSON(&editaddress)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house", Value: editaddress.House}, {Key: "address.0.street", Value: editaddress.Street}, {Key: "address.0.city", Value: editaddress.City},
			{Key: "address.0.pin_code", Value: editaddress.Pincode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(500, "Something went wrong!")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "sucessfully updated the home address")

	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			log.Println("userID is empty!")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error!")
		}
		var editaddress models.Address
		err = c.BindJSON(&editaddress)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house", Value: editaddress.House}, {Key: "address.1.street", Value: editaddress.Street}, {Key: "address.1.city", Value: editaddress.City},
			{Key: "address.1.pin_code", Value: editaddress.Pincode},
		}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update);

		if err != nil {
			c.IndentedJSON(500,"something went wrong!");
			return;
		}

		defer cancel();
		ctx.Done();
		c.IndentedJSON(200,"successfully updated the work address!")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			log.Println("userID is empty!")
			c.Header("Content-Type", "application/json")
			c.JSON(400, gin.H{"error": "userID is empty!"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		userID, err := primitive.ObjectIDFromHex(userID)

		if err != nil {
			c.IndentedJSON(500, "Internal server error!")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userID}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "successfully deleted address")

	}
}
