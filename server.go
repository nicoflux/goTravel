package main

import (
	"fmt"      // import the fmt package
	"net/http" // import the http package
)

func searchHandler(w http.ResponseWriter, r *http.Request) { // function that handles the request
	if r.Method != http.MethodGet { // if the request method is not GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // return an error
		return
	}

	fmt.Fprintf(w, "searching") // write "searching" in the response
}

func bookingHandler(w http.ResponseWriter, r *http.Request) { // function that handles the request
	/* if r.Method != http.MethodGet { // if the request method is not GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // return an error
		return
	} */
	switch r.Method { //https://www.golangprograms.com/example-to-handle-get-and-post-request-in-golang.html
	case "GET":
		fmt.Fprintf(w, "booking GET") // write "booking" in the response
	case "POST":
		fmt.Fprintf(w, "booking POST") // write "booking" in the response
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/api/search", searchHandler)   // handle the request on /api/search
	http.HandleFunc("/api/booking", bookingHandler) // handle the request on /api/booking
	fmt.Println("Server is listening on :5000")     // print a message
	http.ListenAndServe("127.0.0.1:5000", nil)      // listen on port 5000
}
