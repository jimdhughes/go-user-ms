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

type AccessTokenPayload struct {
	AccessToken string `json:"token"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refreshToken"`
}

func InitRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/login", HandleLogin).Methods("POST")
	r.HandleFunc("/register", HandleRegistration).Methods("POST")
	r.HandleFunc("/validateToken", HandleValidateToken).Methods("POST")
	r.HandleFunc("/refreshToken", HandleRefreshToken).Methods("POST")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	httpPort := os.Getenv("USERMS_HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}
	log.Panic(http.ListenAndServe(httpPort, loggedRouter))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error decoding user"))
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
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error decoding user"))
		return
	}
	success, err := DB.CreateUser(user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteResponse(w, http.StatusCreated, success)
}

func HandleValidateToken(w http.ResponseWriter, r *http.Request) {
	var Token AccessTokenPayload
	err := json.NewDecoder(r.Body).Decode(&Token)
	if err != nil {
		fmt.Fprintf(w, "Error Decoding Body: %v", err.Error())
		return
	}
	payload, err := TS.ValidateAccessToken(Token.AccessToken)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteResponse(w, http.StatusOK, payload)
}

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var Token AccessTokenPayload
	err := json.NewDecoder(r.Body).Decode(&Token)
	if err != nil {
		fmt.Fprintf(w, "Error Decoding Body: %v", err.Error())
		return
	}
	payload, err := TS.ValidateRefreshToken(Token.AccessToken)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := DB.GetUserById(payload)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := TS.GenerateTokenPairForUser(user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	WriteResponse(w, http.StatusOK, resp)
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
