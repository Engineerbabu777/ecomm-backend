package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func HashPassword() string {

}


func VerifyPassword(userPassword string, givenPassword string)(bool,string){

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

}

func SearchProductByQuery() gin.HandlerFunc{}