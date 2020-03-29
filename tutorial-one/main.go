package main

import (
	"fmt"
	"net/http"

	"github.com/cbigge/go-web/tutorial-one/controllers"
	"github.com/cbigge/go-web/tutorial-one/middleware"
	"github.com/cbigge/go-web/tutorial-one/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "a9725770137"
	dbname   = "example_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	indexGallery := requireUserMw.ApplyFn(galleriesC.Index)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)
	newGallery := requireUserMw.Apply(galleriesC.New)
	editGallery := requireUserMw.ApplyFn(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFn(galleriesC.Update)
	deleteGallery := requireUserMw.ApplyFn(galleriesC.Delete)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	r.Handle("/galleries", indexGallery).Methods("GET").
		Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").
		Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")

	http.ListenAndServe(":4000", userMw.Apply(r))
}
