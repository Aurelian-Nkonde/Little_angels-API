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

type App struct {
	Db *gorm.DB
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!"))
}

func (app *App) CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	err := json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println("Error request body!")
		return
	}

	var existingTeacher models.Teacher
	result := app.Db.Where("first_name = ? AND last_name = ?", teacher.FirstName, teacher.LastName).First(&existingTeacher)
	if result.Error == gorm.ErrRecordNotFound {
		// teacher does not exist!
		newTeacher := app.Db.Create(&teacher)
		if newTeacher.Error != nil {
			log.Printf("Error creating the teacher:%v", newTeacher.Error)
			http.Error(w, "Error creating a teacher", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(teacher)
	} else {
		log.Println("Teacher already exist!")
		http.Error(w, "Teacher already exist!", http.StatusConflict)
	}
}

func (app *App) UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Incorrect id parameter")
		http.Error(w, "Incorrect id parameter", http.StatusBadRequest)
		return
	}
	var newTeacherDetails models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeacherDetails)
	if err != nil {
		log.Println("Error decoding incoming data")
		http.Error(w, "Error decoding data", http.StatusInternalServerError)
		return
	}
	var existingTeacher models.Teacher
	result := app.Db.First(&existingTeacher, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("teacher value was not found")
			http.Error(w, "teacher was not found", http.StatusNotFound)
			return
		} else {
			log.Println("error getting the teacher value")
			http.Error(w, "Error getting the value", http.StatusInternalServerError)
			return
		}
	}
	if newTeacherDetails.FirstName != "" {
		existingTeacher.FirstName = newTeacherDetails.FirstName
	}
	if newTeacherDetails.LastName != "" {
		existingTeacher.LastName = newTeacherDetails.LastName
	}
	if newTeacherDetails.Gender != "" {
		existingTeacher.Gender = newTeacherDetails.Gender
	}
	if newTeacherDetails.Age != 0 {
		existingTeacher.Age = newTeacherDetails.Age
	}
	if newTeacherDetails.Email != nil {
		existingTeacher.Email = newTeacherDetails.Email
	}

	result = app.Db.Save(&existingTeacher)
	if result.Error != nil {
		log.Println("Error saving the teacher")
		http.Error(w, "Error saving the teacher", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingTeacher)
	log.Println("teacher details updated")
}

func (app *App) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Wrong parameter")
		http.Error(w, "Wrong parameter", http.StatusBadRequest)
		return
	}
	var teacher models.Teacher
	result := app.Db.First(&teacher, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("Teacher record is not found")
		http.Error(w, "Teacher record is not found", http.StatusInternalServerError)
		return
	}
	deleted := app.Db.Delete(&teacher)
	if deleted.RowsAffected == 0 {
		log.Println("Teacher was not deleted")
		http.Error(w, "value was not found", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Println("teacher value was deleted")
}

func (app *App) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers = make([]models.Teacher, 0)
	result := app.Db.Find(&teachers)
	if result.Error != nil {
		log.Println("Error retrieving the teachers")
		http.Error(w, "Error retrieving teahers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teachers)
	log.Println("teachers found!")
}

func (app *App) GetATeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Error id parameter")
		http.Error(w, "Wrong id parameter", http.StatusBadRequest)
		return
	}
	var teacher models.Teacher
	result := app.Db.First(&teacher, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("Teacher is not found!")
		http.Error(w, "Teacher is not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher)
	log.Println("teacher is found")

}

func (app *App) GetAllAdmins(w http.ResponseWriter, r *http.Request) {
	var admins []models.Teacher
	result := app.Db.Where("role = ?", "admin").Find(&admins)
	if result.Error != nil {
		log.Println("Error retrieving all admins")
		http.Error(w, "Admins are not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(admins)
	log.Println("admins found")
}

func (app *App) MakeATeacherTheAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Wrong parameter")
		http.Error(w, "Wrong parameter value", http.StatusBadRequest)
		return
	}
	var existingTeacher models.Teacher
	result := app.Db.First(&existingTeacher, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("teacher was not found")
			http.Error(w, "teacher was not found", http.StatusNotFound)
			return
		} else {
			log.Println("Error getting a teacher")
			http.Error(w, "Error getting a teacher", http.StatusInternalServerError)
			return
		}
	}
	existingTeacher.Role = "admin"
	result = app.Db.Save(&existingTeacher)
	if result.Error != nil {
		log.Println("Error saving the updated value")
		http.Error(w, "Error saving the value", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingTeacher)
	log.Println("value is updated")
}
