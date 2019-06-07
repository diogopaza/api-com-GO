package main

import(

	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	_ "github.com/lib/pq" //postgresql
	
	"github.com/gorilla/mux"
	
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

var getUsers = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	connectingDB:= initDb()
	myUsers,err := returnArrayUsers(connectingDB)
	if err != nil{
		w.Header().Set("Content-Type","application/json; charset=UTF-8")
		w.WriteHeader(400)
	}
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(myUsers)
	w.WriteHeader(http.StatusOK)
	
	

})

var getLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){


	conn:= initDb()
	
	var user Users
	
	erro := r.ParseForm()
	if erro != nil{
		panic(erro)
	}
	login := r.FormValue("login")
	password:= r.FormValue("password")
	fmt.Println("Variavel do POST: ", login) 
	
	sqlQuery := "SELECT id, name, password FROM public.USER WHERE name=$1"

	row := conn.QueryRow(sqlQuery, login)

	err := row.Scan(&user.ID, &user.NAME, &user.PASSWORD)

	if(err!=nil){
		if(err == sql.ErrNoRows){
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
			fmt.Println("Usuario não encontrado")
		}
	
		fmt.Println("Erros:%v", err)
	}else{
		fmt.Println("Achou o usuario")
		fmt.Println("ID:"+user.ID)
		fmt.Println("NAME:"+user.NAME)
		fmt.Println("PASSWORD:"+user.PASSWORD)

		if(password!=user.PASSWORD){
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
			fmt.Println("Erro senha incorreta")
		}else{
			w.Write([]byte("Retornando token"))
			fmt.Println("Aqui retorna o token")
		}

	}

	//var count int
	/*
	var u Users
	
	r.ParseForm()
	name := r.Form["user"]
	password := r.Form["password"]
	u.NAME = name
	u.PASSWORD = password
	fmt.Println(name)
	fmt.Println(password)
	*/
	//rowsCount, err:= connecting.Query("SELECT COUNT(*) as count FROM public.user WHERE id=6")
	
	//defer rowsCount.Close();

	/*
	if err != nil{
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
			fmt.Println("Erro")
		}



	for rowsCount.Next(){
		err:= rowsCount.Scan(&count)
		if err != nil{
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
		}
	}
	*/

	

})

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
	
	r := mux.NewRouter()
	
	http.Handle("/", r)
	r.Handle("/users", getUsers).Methods("GET")
	r.Handle("/login", getLogin).Methods("POST")
	
	http.ListenAndServe(":3000", nil)

}