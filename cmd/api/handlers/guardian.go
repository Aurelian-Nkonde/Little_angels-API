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

func (app *App) GetAllGuardians(w http.ResponseWriter, r *http.Request) {
	var guardians []models.Guardian
	result := app.Db.Find(&guardians)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("guardians are not found")
		http.Error(w, "guardians are not found", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		log.Println("error getting guardians")
		http.Error(w, "error gettting guardians", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guardians)
	log.Println("guardians are found!")
}

func (app *App) GetAGuardian(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	var guardian models.Guardian
	result := app.Db.Find(&guardian, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("guardian is not found")
		http.Error(w, "guardia id not found", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		log.Println("Error getting a guardian")
		http.Error(w, "error getting a  guardia", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guardian)
	log.Println("guardian is found")
}

func (app *App) UpdateGuardian(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid paramter", http.StatusBadRequest)
		return
	}
	var newGuardian models.Guardian
	err := json.NewDecoder(r.Body).Decode(&newGuardian)
	if err != nil {
		log.Println("Error getting request data")
		http.Error(w, "Error getting request data", http.StatusBadRequest)
		return
	}
	var existingGuardian models.Guardian
	result := app.Db.First(&existingGuardian, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Panicln("guardian is not found")
		http.Error(w, "Guardian already exist", http.StatusInternalServerError)
		return
	}
	if newGuardian.FirstName != "" {
		existingGuardian.FirstName = newGuardian.FirstName
	}
	if newGuardian.LastName != "" {
		existingGuardian.LastName = newGuardian.LastName
	}
	if newGuardian.Email != nil {
		existingGuardian.Email = newGuardian.Email
	}
	if newGuardian.Gender != "" {
		existingGuardian.Gender = newGuardian.Gender
	}
	if newGuardian.Proffession != nil {
		existingGuardian.Proffession = newGuardian.Proffession
	}
	result = app.Db.Save(&existingGuardian)
	if result.Error != nil {
		log.Println("Error updating a guardian")
		http.Error(w, "Error updating a guardian", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingGuardian)
	log.Println("guardian is successfully updated")
}

func (app *App) DeleteGuardian(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
	var guardian models.Guardian
	result := app.Db.First(&guardian, id)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("guardian is not found")
		http.Error(w, "guardian is not found", http.StatusNotFound)
		return
	}
	result = app.Db.Delete(&guardian)
	if result.Error != nil {
		log.Println("Error deleting the guardian")
		http.Error(w, "Error deleting the guardian", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Println("gaurdian was deleted")
}

func (app *App) CreateGuardian(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAGuardianStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id, err := strconv.Atoi(id); err != nil || id < 0 {
		log.Println("Invalid parameter")
		http.Error(w, "Invalid parameter", http.StatusBadRequest)
		return
	}
}
