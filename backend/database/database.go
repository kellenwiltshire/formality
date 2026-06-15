package database

import (
	"database/sql"
	"fmt"
	loadenv "formality/backend/load_env"
	"strconv"

	_ "github.com/lib/pq"
)

// Db is the database connection used by this package. It should be set by the
// application during initialization.
var Db *sql.DB

func ConnectDatabase() {

   host := loadenv.LoadDotEnvVariable("DB_HOST")
   port, _ := strconv.Atoi(loadenv.LoadDotEnvVariable("DB_PORT"))
   user := loadenv.LoadDotEnvVariable("DB_USERNAME")
   dbname := loadenv.LoadDotEnvVariable("DB_DATABASE_NAME")
   pass := loadenv.LoadDotEnvVariable("DB_PASSWORD")

   // set up postgres sql to open it.
   psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
       host, port, user, dbname, pass)
   db, errSql := sql.Open("postgres", psqlSetup)
   if errSql != nil {
      fmt.Println("There is an error while connecting to the database ", errSql)
      panic(errSql)
   } else {
      Db = db
      fmt.Println("Successfully connected to database!")
   }
}