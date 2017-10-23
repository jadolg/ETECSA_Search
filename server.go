package main

import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/labstack/echo"
	"net/http"
	"regexp"
)

type PhoneArray struct {
	Phones []Phone `json:"phones" xml:"phones"`
	Status string  `json:"status" xml:"status"`
}

type Phone struct {
	Name     string `json:"name" xml:"name" `
	Number   string `json:"number" xml:"number"`
	Province string `json:"province"xml:"province"`
	Address  string `json:"address" xml:"address"`
}

func getDB() (*sql.DB) {
	db, err := sql.Open("sqlite3", "/home/akiel/Desktop/etecsa.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func getPhoneFromTable(phonenumber string, db *sql.DB, table string) ([]Phone, error) {
	//ToDo: Try striping province code from number in order to search
	rows, err := db.Query("select number, province, name, address from " + table + " where number = '" + phonenumber + "'")
	if err != nil {
		return []Phone{}, err
	}

	phones := make([]Phone, 0)

	for rows.Next() {
		var Number string
		var Province string
		var Name string
		var Address string
		rows.Scan(&Number, &Province, &Name, &Address)

		if Number == "" {
			return []Phone{}, errors.New("Phone not found on table " + table)
		}
		if table == "movil" {
			phones = append(phones, Phone{Number: Number, Province: "53", Name: Name, Address: Address})
		} else {
			phones = append(phones, Phone{Number: Number, Province: Province, Name: Name, Address: Address})
		}

	}

	return phones, nil
}

func getPhones(phonenumber string, db *sql.DB) ([]Phone, error) {
	movil, err := getPhoneFromTable(phonenumber, db, "movil")
	if err == nil && len(movil) > 0 {
		return movil, nil
	}

	fix, err2 := getPhoneFromTable(phonenumber, db, "fix")
	if err2 == nil && len(fix) > 0 {
		return fix, nil
	}

	return []Phone{}, errors.New("Phone not found")
}

func handleSearch(c echo.Context) error {
	phonenumber := c.Param("phone")

	match, err := regexp.MatchString("^[0-9]+$", phonenumber)
	if phonenumber == "" || !match {
		return c.JSONPretty(http.StatusOK, PhoneArray{Status: "Phone not specified or no valid input", Phones: []Phone{}}, "    ")
	}
	db := getDB()
	phones, err := getPhones(phonenumber, db)
	if err == nil {
		return c.JSONPretty(http.StatusOK, PhoneArray{Status: "OK", Phones: phones}, "    ")
	} else {
		return c.JSONPretty(http.StatusNotFound, PhoneArray{Status: "Phone not found", Phones: []Phone{}}, "    ")
	}
}

func handleMain(c echo.Context) error {
	return c.String(http.StatusOK, "try curl http://"+c.Request().Host+"/phones/58999999\n")
}

func main() {
	e := echo.New()
	e.GET("/phones/:phone", handleSearch)
	e.File("/", "site/index.html")
	e.Static("/assets", "site/assets")
	e.Logger.Fatal(e.Start(":6060"))
}
