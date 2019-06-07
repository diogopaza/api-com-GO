package main

import(

	"net/http"
	"fmt"

)

type User struct{

	Nome string
}


func getUser(w http.ResponseWriter, r *http.Request){
	/*
	data, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(err)
	}
	

	fmt.Println(string(data))
	*/
	
	var u User
	
	r.ParseForm()
	nome := r.Form["usuario"]
	u.Nome = nome
	fmt.Println(u.Nome)


}


func main(){

	http.HandleFunc("/user", getUser)


	fmt.Println("Rodando na 3000")
	http.ListenAndServe(":3000", nil)

}