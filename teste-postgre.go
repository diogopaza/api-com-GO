package main

import(

	"fmt"
	"database/sql"
	"net/http"
	"log"
	"encoding/json"
	_ "github.com/lib/pq" //postgresql
)



type Users struct{

	ID string
	LOGIN string
	NAME string
	PASSWORD string
	TOKEN string
	PROFILE_ID string
	
}

const(
	host = "localhost"
	port=5432
	user="postgres"
	password_admin="123321"
	dbname="gestor"
)

func getUsers(w http.ResponseWriter, r *http.Request){

	connectingDB:= initDb()
	myUsers,err := returnArrayUsers(connectingDB)
	if err != nil{
		w.Header().Set("Content-Type","application/json; charset=UTF-8")
		w.WriteHeader(400)
	}
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(myUsers)
	w.WriteHeader(http.StatusOK)
	
	

}

func getLogin(w http.ResponseWriter, r *http.Request){

	connecting:= initDb()
	var count int
	var u Users

	r.ParseForm()
	name := r.Form["user"]
	password := r.Form["password"]
	u.NAME = name
	u.PASSWORD = password
	fmt.Println(name)
	fmt.Println(password)

	rowsCount, err:= connecting.Query("SELECT COUNT(*) as count FROM public.user WHERE id=1")
		if err != nil{
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
		}
	for rowsCount.Next(){
		err:= rowsCount.Scan(&count)
		if err != nil{
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
		}
	}

	

}

func returnArrayUsers(connecting *sql.DB) ([]Users,error) {

	rows, err := connecting.Query("SELECT * FROM public.user")
	if err != nil{
		return nil, err
	}

	var id string
	var login string
	var name string
	var password string
	var token string
	var profile_id string
	
	u := Users{}
	res :=[]Users{}

	for rows.Next(){
		
		err = rows.Scan(&id,&login,&name,&password,&token,&profile_id)
		if err != nil{
			return nil,err
		}
	
		u.ID = id
		u.LOGIN = login
		u.NAME = name
		u.PASSWORD = password
		u.TOKEN = token
		u.PROFILE_ID = profile_id

		res = append(res, u)
	
	
	}
		
		return res, nil


}

func initDb() *sql.DB{

	banco := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable\n", host, port, user, password_admin, dbname)
	

	db, err := sql.Open("postgres", banco)
	if err != nil{
		panic(err)
	}
	

	err = db.Ping()
	if err != nil{
		panic(err)
	}

	fmt.Println("Sucessfully connected!!!")
	return db

}


func main(){	
	
	
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/login", getLogin)
	

	var porta = 8000
	fmt.Printf("Rodando na %d", porta)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", porta), nil))

}