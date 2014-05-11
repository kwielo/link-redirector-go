package main

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

func main() {
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
		http.Error(w, fmt.Sprintf("No redirection for: %s", rdct.slug), http.StatusNotFound)
	} else {
		http.Redirect(w, r, rdct.url, http.StatusMovedPermanently)
	}

}

func getLink(rdct Redirection) Redirection {
	mc := memcache.New("127.0.0.1:11211")

	mcRedirection, err := mc.Get(rdct.slug)
	if err != nil {
		log.Printf("Cache miss for %s", rdct.slug)
	}

	if mcRedirection != nil {
		log.Printf("Found in cache: %s => %s", rdct.slug, string(mcRedirection.Value))
		rdct.url = string(mcRedirection.Value)
	}

	config := getConfig()
	db, err := sql.Open("mysql", config.Db.User+":"+config.Db.Pass+"@/"+config.Db.Name)

	var redirection string
	if err = db.QueryRow("SELECT url FROM link WHERE tag = ?", rdct.slug).Scan(&redirection); err != nil {
		rdct.err = "Redirection doesn't exist"
	} else {
		mc.Set(&memcache.Item{Key: rdct.slug, Value: []byte(redirection)})
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
