package main

import(

	"fmt"
	"encoding/base64"
)

//https://dev.to/hitman666/jwt-authentication-in-an-angular-application-with-a-go-backend--13cg
func main(){

	data := `{"user_id":"a1b2c3","username":"nikola"}`
	uEnc := base64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)

	uDec, _ := base64.URLEncoding.DecodeString("https://www.google.com/search?ei=N9XvXJ_QKqWW0AbBuJ64CA&q=cgn&oq=cgn&gs_l=psy-ab.3..35i39j0i67l6j0l2j0i67.190093.190601..190733...0.0..0.126.251.0j2......0....1..gws-wiz.......0i71j0i131.mmwvAB4CKDw")
	fmt.Println(uDec)

	header := `{"alg":"HS256","typ":"JWT"}`
	encHeader := base64.URLEncoding.EncodeToString([]byte(header))
	fmt.Println(encHeader)
}