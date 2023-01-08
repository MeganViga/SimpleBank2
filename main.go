package main

import (
	"database/sql"
	"log"

	"github.com/MeganViga/SimpleBank2/api"
	//"github.com/MeganViga/SimpleBank2/db"
	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	_ "github.com/lib/pq"
)

var (
	driver = "postgres"
	connectionString = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)
func main(){
	dbs, err := sql.Open(driver,connectionString)
	if err != nil{
		log.Fatal(err)
	}
	store := db.NewStore(dbs)
	apiServer := api.NewServer(store)
	apiServer.StartServer("0.0.0.0:8082")


	//db.Ping()
	//db.Close()
}