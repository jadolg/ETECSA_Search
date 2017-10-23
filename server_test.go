package main

import "testing"

func TestServerPhoneFromTable(t *testing.T) {
	db := getDB()
	movil, err := getPhoneFromTable("54831077", db, "movil")
	if err == nil {
		t.Log(movil)
	}
}