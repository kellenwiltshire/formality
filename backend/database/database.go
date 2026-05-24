package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Db is the database connection used by this package. It should be set by the
// application during initialization.
var Db *sql.DB

func ConnectDatabase() {

   err := godotenv.Load()//by default, it is .env so we don't have to write
   if err != nil {
      fmt.Println("Error is occurred  on .env file please check")
   }
   //we read our .env file
   host := os.Getenv("DB_HOST")
   port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // don't forget to convert int since port is int type.
   user := os.Getenv("DB_USERNAME")
   dbname := os.Getenv("DB_DATABASE_NAME")
   pass := os.Getenv("DB_PASSWORD")

   // set up postgres sql to open it.
   psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
       host, port, user, dbname, pass)
   db, errSql := sql.Open("postgres", psqlSetup)
   if errSql != nil {
      fmt.Println("There is an error while connecting to the database ", err)
      panic(err)
   } else {
      Db = db
      fmt.Println("Successfully connected to database!")
   }
}