package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Visited: %s\n", r.URL.Path[1:])
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

var statuses = []int{
	http.StatusInternalServerError,
	http.StatusNotFound,
	http.StatusMovedPermanently,
	http.StatusOK,
	http.StatusContinue,
}

func randPage(w http.ResponseWriter, r *http.Request) {
	status := statuses[rand.Intn(len(statuses))]
	w.WriteHeader(status)
	fmt.Fprintf(w, "Testing...")
	fmt.Printf("Sent -> %d\n", status)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/rand", randPage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
