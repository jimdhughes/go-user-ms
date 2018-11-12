package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ApiErrorResponse struct {
	Error string `json:"error"`
}

type ApiStandardResponse struct {
	Payload interface{} `json:"payload"`
}

type TokenPayload struct {
	Token string `json:"token"`
}

func InitRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/login", HandleLogin).Methods("POST")
	r.HandleFunc("/register", HandleRegistration).Methods("POST")
	r.HandleFunc("/validateToken", HandleValidateToken).Methods("POST")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Panic(http.ListenAndServe(":8080", loggedRouter))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error Decoding User"))
	}
	status, err := DB.Login(user.Email, user.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteResponse(w, http.StatusOK, status)
}

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error Decoding User"))
	}
	success, err := DB.CreateUser(user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteResponse(w, http.StatusCreated, success)
}

func HandleValidateToken(w http.ResponseWriter, r *http.Request) {
	var Token TokenPayload
	err := json.NewDecoder(r.Body).Decode(&Token)
	if err != nil {
		fmt.Fprintf(w, "Error Decoding Body: %v", err.Error())
		return
	}
	payload, err := TS.ValidateToken(Token.Token)
	if err != nil {
		fmt.Fprintf(w, "Error Decoding Token: %v", err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, payload)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	response := ApiErrorResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(response)
}

func WriteResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	response := ApiStandardResponse{Payload: payload}
	json.NewEncoder(w).Encode(response)
}
