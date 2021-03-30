package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/simple-app/database"
	"github.com/simple-app/entity"
)

// connection postgresql const
const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "yourUserDB"  // default user: postgres
	password = "yourPassDB"   // you can set your own password
	dbname   = "yourNameDB" // your dbname
)

type handlerUser struct {
	db *sql.DB
}

// TrackerDuration execution
func TrackerDuration(start time.Time, funcName string) {
	duration := time.Since(start)
	log.Printf("%s took %s", funcName, duration)
}

// GetUserDetail from given param user id
func (h *handlerUser) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	//for tracking time execution
	startTime := time.Now()

	// get param request
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	userProfile, err := h.GetUserProfile(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	TrackerDuration(startTime, "GetUserProfile")

	userFamily, err := h.GetUserFamily(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	TrackerDuration(startTime, "GetUserFamily")

	userTransportation, err := h.GetUserTransportation(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	TrackerDuration(startTime, "GetUserTransportation")

	//aggregate data user detail
	userDetail := entity.UserDetail{}
	userDetail.Profile = userProfile
	userDetail.Family = userFamily
	userDetail.Transportation = userTransportation
	TrackerDuration(startTime, "All function")
	log.Println("==========================================")
	json.NewEncoder(w).Encode(userDetail)
	return
}

// GetUserProfile from DB
func (h *handlerUser) GetUserProfile(userID int64) (entity.UserProfile, error) {
	return database.GetUserProfile(userID, h.db)
}

// GetUserFamily from DB
func (h *handlerUser) GetUserFamily(userID int64) ([]entity.UserFamily, error) {
	return database.GetUserFamily(userID, h.db)
}

// GetUserTransportation from DB
func (h *handlerUser) GetUserTransportation(userID int64) ([]entity.UserTransportation, error) {
	return database.GetUserTransportation(userID, h.db)
}

func handleRequests(h *handlerUser) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/user/{id}", h.GetUserDetail)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	// create string conection for psql
	psql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// connect to DB
	dbconn, err := sql.Open("postgres", psql)
	if err != nil {
		log.Println(err)
		return
	}
	defer dbconn.Close()

	handler := &handlerUser{
		db: dbconn,
	}

	fmt.Println("===== || Starting Apps ||=====")
	handleRequests(handler)

}
