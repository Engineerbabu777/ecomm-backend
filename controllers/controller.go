package controllers

import (
	"context"
	"ecomm/database"
	"ecomm/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


var UserCollection *mongo.Collection= database.UserData(database.Client, "Users");
var ProductCollection *mongo.Collection= database.ProductData(database.Client, "Products");

var Validate = validator.New();


func HashPassword(password string) string {
bytes,err := bcrypt.GenerateFromPassword([]byte(password),14);
if err != nil {
	log.Panic(err);
}
return string(bytes);
}

func VerifyPassword(userPassword string, givenPassword string)(bool,string){

	err := bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword));

	check := true;
	msg := "";

	if err != nil {
		msg = "Login or Password is incorrect!";
		check = false;
	}

	return check,msg;
}

func Signup() gin.HandlerFunc{

	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);
		defer cancel();

		var user models.User;
		if err:=c.BindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return;
		}

		validateError := validate.Struct(user);
		if validateError != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validateError});
			return;
		}

		count,err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email});
		if err != nil{
			log.Panic(err);
			c.JSON(http.StatusInternalServerError, gin.H{"error": err});
		}
		if count >0 {
			c.JSON(http.StatusBadRequest, gin.H{"err":"user already exists"});
		}

		count,err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone});
		if err != nil{
			log.Panic(err);
			c.JSON(http.StatusInternalServerError, gin.H{"error": err});
		}

		if count >0 {
			c.JSON(http.StatusBadRequest, gin.H{"err":"phone number already exists"});
			return;
		}

		password := HashPassword(user.Password);

		user.Password = password;

		user.Created_At,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
		user.Updated_At,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));

		user.ID = primitive.NewObjectID();
		user.User_ID = user.ID.String();

		token, refreshToken := generate.TokenGenerator(user.Email,user.First_Name, user.Last_Name, user.User_ID);

		user.Token= token ;
		user.Refresh_Token = refreshToken;
		user.UserCart = make([]models.ProductUser,0);
		user.AddressDetails= make([]models.Address,0);
		user.Order_Status = make([]models.OrderStatus,0);

		_,inserterr := UserCollection.InsertOne(ctx,user);

		if inserterr!=nil{
          c.JSON(http.StatusInternalServerError, gin.H{"error": inserterr})
		  return;
		}
		defer cancel();
		c.JSON(http.StatusCreated,gin.H{"message":"Successfully account created!","user":user});
	}
}

func Login() gin.HandlerFunc{
    return func(c *gin.Context){
		ctx,cancel := context.WithTimeoutOut(context.Background(),time.Second*100);

		defer cancel();

		var user models.User;
		err := c.BindJSON(&user);
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err});
		}

		var foundUser models.User;
		err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser);

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err});
		}

		PasswordIsValid,msg := VerifyPassword(user.Password, foundUser.password);

		defer cancel();
		if(!PasswordIsValid){
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg});
			return;
		}

		token,refreshToken ,_ := generate.TokenGenerator(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_ID);

		generate.UpdateAllTokens(token, refreshToken, foundUser.User_Id);

		c.JSON(http.StatusFound,foundUser)

	}
}

func ProductViewerAdmin() gin.HandlerFunc{
}

func SearchProduct() gin.HandlerFunc{

	return func(ctx *gin.Context){

		var ProductList []models.Product;
		var c,cancel = context.WithTimeout(context.Background(),100*time.Second);

		defer cancel();

		cursor,err := ProductCollection.Find(c,bson.M{});

		if err!= nil {
			ctx.IndentedJSON(http.StatusInternalServerError,"something went wrong, please try again later!");
			return;
		}

		err = cursor.All(c,&ProductList);

		if err != nil {
			log.Println(err);
			ctx.AbortWithStatus(http.StatusInternalServerError);
			return;
		}

		defer cursor.Close(c);

		if err := cursor.Err(); err!= nil{
			log.Println(err);
			ctx.IndentedJSON(400,"invalid");
			return;
		}
		defer cancel();
		ctx.IndentedJSON(200,ProductList); 
	}
}

func SearchProductByQuery() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		// CREATING A NEW VAR THAT WILL HOLDS ALL THE PRODUCTS!
		var searchProducts []models.Product;
	// GET THE NAME FROM THE QUERY AND STORED!!
		queryParam := ctx.Query("name");

		// IF YOU DO NOT HAVE QUERY!
		if queryParam == "" {
		    log.Println("query is empty!");
			ctx.Header("Content-Type", "application/json");
			ctx.JSON(http.StatusNotFound,gin.H{"Error": "query is empty!"});
			ctx.Abort();
			return;
		}

		// CREATING A CONTEXT TIME OUT!
		var c,cancel = context.WithTimeout(context.Background(),100*time.Second);

		defer cancel();

		// FINDING THE DATA FROM THE DATABASE USING REGEX!
		searchquerydb,err := ProductCollection.Find(c, bson.M{"product_name":bson.M{"$regex":queryParam}});

		if err!= nil {
		    ctx.IndentedJSON(404,"something went wrong while fetching the data!");
			return;
		}

		// TRANSLATING DATA TO GO STRUCT!
		err = searchquerydb.All(c,&searchProducts);


		if err != nil {
		    log.Println(err);
			ctx.IndentedJSON(400,"invalid!");
			return;
		}

		// WHEN ALL THE DATA IS TRANSFERED CLOSE THE CURSOR!
		defer searchquerydb.Close(c);


		if err := searchquerydb.Err(); err!= nil{
		    log.Println(err);
			ctx.IndentedJSON(400,"invalid request!");
			return;
		}
		defer cancel();

		// RETURNED DATA BACK!
		ctx.IndentedJSON(200, searchProducts);
	}
}