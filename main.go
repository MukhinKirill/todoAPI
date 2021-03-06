package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	todos "github.com/MukhinKirill/todo-api/lib"
)

func main() {
	file, errf := os.Open("config.json")
	if errf != nil {
		panic(errf)
	}
	decoder := json.NewDecoder(file)
	config := new(todos.Config)
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	var postgres *todos.Postgres
	var errDb error
	for i := 0; i < 3; i++ {
		time.Sleep(3 * time.Second)
		postgres, errDb = todos.ConnectDb(config.ConnectionString)
	}
	if errDb != nil {
		panic(errDb)
	} else if postgres == nil {
		panic("postgres is nil")
	}
	defer postgres.Close()
	postgres.DbInit()
	router := todos.SetUpRouting(postgres)
	portStr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("http://localhost%s", portStr)

	log.Fatal(http.ListenAndServe(portStr, router))
}
