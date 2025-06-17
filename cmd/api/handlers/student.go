package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"little.angels/internal/models"
)

func (app *App) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	result := app.Db.Find(&students)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Students are not found!")
			http.Error(w, "students are not found", http.StatusInternalServerError)
			return
		} else {
			log.Println("Error getting students")
			http.Error(w, "Error getting sttudents", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
	log.Println("students found!")
}

func (app *App) GetStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	var student models.Student
	result := app.Db.First(&student, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("student is not found")
			http.Error(w, "student is not found", http.StatusNotFound)
			return
		} else {
			log.Println("error getting a student")
			http.Error(w, "error getting a student", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
	log.Println("student is found!")
}

func (app *App) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var newStudent models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudent)
	if err != nil {
		log.Println("Error decoding body")
		http.Error(w, "Error getting requet body data", http.StatusBadRequest)
		return
	}

	if newStudent.Age > 5 {
		log.Println("student age is above the allowed")
		http.Error(w, "student is above the age", http.StatusBadRequest)
		return
	}

	var existingStudent models.Student
	result := app.Db.Where("first_name = ? AND last_name = ?", newStudent.FirstName, newStudent.LastName).First(&existingStudent)
	if result.Error == gorm.ErrRecordNotFound {
		result = app.Db.Create(&newStudent)
		if result.Error != nil {
			log.Println("Error creating a student")
			http.Error(w, "Error creating a student", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newStudent)
		log.Println("student was successfully created!")
	} else {
		log.Println("Student exists")
		http.Error(w, "Student exists", http.StatusConflict)
		return
	}
}

func (app *App) UpdateAStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	var newStudent models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudent)
	if err != nil {
		log.Println("Error getting body data")
		http.Error(w, "Error getting data", http.StatusBadRequest)
		return
	}
	var existingStudent models.Student
	result := app.Db.First(&existingStudent, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("student is not found")
			http.Error(w, "student is not found", http.StatusNotFound)
			return
		} else {
			log.Println("error getting a student")
			http.Error(w, "error getting a student", http.StatusInternalServerError)
			return
		}
	}
	if newStudent.FirstName != "" {
		existingStudent.FirstName = newStudent.FirstName
	}
	if newStudent.LastName != "" {
		existingStudent.LastName = newStudent.LastName
	}
	if newStudent.Gender != "" {
		existingStudent.Gender = newStudent.Gender
	}
	if newStudent.Age != 0 {
		existingStudent.Age = newStudent.Age
	}
	if newStudent.Guardian != nil {
		existingStudent.Guardian = newStudent.Guardian
	}
	result = app.Db.Save(&existingStudent)
	if result.Error != nil {
		log.Println("Error updating the student")
		http.Error(w, "Error updating the student", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingStudent)
	log.Println("student is updated")

}

func (app *App) DeleteAStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid paramater!", http.StatusBadRequest)
		return
	}
	var student models.Student
	result := app.Db.First(&student, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("value is not found")
		http.Error(w, "value is not found", http.StatusNotFound)
		return
	}
	result = app.Db.Delete(&student)
	if result.RowsAffected == 0 {
		log.Println("value was not deleted")
		http.Error(w, "value was not deleted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Println("student was deleted")
}
