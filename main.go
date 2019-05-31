package main

import(

	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
)

var users = []User{

	User{Name:"diogo", Password:"123"},
	User{Name:"joao", Password:"56"},
	User{Name:"rodrigo", Password:"12"},
	User{Name:"beto", Password:"50"},

}

type User struct{

	Name string
	Password string
}


func getUser(w http.ResponseWriter, r *http.Request){

	json.NewEncoder(w).Encode(users)
	

}

func createUser(w http.ResponseWriter, r *http.Request){

	
	var u User
	
	body, err:= ioutil.ReadAll(r.Body)

	if err != nil{
		panic(err)
	}
	
	err = json.Unmarshal(body, &u)
	if err != nil{
		w.Header().Set("Content-Type","application/json; charset=UTF-8")
		w.WriteHeader(422)
	}
	users = append(users, u)
	
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	
	//json.NewEncoder(w).Encode(u)
}

func main(){

	var port=":8000"
	
	http.HandleFunc("/", getUser)
	http.HandleFunc("/create", createUser)
	
	fmt.Printf("Servidor rodando na porta:%s", port )
	log.Fatal(http.ListenAndServe(port, nil))

}