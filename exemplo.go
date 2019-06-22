package main

import(

	"fmt"
	
	
	
	"github.com/dgrijalva/jwt-go"
	
	

)

func parseToken(){

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTYxMjMyNjg2LCJpZCI6IjEiLCJuYW1lIjoiZGlvZ28ifQ.7nHwRv646U_REXTYCOPrx2iRNKGR2gLwBmpQKdanfmI"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil{
		fmt.Println("Erro")
	}else{
		fmt.Println(token)
	}
}



func main(){
	
	parseToken()
	/*
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		  return []byte("secret"), nil
		},
	   
		SigningMethod: jwt.SigningMethodHS256,
	  })
	*/
	//app := jwtMiddleware.Handler(myHandler)
	//http.ListenAndServe(":3000", app)

}