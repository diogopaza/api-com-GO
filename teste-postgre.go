package main

import(

	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	_ "github.com/lib/pq" //postgresql
	"github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dgrijalva/jwt-go"
	"time"
	
)

var myKey = []byte("secret")

type Users struct{

	ID string `json:"id"`
	LOGIN string `json:"login"`
	NAME string `json:"nome"`
	PASSWORD string `json:"password"`
	TOKEN string `json:"token"`
	PROFILE_ID string `json:"profileId"`
	
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

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

var getLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	fmt.Println("ENTREI LOGIN")
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

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
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			fmt.Println("Aqui retorna o token")
		}

	}

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


var getTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	token := jwt.New(jwt.SigningMethodHS256)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"]=true
	claims["name"]="Diogo Paza"
	claims["exp"]=time.Now().Add(time.Hour * 24).Unix()
	
	tokenString, _ := token.SignedString(myKey)
	w.Write([]byte(tokenString))


})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	  return myKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})


func main(){	
	
	router := mux.NewRouter()

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
    allowedOrigins := handlers.AllowedOrigins([]string{"*"})
    allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	router.Handle("/users", jwtMiddleware.Handler(getUsers)).Methods("GET")
	router.Handle("/login", jwtMiddleware.Handler(getLogin)).Methods("POST")
	
	http.ListenAndServe(":3000", handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
	//http.ListenAndServe(":3000", nil)

}