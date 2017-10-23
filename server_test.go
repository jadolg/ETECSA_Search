package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestServerPhoneFromTable(t *testing.T) {
	db := getDB()
	movil, err := getPhonesFromTable("54831077", db, "movil")
	if err == nil {
		assert.EqualValues(t, 1, len(movil), "found more or less than 1 result for unique phone")
		assert.EqualValues(t, "54831077", movil[0].Number, "result number does not match")
	} else {
		t.Error(err)
	}
}

func TestServerPhoneFromTableNotFound(t *testing.T) {
	db := getDB()
	movil, err := getPhonesFromTable("000000", db, "movil")
	if err == nil {
		assert.EqualValues(t, 0, len(movil), "invalid amount of data returned for inexistent number")
	} else {
		t.Error(err)
	}
}
