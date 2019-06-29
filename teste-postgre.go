package main

import(

	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	_ "github.com/lib/pq" //postgresql

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"time"
	"strings"
	
	
	
	
)

var myKey = []byte("secret")
var loginGlobal string
var passwordGlobal string

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
	
	connectingDB := initDb()
	
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
	loginGlobal := r.FormValue("login")
	passwordGlobal := r.FormValue("password")
	fmt.Println("Variavel do POST: ", loginGlobal) 
	
	sqlQuery := "SELECT id, name, password FROM public.USER WHERE name=$1"

	row := conn.QueryRow(sqlQuery, loginGlobal)

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

		if(passwordGlobal != user.PASSWORD){
			w.Header().Set("Content-Type","application/json; charset=UTF-8")
			w.WriteHeader(400)
			fmt.Println("Erro senha incorreta")
		}else{
			w.Header().Set("Content-Type", "application/json")
			
			myToken:= getToken(user)
			TokenSplit := strings.Split(myToken, ".")
			fmt.Println("Split: ", TokenSplit[2])
			
			var m = make(map[string]string)
			m["token"]=myToken

			saveToken(conn, TokenSplit[2], user.ID)
			json.NewEncoder(w).Encode(m)

			
		}

	}

})

func returnArrayUsers(connecting *sql.DB) ([]Users,error) {
	fmt.Println("estou returnUsers")
	rows, err := connecting.Query("SELECT * FROM public.user")
	if err != nil{
		fmt.Println("Não foi pesquisar usuários")
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
		fmt.Println(" next usuários")
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
		fmt.Println("users:",res)
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


func getToken(u Users) string{

	token := jwt.New(jwt.SigningMethodHS256)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"]=true
	claims["id"]= u.ID
	claims["name"]= u.NAME
	claims["password"]= u.PASSWORD
	claims["exp"]=time.Now().Add(time.Minute * 45).Unix()
	
	tokenString, _ := token.SignedString(myKey)
	
	return tokenString

}

func middlewareJWT( h http.HandlerFunc ) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request){
		
		var user Users
		//pega os dados do Header
		auth := r.Header.Get("Authorization")		
		tokenString := strings.Split(auth, " ")	
		tokenSplit := strings.Split(tokenString[1], ".")
		fmt.Println(tokenSplit[2])
				
		//verifica se o token existe e se é válido e se a chave de autenticação está correta
		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil{
				w.Header().Set("Content-Type","application/json; charset=UTF-8")
				w.WriteHeader(401)
			  json.NewEncoder(w).Encode("Invalid authorization token")
			 return 

		}else{
				
		  //pega nome de usuário que é único no banco de dados
			claims:= token.Claims.(jwt.MapClaims)
			mapstructure.Decode(claims, &user)
			signedTokenLocalBank := signedToken(user.NAME)

			fmt.Println(signedTokenLocalBank)
			if signedTokenLocalBank == tokenSplit[2]{
				h.ServeHTTP(w,r)
			}else{
				w.Header().Set("Content-Type","application/json; charset=UTF-8")
				w.WriteHeader(401)
				json.NewEncoder(w).Encode("token nao corresponde")
			}
			
		}		
			
	

		
	
		
	}
}

var setupResponse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
})

func saveToken(conn *sql.DB, signatureToken string, id string){
	
	sqlQuery := "UPDATE public.user SET token=$2 WHERE id=$1"
	_ ,err := conn.Exec(sqlQuery, id, signatureToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token atualizado com sucesso")

}  

func signedToken(name string) (token string){

	var userToken string

	connectingDB:= initDb()
	sqlQuery := "SELECT token FROM public.user WHERE name=" + "'" + name + "'"
	rows, err := connectingDB.Query(sqlQuery)
	if err != nil{
		fmt.Println("Erro ao consultar usuario no banco de dados")
		return 
	}
	
	for rows.Next(){

		err = rows.Scan(&token)
		if err != nil{
			fmt.Println("Erro ao percorrer no banco de dados")
			return 
		}

		userToken = token
	
	

	}
	
	return userToken

} 

var searchName = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){


	var usersName []string
  var userName string

	erro := r.ParseForm()
	if erro != nil{
		panic(erro)
	}
	campoPesquisa := r.FormValue("pesquisa")
	
	connectingDB:= initDb()

	sqlQuery := `SELECT name FROM public.user WHERE name ILIKE '%' || $1 || '%' `  
		
	rows, err := connectingDB.Query(sqlQuery, campoPesquisa)
	if err != nil {
		fmt.Println("Erro ao percorrer dados")
}
if rows == nil {
	fmt.Println("Nenhuma informação localizada ")
}

for rows.Next(){
	err := rows.Scan(&userName)
	if err != nil{
		fmt.Println("Erro ao indexar banco de dados")
	}

	usersName = append(usersName, userName)
}

/*
	//pesquisa antiga no banco de dados sem parametros
	rows, err := connectingDB.Query(sqlQuery)	
	if err != nil{
		fmt.Println("Erro ao consultar usuario no banco de dados")
		return 
	}	
	for rows.Next(){

		err = rows.Scan(&userName)
		usersName = append(usersName, userName)
		if err != nil{
			fmt.Println("Erro ao percorrer no banco de dados")
			return 
		}
	}			
	*/
		w.Header().Set("Content-Type","application/json")
	  json.NewEncoder(w).Encode(usersName)
		
})


func main(){	
	
	router := mux.NewRouter()

	allowedHeaders := handlers.AllowedHeaders([]string{"Cache-Control","X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
  allowedOrigins := handlers.AllowedOrigins([]string{"*"})
  allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	//router.Handle("/login", setupResponse).Methods("OPTIONS")
	router.Handle("/users", middlewareJWT(getUsers)).Methods("GET")
	router.Handle("/login", getLogin).Methods("POST")
	router.Handle("/search-name", middlewareJWT(searchName)).Methods("POST")

	
	
	fmt.Println("Rodando na 3000")
	http.ListenAndServe(":3000", handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
	

}