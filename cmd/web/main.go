package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mihirkelkar/golang-jwt/pkg/models"
)

func main() {
	//get the port to run the server on from the command line argument
	addr := flag.String("port", ":4000", "Port on which to operate the server")
	conn := flag.String("postgres", "host=localhost port=5432 user=postgres "+
		"dbname=jwttest sslmode=disable", "dsn connection for postgres")
	flag.Parse()

	//define an error log and an info log.
	infoLog := log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\n", log.Ldate|log.Ltime)

	db, err := gorm.Open("postgres", *conn)
	if err != nil {
		errLog.Panic(err)
	}

	newUserService := models.NewUserService(db)

	app := Application{ErrLog: errLog, InfLog: infoLog, UserService: newUserService}
	app.UserService.AutoMigrate()

	srvr := http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srvr.ListenAndServe()
	errLog.Fatal(err)
}
