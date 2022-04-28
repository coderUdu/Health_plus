package handlers

import (
	// "encoding/json"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/decadev/squad10/healthplus/db"
	"github.com/decadev/squad10/healthplus/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var Sessions *scs.SessionManager

//Indexhandler gets the homepage
func Indexhandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/index.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// RegisterPatientHandler gets Patient's SignUp page
func RegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/registerPatient.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//PostRegisterPatientHandler successfully register's patient's name in the db if valid
func PostRegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Patient
	r.ParseForm()
	name := r.FormValue("name")
	ageString := r.FormValue("ageString")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	age, _ := strconv.Atoi(ageString)

	user.ID = uuid.NewString()
	user.Name = name
	user.Age = uint(age)
	user.Email = email
	user.Username = username
	user.Password = password

	_, err := db.FindUserByEmailandUserName(user.Email, user.Username)
	if err == nil {
		// this user already exists
		// return a message to the user

		t, e := template.ParseFiles("pages/registerPatient.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists. Check Email or Username")
		if e != nil {
			fmt.Println(e)
			return
		}

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)

	db.CreatePatientInTable(user)

	http.Redirect(w, r, "/patientLogin", http.StatusFound)

}

func PatientLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/patientLogin.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------PostPatientLoginHandler logs in doctor if valid-----------------------------------------------------
func PostLoginPatientdHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	_, err := db.AuthenticatePatient(user.Username, user.Password)
	if err != nil {
		t, e := template.ParseFiles("pages/patientLogin.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "Check username or Password")
		if e != nil {
			fmt.Println(e)
			return
		}
		return
	}
	Sessions.Put(r.Context(), "username", user.Username)
	http.Redirect(w, r, "/patientDashboard", http.StatusFound)
}

//------------------------------PatientDashboardHandler gets Patient's Dashboard page-----------------------------------------------
func PatientHomeHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/patientDashboard.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	patient, err := db.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, patient)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------PatientLogoutHandler logsout ---------------------------------------------------------------------
func PatientLogoutHandler(w http.ResponseWriter, r *http.Request) {
	Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

//-------------------------RegisterDoctorHandler gets Doctor's SignUp page-----------------------------------------------
func RegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorRegister.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//-------------------PostRegisterDoctorHandler successfully registers doctor's name in the db if valid----------------------------
func PostRegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor

	ageString := r.FormValue("ageString")
	age, _ := strconv.Atoi(ageString)

	user.ID = uuid.NewString()
	user.Name = strings.TrimSpace(r.FormValue("name"))
	user.Age = uint(age)
	user.Email = strings.TrimSpace(r.FormValue("email"))
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	user.Specialty = strings.TrimSpace(r.FormValue("specialty"))

	_, err := db.FindDocByEmailandUserName(user.Email, user.Username)
	if err == nil {

		t, e := template.ParseFiles("pages/doctorRegister.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists, confirm email or username")
		if e != nil {
			fmt.Println(e)
			return
		}

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)

	db.CreateDocInTable(user)

	http.Redirect(w, r, "/doctorLogin", http.StatusFound)

}

//------------------------------DoctorLoginHandler gets Doctor's Login page---------------------------------------------------------
func DoctorLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorLogin.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func PostLoginDoctordHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	usa, err := db.Authenticate(user.Username, user.Password)
	if err != nil {
		t, e := template.ParseFiles("pages/doctorLogin.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "Check Username or Password")
		if e != nil {
			fmt.Println(e)
			return
		}
		return
	}
	Sessions.Put(r.Context(), "username", usa.Username)
	http.Redirect(w, r, "/doctorDashboard", http.StatusFound)
}

//------------------------------DoctorDashboardHandler gets Doctor's Dashboard page-----------------------------------------------
func DoctorHomeHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorHome.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	doc, err := db.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, doc)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------DoctorLogoutHandler logsout ---------------------------------------------------------------------
func DoctorLogoutHandler(w http.ResponseWriter, r *http.Request) {
	Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

