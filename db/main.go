package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
)

type Message struct {
	Message string `json:"message"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Id struct {
	Id string `json:"id"`
}

type PostUser struct {
	Name string
	Age  int
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case "GET":
		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := godotenv.Load("../.env")
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp(localhost:3306)/test_database")
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rows, err := db.Query("SELECT * FROM user WHERE name = ?", name)
		if err != nil {
			tx.Rollback()
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		user := User{}
		result := make([]User, 0)
		for rows.Next() {
			error := rows.Scan(&user.Id, &user.Name, &user.Age)
			if error != nil {
				tx.Rollback()
				println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				result = append(result, user)
			}
		}

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(result); err != nil {
			tx.Rollback()
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tx.Commit()
		fmt.Fprint(w, buf.String())
	case "POST":
		var user PostUser
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if user.Name == "" || len(user.Name) > 50 || user.Age < 20 || user.Age > 80 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := godotenv.Load("../.env")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp(localhost:3306)/test_database")
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sql := "INSERT INTO user(id,name,age) VALUES(?,?,?)"

		t := time.Now()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		id := ulid.MustNew(ulid.Timestamp(t), entropy)

		_, err = db.Exec(sql, id.String(), user.Name, user.Age)
		if err != nil {
			tx.Rollback()
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id_json := &Id{Id: id.String()}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(id_json); err != nil {
			tx.Rollback()
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tx.Commit()
		fmt.Fprint(w, buf.String())
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
