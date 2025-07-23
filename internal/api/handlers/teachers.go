package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/internal/represitory/sqlconnect"
	"strconv"
	"strings"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
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

func sortParams(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func GetTeachers(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to Mariadb:", err)
		return
	}

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = addFilters(r, query, args)

	query = sortParams(r, query)

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
}

func GetTeacher(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to Mariadb:", err)
		return
	}

	strID := r.PathValue("id")

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

func AddTeachers(w http.ResponseWriter, r *http.Request) {

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

func UpdateTeacher(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		http.Error(w, "Failed to connect to database!", http.StatusBadRequest)
		return
	}
	defer db.Close()

	strID := r.PathValue("id")

	userID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("Failed to parse string to int:", err)
		http.Error(w, "Failed to parse string to int!", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request body!", http.StatusInternalServerError)
		return
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", userID).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println("Teacher not found:", err)
		http.Error(w, "Teacher not found!", http.StatusInternalServerError)
		return
	} else if err != nil {
		log.Println("Failed to scan over query row:", err)
		http.Error(w, "Failed to scan over query row!", http.StatusInternalServerError)
		return
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		log.Println("Failed to update data in database:", err)
		http.Error(w, "Failed to update data in database!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Status         string
		RemovedTeacher models.Teacher
		UpdatedTeacher models.Teacher
	}{
		Status:         "success",
		RemovedTeacher: existingTeacher,
		UpdatedTeacher: updatedTeacher,
	}

	json.NewEncoder(w).Encode(response)
}

func PatchTeachers(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		http.Error(w, "Failed to connect to database!", http.StatusBadRequest)
		return
	}
	defer db.Close()

	strID := r.PathValue("id")

	userID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("Failed to parse string to int:", err)
		http.Error(w, "Failed to parse string to int!", http.StatusBadRequest)
		return
	}

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request body!", http.StatusInternalServerError)
		return
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", userID).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println("Teacher not found:", err)
		http.Error(w, "Teacher not found!", http.StatusInternalServerError)
		return
	} else if err != nil {
		log.Println("Failed to scan over query row:", err)
		http.Error(w, "Failed to scan over query row!", http.StatusInternalServerError)
		return
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	fmt.Println(teacherVal)
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		log.Println("Failed to update data in database:", err)
		http.Error(w, "Failed to update data in database!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Status         string
		UpdatedTeacher models.Teacher
	}{
		Status:         "success",
		UpdatedTeacher: existingTeacher,
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteTeachers(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		http.Error(w, "Failed to connect to database!", http.StatusBadRequest)
		return
	}
	defer db.Close()

	strID := r.PathValue("id")

	userID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("Failed to parse string to int:", err)
		http.Error(w, "Failed to parse string to int!", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", userID).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println("Teacher not found:", err)
		http.Error(w, "Teacher not found!", http.StatusInternalServerError)
		return
	} else if err != nil {
		log.Println("Failed to scan over query row:", err)
		http.Error(w, "Failed to scan over query row!", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM teachers WHERE id =?", userID)
	if err != nil {
		log.Println("Failed to delete data in database:", err)
		http.Error(w, "Failed to delete data in database!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Status         string
		DeletedTeacher int
	}{
		Status:         "Teacher successfully deleted",
		DeletedTeacher: existingTeacher.ID,
	}

	json.NewEncoder(w).Encode(response)
}
