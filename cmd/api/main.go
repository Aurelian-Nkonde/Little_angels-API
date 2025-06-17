package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"little.angels/cmd/api/handlers"
	"little.angels/internal/models"
)

type App struct {
	Db *gorm.DB
}

func main() {

	port := flag.String("port", ":5000", "port address for this api!")
	flag.Parse()

	dns := "root:Thousand@90@tcp(127.0.0.1:3306)/angels?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to a db,", err)
	}

	sqlDb, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database object", err)
	}

	if err := sqlDb.Ping(); err != nil {
		log.Fatal("Database ping failed", err)
	}
	log.Println("Database connected successfully!")

	err = DB.AutoMigrate(&models.Guardian{}, &models.Student{}, &models.Teacher{})
	if err != nil {
		log.Fatal("Error running migrations!")
	}
	log.Println("Migrations run successful!")

	app := &handlers.App{Db: DB}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.HomeHandler)
	r.Route("/teacher", func(r chi.Router) {
		r.Get("/", app.GetAllTeachersHandler)
		r.Get("/admins", app.GetAllAdmins)
		r.Get("/{id}", app.GetATeacherHandler)
		r.Delete("/{id}", app.DeleteTeacherHandler)
		r.Put("/{id}", app.UpdateTeacherHandler)
		r.Post("/", app.CreateTeacherHandler)
		r.Put("/admin/{id}", app.MakeATeacherTheAdmin)
	})
	r.Route("/student", func(r chi.Router) {
		r.Get("/", app.GetAllStudents)
		r.Get("/{id}", app.GetStudent)
		r.Put("/{id}", app.UpdateAStudent)
		r.Delete("/{id}", app.DeleteAStudent)
		r.Post("/", app.CreateStudent)
	})

	server := &http.Server{
		Addr:    *port,
		Handler: r,
	}

	log.Println("*******************************************************")
	log.Println("Starting up a server for little angel api!")
	log.Println("server is running on port:", *port)
	log.Println("*******************************************************")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
