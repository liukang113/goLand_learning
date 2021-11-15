package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"path"
)

var db *sql.DB

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:qqq110@/svn")
	if err != nil {
		panic(err)
	}
}

func getUserByName(username string) (u user) {
	u = user{}
	rows, err := db.Query("select password from user where username = ?", username)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var password string
		err := rows.Scan(&password)
		if err != nil {
			panic(err)
		}
		u = user{
			Username: username,
			Password: password,
		}
	}
	return u
}

func create(u user) {
	stmt, err := db.Prepare("insert into user (username, password) values (?, ?)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(u.Username, u.Password)
	if err != nil {
		panic(err)
	}
}
func modify(u user) {
	stmt, err := db.Prepare("update user set password = ? where username = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(u.Password, u.Username)
	if err != nil {
		panic(err)
	}
}
func dele(username string) {
	stmt, err := db.Prepare("delete from user where username = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(username)
	if err != nil {
		panic(err)
	}
}

//多路复用器
func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	case "PUT":
		handleModify(w, r)
	case "DELETE":
		handleDelete(w, r)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	//返回路径中最后一个元素
	username := path.Base(r.URL.Path)
	u := getUserByName(username)
	//b, err := json.MarshalIndent(&u,"","\t\t")
	b, err := json.Marshal(u)
	if err != nil {
		fmt.Println("json error", err)
	}
	fmt.Println(u)
	fmt.Println(string(b))
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func handlePost(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var u user
	json.Unmarshal(body, &u)
	create(u)
	w.WriteHeader(200)

}
func handleModify(w http.ResponseWriter, r *http.Request) {
	username := path.Base(r.URL.Path)
	u := getUserByName(username)

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &u)
	modify(u)
	w.WriteHeader(200)
}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	username := path.Base(r.URL.Path)
	dele(username)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/user/", handleRequest)
	server.ListenAndServe()
}
