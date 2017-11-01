package main

import (
	"database/sql"
	"flag"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

var provinces []string
var dbPath *string

func getDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func getPhonesFromTable(phonenumber string, code string, db *sql.DB, table string) ([]Phone, error) {
	var rows *sql.Rows
	var err error
	if code == "" {
		rows, err = db.Query("select number, province, name, address from " + table + " where number = '" + phonenumber + "'")
	} else {
		rows, err = db.Query("select number, province, name, address from " + table + " where number = '" + phonenumber + "' and province = '" + code + "'")
	}

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

func appendIfNoError(phones []Phone, err error, result *[]Phone) {
	if err == nil && len(phones) > 0 {
		for _, amovil := range phones {
			*result = append(*result, amovil)
		}
	}
}

func getPhones(phonenumber string, db *sql.DB) ([]Phone, error) {
	result := make([]Phone, 0)
	movil, err := getPhonesFromTable(phonenumber, "", db, "movil")
	appendIfNoError(movil, err, &result)

	if strings.HasPrefix(phonenumber, "53") {
		movil, err := getPhonesFromTable(strings.Replace(phonenumber, "53", "", 1), "", db, "movil")
		appendIfNoError(movil, err, &result)
	}

	fix, err2 := getPhonesFromTable(phonenumber, "", db, "fix")
	appendIfNoError(fix, err2, &result)

	for _, code := range provinces {
		if strings.HasPrefix(phonenumber, code) {
			fix, err2 := getPhonesFromTable(strings.Replace(phonenumber, code, "", 1), code, db, "fix")
			appendIfNoError(fix, err2, &result)
		} else {
			if strings.HasPrefix(phonenumber, "0"+code) {
				fix, err2 := getPhonesFromTable(strings.Replace(phonenumber, "0"+code, "", 1), code, db, "fix")
				appendIfNoError(fix, err2, &result)
			}
		}
	}

	if len(result) == 0 {
		return []Phone{}, errors.New("Phone not found")
	} else {
		return result, nil
	}
}

func getProvinces(db *sql.DB) []string {
	log.Info("Loading provinces data")
	rows, err := db.Query("SELECT DISTINCT province FROM fix")
	if err != nil {
		panic(err)
	}

	provinces := make([]string, 0)

	for rows.Next() {
		var province string
		rows.Scan(&province)
		provinces = append(provinces, province)
	}
	log.Info("Done loading provinces data")
	return provinces
}

func handleSearch(c echo.Context) error {
	phonenumber := c.Param("phone")

	match, err := regexp.MatchString("^[0-9]+$", phonenumber)
	if phonenumber == "" || !match {
		return c.JSONPretty(http.StatusOK, PhoneArray{Status: "Phone not specified or no valid input", Phones: []Phone{}}, "    ")
	}
	db := getDB(*dbPath)
	phones, err := getPhones(phonenumber, db)
	if err == nil {
		return c.JSONPretty(http.StatusOK, PhoneArray{Status: "OK", Phones: phones}, "    ")
	} else {
		return c.JSONPretty(http.StatusNotFound, PhoneArray{Status: "Phone not found", Phones: []Phone{}}, "    ")
	}
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	dbPath = flag.String("db", "etecsa.db", "ETECSA's database path")
	port := flag.String("port", "6060", "port to be used by the service")

	flag.Parse()

	if !exists(*dbPath) {
		log.Error("Specified database does not exist")
		os.Exit(1)
	}

	provinces = getProvinces(getDB(*dbPath))
	e := echo.New()
	e.GET("/phones/:phone", handleSearch)
	e.File("/", "site/index.html")
	e.Static("/assets", "site/assets")
	e.Logger.Fatal(e.Start(":" + *port))
}
