package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getUserID(path, route string) (int, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == route {
		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Failed to convert ID to integer value:", err)
			return 0, false
		}
		return userID, true
	}
	return 0, false
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, root route!"))
	fmt.Println("Hello, root route!")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// fmt.Println("Path:", r.URL.Path)
		// path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		// userID := strings.TrimSuffix(path, "/")
		// fmt.Println("The ID is:", userID)

		//teachers/?key=value&query=value2&sortby=email&sortorder=ASC
		fmt.Println("Path:", r.URL.Path)
		userID, isID := getUserID(r.URL.Path, "teachers")
		if !isID {
			fmt.Println("Failed to fetch user ID!")
		} else {
			fmt.Println("The ID is:", userID)
		}

		fmt.Println("Query params:", r.URL.Query())

		queryParams := r.URL.Query()
		key := queryParams.Get("key")
		query := queryParams.Get("query")
		sortby := queryParams.Get("sortby")
		sortorder := queryParams.Get("sortorder")

		if sortorder == "" {
			sortorder = "DESC"
		}

		fmt.Println("Key:", key)
		fmt.Println("Query:", query)
		fmt.Println("Sort by:", sortby)
		fmt.Println("Sort order:", sortorder)

		w.Write([]byte("Hello GET method on teachers route!"))
		fmt.Println("Hello GET method on teachers route!")
	case http.MethodPost:
		w.Write([]byte("Hello Post method on teachers route!"))
		fmt.Println("Hello Post method on teachers route!")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on teachers route!"))
		fmt.Println("Hello Put method on teachers route!")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on teachers route!"))
		fmt.Println("Hello Patch method on teachers route!")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on teachers route!"))
		fmt.Println("Hello Delete method on teachers route!")
	}
	// w.Write([]byte("Hello, teachers route!"))
	// fmt.Println("Hello,teachers route!")
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on students route!"))
		fmt.Println("Hello GET method on students route!")
	case http.MethodPost:
		w.Write([]byte("Hello Post method on students route!"))
		fmt.Println("Hello Post method on students route!")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on students route!"))
		fmt.Println("Hello Put method on students route!")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on students route!"))
		fmt.Println("Hello Patch method on students route!")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on students route!"))
		fmt.Println("Hello Delete method on students route!")
	}
	// w.Write([]byte("Hello, students route!"))
	// fmt.Println("Hello,students route!")
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on execs route!"))
		fmt.Println("Hello GET method on execs route!")
	case http.MethodPost:
		w.Write([]byte("Hello Post method on execs route!"))
		fmt.Println("Hello Post method on execs route!")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on execs route!"))
		fmt.Println("Hello Put method on execs route!")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on execs route!"))
		fmt.Println("Hello Patch method on execs route!")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on execs route!"))
		fmt.Println("Hello Delete method on execs route!")
	}
	// w.Write([]byte("Hello, execs route!"))
	// fmt.Println("Hello,execs route!")
}

func main() {
	port := ":3000"

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/teachers/", teachersHandler)
	http.HandleFunc("/students/", studentsHandler)
	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server started running on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Failed to listen on port", port)
	}
}
