//https://auth0.com/blog/authentication-in-golang/#Building-an-API-in-Go

package main

import(

	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	
	"os"
	"encoding/json"
)

type Product struct{

	Id int
	Name string
	Slug string
	Description string

}

var myKey = []byte("secret")

var products = []Product{
	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description : "Shoot your way to the top on 14 different hoverboards"},
	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description : "Explore the depths of the sea in this one of a kind underwater experience"},
	Product{Id: 3, Name: "Dinosaur Park", Slug : "dinosaur-park", Description : "Go back 65 million years in the past and ride a T-Rex"},
	Product{Id: 4, Name: "Cars VR", Slug : "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description : "Pick up the bow and arrow and master the art of archery"},
	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description : "Explore the seven wonders of the world in VR"},

}

var index = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Index"))
})

var statusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("API is up and runnig"))

})

var productHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	
	payload, _ := json.Marshal(products)
	w.Header().Set("Content-Type","application/json")
	
	w.Write([]byte(payload))

})

var addFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	
	decoder := json.NewDecoder(r.Body)
	var product Product

	err := decoder.Decode(&product)
	if err != nil{
		panic(err)
	}
	payload, _ := json.Marshal(product)
	w.Write([]byte(payload))
	//json.NewEncoder(w).Encode(product)
	
})

var getTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	w.Write([]byte(myKey))

})

	
func main(){

	r := mux.NewRouter()
	
	r.Handle("/status", statusHandler).Methods("GET")
	r.Handle("/products", productHandler).Methods("GET")
	r.Handle("/products/feedback", addFeedbackHandler).Methods("POST")
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.Handle("/get-token", getTokenHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
	
	
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	                            
}