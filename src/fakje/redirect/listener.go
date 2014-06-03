package main

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

var config Config
var db *sql.DB
var cache = map[string]string{}

func main() {
	config = getConfig()
	db, _ = sql.Open("mysql", config.Db.User+":"+config.Db.Pass+"@/"+config.Db.Name)
	// if err != nil {
	// 	log.Printf("Error with db connection: %s", err.Error())
	// }
	http.HandleFunc("/", linkHandler)
	http.ListenAndServe(":8088", nil)
}

type Redirection struct {
	slug, url, err string
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.Replace(fmt.Sprintf("%s", r.URL), "/", "", -1)
	if slug == "favicon.ico" {
		http.Error(w, "No redirection.", http.StatusNotFound)
		return
	}

	rdct := Redirection{slug, "", ""}
	rdct = getLink(rdct)

	if rdct.err != "" {
		http.Error(w, fmt.Sprintf("No redirection for: %s, cause: %s", rdct.slug, rdct.err), http.StatusNotFound)
	} else {
		http.Redirect(w, r, rdct.url, http.StatusMovedPermanently)
		logRedirection(rdct)
	}

}

func logRedirection(rdct Redirection) {
	_, err := db.Exec("UPDATE link SET clicks=(clicks+1) WHERE slug = ?", rdct.slug)
	if err != nil {
		log.Printf("There was a problem during saving the click count.")
	}
}

func getLink(rdct Redirection) Redirection {

	url := cache[rdct.slug]
	if url == "" {
		// log.Printf("Cache miss for %s", rdct.slug)
	} else {
		// log.Printf("Found in cache: %s => %s", rdct.slug, url)
		rdct.url = url
		return rdct
	}

	var redirection string
	if err := db.QueryRow("SELECT url FROM link WHERE slug = ?", rdct.slug).Scan(&redirection); err != nil {
		rdct.err = "Redirection doesn't exist"
	} else {
		cache[rdct.slug] = redirection
		rdct.url = redirection
	}

	return rdct
}

type Config struct {
	Db struct {
		Host string
		User string
		Pass string
		Name string
	}
}

func getConfig() Config {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		panic(fmt.Sprintf("error reading config file: ", err.Error()))
	}
	return cfg
}
