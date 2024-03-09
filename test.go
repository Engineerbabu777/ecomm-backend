package main;

import "fmt"

func main(){

	var testVar = 15;
	var testVarPtr *int= &testVar;
	fmt.Println(*testVarPtr);

}