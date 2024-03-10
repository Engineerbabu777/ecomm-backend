package controllers

import (
	"context"
	"ecomm/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func AddAddress() gin.HandlerFunc{
	
}

func EditAddress() gin.HandlerFunc{
	
}


func EditWorkAddress() gin.HandlerFunc{
	
}

func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context) {
	    userID := c.Query("id");

		if userID == "" {
			log.Println("userID is empty!");
			c.Header("Content-Type", "application/json");
			c.JSON(400, gin.H{"error": "userID is empty!"});
			c.Abort();
			return;
		}

		addresses := make([]models.Address,0);
		userID,err := primitive.ObjectIDFromHex(userID);

		if err != nil {
			c.IndentedJSON(500,"Internal server error!");
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);
		defer cancel();
		
		filter := bson.D{primitive.E{Key: "_id", Value: userID}};

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}}
		_,err = UserCollection.UpdateOne(ctx,filter,update);

		if err != nil {
			c.IndentedJSON(404,"Wrong command");
			return;
		}
		defer cancel();
		ctx.Done();
		c.IndentedJSON(200, "successfully deleted address");

	}
}