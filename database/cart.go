

package database;


var (
	ErrCartFindProduct = errors.New("can't find the product")
	ErrCartDecodeProducts = errors.New("can't decode the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cant add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cant remove this item from the cart")
	ErrCantGetItem = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the purchase")
)


func AddProductToCart(){

}

func RemoveCartItem(){

}



func ButItemFromCart(){

}


func InstantBuyer(){

}

