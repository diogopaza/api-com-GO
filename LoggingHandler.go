package main

import(

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a catch-all route"))
	})
	loggeRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":3000", loggeRouter)

}