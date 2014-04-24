package main

import (
	"net/http"
	// "github.com/bradfitz/gomemcache/memcache"
	// "fmt"
	"log"
)

func main() {
	http.HandleFunc("/", linkHandler)
	http.ListenAndServe(":8088", nil)
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Page loaded")

	// w.WriteHeader(http.StatusTemporaryRedirect)

	http.Redirect(w, r, "http://google.pl", http.StatusTemporaryRedirect)
	// fmt.Fprintf(w, "<h1>Yo!%v</h1>", r.URL)
}

// func getLink() {
// 	mc := memcache.New("127.0.0.1:11211")
// 	mc.Set(&memcache.Item{Key: "fakje_1", Value: []byte("new value")})

// 	link, err := mc.Get("fakje_1")
// 	//TODO: this is a very bad error handling, change it!
// 	if err != nil {
// 		fmt.Errorf("err: %v\n", err)
// 	}
// }