package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerPhoneFromTable(t *testing.T) {
	db := getDB("/home/akiel/Desktop/etecsa.db")
	movil, err := getPhonesFromTable("54831077", "", db, "movil")
	if err == nil {
		assert.EqualValues(t, 1, len(movil), "found more or less than 1 result for unique phone")
		assert.EqualValues(t, "54831077", movil[0].Number, "result number does not match")
	} else {
		t.Error(err)
	}
}

func TestServerPhoneFromTableNotFound(t *testing.T) {
	db := getDB("/home/akiel/Desktop/etecsa.db")
	movil, err := getPhonesFromTable("000000", "", db, "movil")
	if err == nil {
		assert.EqualValues(t, 0, len(movil), "invalid amount of data returned for inexistent number")
	} else {
		t.Error(err)
	}
}

func TestAppendIfNoError(t *testing.T) {
	slice := make([]Phone, 0)
	sample := Phone{Number: "54831077", Province: "53", Name: "Jorge", Address: "none"}
	slice2 := []Phone{sample}
	appendIfNoError(slice2, nil, &slice)
	assert.EqualValues(t, sample, slice[0])
}
