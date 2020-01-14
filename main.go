package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	hostname := os.Getenv("USER")
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			fmt.Fprint(w, "Get out")
		} else {
			fmt.Fprintf(w, "<h1>Hello, %q, %s</h1>", html.EscapeString(r.URL.Path), hostname)
		}
		helloresp, _ := http.Get("http://localhost:8080/hello?fname=frompath")
		hellobody, _ := ioutil.ReadAll(helloresp.Body)
		fmt.Println(string(hellobody))
	})
	http.HandleFunc("/json", jsonHandler)
	go http.ListenAndServe(":8080", nil)
	http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, `{"name":"mehrab"}`)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	clientname := r.FormValue("fname")
	fmt.Println(r.FormValue("lname"))
	fmt.Fprint(w, "Hello, ", clientname)
}
