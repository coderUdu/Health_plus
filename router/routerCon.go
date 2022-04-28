package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/decadev/squad10/healthplus/handlers"
	"github.com/gorilla/mux"
	"os"
)

func SetupRouter() {
  handlers.Sessions  = scs.New()
	handlers.Sessions.Lifetime = 24 * time.Hour
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Indexhandler).Methods("GET")
	router.HandleFunc("/registerPatient", handlers.RegisterPatientHandler).Methods("GET")
	router.HandleFunc("/registerPatient", handlers.PostRegisterPatientHandler).Methods("POST")
	router.HandleFunc("/patientLogin", handlers.PatientLoginHandler).Methods("GET")
	router.HandleFunc("/patientLogin", handlers.PostLoginPatientdHandler).Methods("POST")
	router.HandleFunc("/registerDoctor", handlers.RegisterDoctorHandler).Methods("GET")
	router.HandleFunc("/registerDoctor", handlers.PostRegisterDoctorHandler).Methods("POST")
	router.HandleFunc("/doctorLogin", handlers.DoctorLoginHandler).Methods("GET")
	router.HandleFunc("/doctorLogin", handlers.PostLoginDoctordHandler).Methods("POST")
 	router.HandleFunc("/doctorLogout", handlers.DoctorLogoutHandler).Methods("GET")
  	router.HandleFunc("/doctorDashboard", handlers.DoctorHomeHandler).Methods("GET")
	router.HandleFunc("/patientDashboard", handlers.PatientHomeHandler).Methods("GET")
	router.HandleFunc("/patientLogout", handlers.PatientLogoutHandler).Methods("GET")
	

	fs := http.FileServer(http.Dir("./pages/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//http.Handle("/", router)

	fmt.Println("Server up and Running")

	port := os.Getenv("DB_PORT")



	e := http.ListenAndServe(port, handlers.Sessions.LoadAndSave(router))

// 	e := http.ListenAndServe(port, router)

	if e != nil {
		fmt.Println(e)
		return
	}
}