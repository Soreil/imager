package webm

import (
	"io/ioutil"
	"os"
	"testing"
)

const dataDirectory = "../testData/"

func TestWebm(t *testing.T) {
	const filename = dataDirectory + "wafel.webm"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	decode(file)
}
