package main

import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/labstack/echo"
	"net/http"
	"regexp"
)

type Phone struct {
	Name    string `json:"name" xml:"name" `
	Number  string `json:"number" xml:"number"`
	Address string `json:"address" xml:"address"`
	Status  string `json:"status" xml:"status"`
}

func getDB() (*sql.DB) {
	db, err := sql.Open("sqlite3", "/home/akiel/Desktop/etecsa.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func getPhoneFromTable(phonenumber string, db *sql.DB, table string) (Phone, error) {
	//ToDo: Add support for multiple results
	//ToDo: Add province field to result
	//ToDo: Try striping province code from number in order to search
	rows, err := db.Query("select number, name, address from " + table + " where number = '" + phonenumber + "'")
	if err != nil {
		return Phone{}, err
	}

	rows.Next()
	var Number string
	var Name string
	var Address string
	rows.Scan(&Number, &Name, &Address)

	if Number == "" {
		return Phone{}, errors.New("Phone not found on table " + table)
	}

	return Phone{Number: Number, Name: Name, Address: Address}, nil
}

func getPhone(phonenumber string, db *sql.DB) (Phone, error) {
	movil, err := getPhoneFromTable(phonenumber, db, "movil")
	if err == nil {
		return movil, nil
	}

	fix, err2 := getPhoneFromTable(phonenumber, db, "fix")
	if err2 == nil {
		return fix, nil
	}

	return Phone{}, errors.New("Phone not found")
}

func handleSearch(c echo.Context) error {
	phonenumber := c.Param("phone")

	match, err := regexp.MatchString("^[0-9]+$", phonenumber)
	if phonenumber == "" || !match{
		return c.JSONPretty(http.StatusOK, Phone{Status:"Phone not specified or no valid input", Number:phonenumber}, "    ")
	}
	db := getDB()
	phone, err := getPhone(phonenumber, db)
	if err == nil {
		phone.Status = "OK"
		return c.JSONPretty(http.StatusOK, phone, "    ")
	} else {
		return c.JSONPretty(http.StatusNotFound, Phone{Status:"Phone not found"}, "    ")
	}
}

func handleMain(c echo.Context) error {
	return c.String(http.StatusOK, "try curl http://"+c.Request().Host+"/phones/58999999\n")
}

func main() {
	e := echo.New()
	e.GET("/phones/:phone", handleSearch)
	//e.GET("/", handleMain)
	e.File("/", "site/index.html")
	e.Static("/assets", "site/assets")
	e.Logger.Fatal(e.Start(":6060"))
}
