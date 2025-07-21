package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/represitory/sqlconnect"
	"strconv"
	"strings"
)

var (
	teachers = make(map[int]models.Teacher)
	nextID   = 1
	// mutex    = &sync.Mutex{}
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on teachers route!"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH method on teachers route!"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on teachers route!"))
	default:
		w.Write([]byte("Hello HEAD/OPTIONS method on teachers route!"))
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to Mariadb:", err)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	strID := strings.TrimSuffix(path, "/")
	fmt.Println(strID)

	if strID == "" {

		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []interface{}

		query, args = addFilters(r, query, args)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Bad query!", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		teachersList := make([]models.Teacher, 0)

		for rows.Next() {
			var teacher models.Teacher
			err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Failed to SCAN over rows and extract data!", http.StatusInternalServerError)
				return
			}
			teachersList = append(teachersList, teacher)
		}

		response := struct {
			Status string
			Count  int
			Data   []models.Teacher
		}{
			Status: "active",
			Count:  len(teachersList),
			Data:   teachersList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := strconv.Atoi(strID)
	if err != nil {
		fmt.Println("This error:", err)
		return
	}

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", userID).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found!", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Bad query!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to Mariadb:", err)
		return
	}

	var newTeachers = make([]models.Teacher, 0)

	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Failed to decode JSON data!", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Failed to PREPARE query string", http.StatusInternalServerError)
		fmt.Println("Error:", err)
		return
	}
	defer stmt.Close()

	var addedTeachers = make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Failed to EXECUTE query string", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to GET last inserted ID", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string
		Count  int
		Data   []models.Teacher
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}
