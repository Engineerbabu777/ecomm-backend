package controllers

import (
	"context"
	"ecomm/database"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection) *Application {

	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CHECKING FOR THE PRODUCT ID!
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty!")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty!"))
			return
		}
		// CHECKING FOR THE USER ID!
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("User id is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty!"))
			return
		}
		// GETTING ID!
		productId,err := primitive.ObjectIDFromHex(productQueryID);

		if err!=nil{
			log.Println(err);
			c.AbortWithStatus(http.StatusInternalServerError);
			return;
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		defer cancel()

		err = database.AddProductToCart(ctx,app.prodCollection,app.userCollection,productId,userQueryId);

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err);
		}
		c.IndentedJSON(200,"Successfully added product");

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
   return func(ctx *gin.Context) {
	productQueryID := ctx.Query("id");
	if productQueryID == "" {
	   log.Println("product id is empty!");
	   
	   ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty!"));
	   return;
	}

	userQueryID := ctx.Query("userID");
	if userQueryID == "" {
	    log.Println("user id is empty!");

		ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"));
		return;
	}

	productID, err := primitive.ObjectIDFromHex(productQueryID);

	if err!=nil{
		log.Println(err);
		ctx.AbortWithStatus(http.StatusInternalServerError);
		return;
	}

	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second);

	defer cancel();

	err = database.RemoveCartItem(c,app.prodCollection,app.userCollection,productID,userQueryID);

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err);
	}
	ctx.IndentedJSON(200, "product removed from cart!");

   }
}

func GetItemFromCart() gin.HandlerFunc {

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
return func(ctx *gin.Context) {
    
	userQueryID := ctx.Query("id");

	if userQueryID == "" {
		log.Panic("user id is empty");

		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"));

		var c,cancel = context.WithTimeout(context.Background(),time.Second*100);

		defer cancel();

		err := database.ButItemFromCart(c,app.userCollection,userQueryID);

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError,err);
		}

		ctx.IndentedJSON(200,"success fully buyed from cart!")

	}
}

}

func (app *Application) InstantBuy() gin.HandlerFunc {
   return func (c *gin.Context)  {
	 
	// PRODUCT ID!
    productQueryID := c.Query("id");

	if productQueryID== ""{
		log.Println("product id is empty");

		c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty!"));
		return;
	}
	// USER ID
	userQueryID := c.Query("userID");
	if userQueryID== ""{
		log.Println("user id is empty");

		c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"));
		return;
	}

	productID,err := primitive.ObjectIDFromHex(productQueryID);

	if err !=nil{
		log.Println(err);

		c.AbortWithStatus(http.StatusInternalServerError);
		return;
	}

	var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second);

	defer cancel();

	err = database.InstantBuyer(ctx, app.prodCollection,app.userCollection, productID,userQueryID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err);
	}

	c.IndentedJSON(200, "Sucessfully placed the order!")
   }
}
