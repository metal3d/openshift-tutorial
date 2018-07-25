package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/globalsign/mgo"
)

const index = "template/index.html"

type Data struct {
	Name string
}

func getDB() *mgo.Database {
	h := os.Getenv("DB_HOST")
	u := os.Getenv("MONGODB_USER")
	p := os.Getenv("MONGODB_PASSWORD")
	n := os.Getenv("MONGODB_DATABASE")
	c, err := mgo.Dial(fmt.Sprintf("%s:%s@%s:27017/%s", u, p, h, n))
	if err != nil {
		panic(err)
	}
	db := c.DB(n)
	return db
}

func SaveData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form.Get("Name"))
	d := Data{r.Form.Get("Name")}

	db := getDB()
	defer db.Session.Close()
	db.C("names").Insert(d)

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func Index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(index)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := []*Data{}
	db := getDB()
	defer db.Session.Close()
	if err := db.C("names").Find(nil).All(&data); err != nil {
		log.Println(err)
		data = []*Data{}
	}

	log.Println("Parsing and rendering")
	tpl.Execute(w, map[string]interface{}{
		"Data": data,
	})
}

func main() {
	http.HandleFunc("/save", SaveData)
	http.HandleFunc("/", Index)
	log.Println("Serving on port 8000")
	http.ListenAndServe(":8000", nil)
}
