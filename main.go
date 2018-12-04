package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	InitializeDatabase()
	InitializeTokenService()
	InitializeHttpServer()
}

func InitializeHttpServer() {
	InitRouter()
}

func InitializeDatabase() {
	DB = &DBClient{}
	DB.Initialize("/app/bolt.db")
}

func InitializeTokenService() {
	TS = &TokenService{}
}
