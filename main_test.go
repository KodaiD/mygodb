package main

import (
	"bytes"
	"testing"
)

func TestDB(t *testing.T) {
	testString := "insert 1 user1 person1@example.com\ninsert 2 user2 person2@example.com\nunrecognized\nselect\n.exit"
	input := bytes.NewBuffer([]byte(testString))
	output := &bytes.Buffer{}
	db := newDB(input)
	db.run(output)
	if output.String() != testString {
		t.Errorf("Wrong output: %v", output.String())
	}
}