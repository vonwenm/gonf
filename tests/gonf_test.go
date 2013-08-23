package tests

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"github.com/xrash/gonf"
)

const (
	teststring = `
# comment 1
# comment 2

testinteger 22

awesome [
    arrays
    "like this"
    # comment inside an array
]

"keys are any string literal" "just as values"

username www-data
group www-data

anything "any thing"
backslashed "wa\\"
escaped "\""

database {
    host 127.0.0.1
    auth {
        # comment over here
        user testuser
        pass testpass
    }
}

# finish it up`
)

type teststruct struct {
	Username string `gonf:"username"`
	Group string `gonf:"group"`
	Database struct {
		Host string `gonf:"host"`
		Auth struct {
			User string `gonf:"user"`
			Pass string `gonf:"pass"`
		} `gonf:"auth"`
	} `gonf:"database"`
	TestInteger int `gonf:"testinteger"`
}

var file *os.File
var config *gonf.Config

func createTempFile(t *testing.T) {
	var err error
	file, err = ioutil.TempFile(os.TempDir(), "__gonf_test_")
	if err != nil {
		t.FailNow()
	}

	if n, err := file.WriteString(teststring); err != nil || n != len(teststring) {
		t.FailNow()
	}

	file.Seek(0, 0)
}

func deleteTempFile(t *testing.T) {
	if err := file.Close(); err != nil {
		t.FailNow()
	}
	if err := os.Remove(file.Name()); err != nil {
		t.FailNow()
	}
}

func readTempFile(t *testing.T) {
	var err error
	config, err = gonf.Read(file)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func testString(t *testing.T) {
	username, _ := config.String("username")
	group, _ := config.String("group")
	anything, _ := config.String("anything")
	backslashed, _ := config.String("backslashed")
	escaped, _ := config.String("escaped")
	host, _ := config.String("database", "host")
	user, _ := config.String("database", "auth", "user")
	arrays, _ := config.String("awesome", 0)
	likethis, _ := config.String("awesome", 1)
	wololo, _ := config.String("keys are any string literal")
	testinteger, _ := config.Int("testinteger")

	stringtests := map[string]string{
		username: "www-data",
		group: "www-data",
		anything: "any thing",
		backslashed: "wa\\",
		escaped: "\"",
		host: "127.0.0.1",
		user: "testuser",
		arrays: "arrays",
		likethis: "like this",
		wololo: "just as values",
	}

	inttests := map[int]int{
		testinteger: 22,
	}

	for i, v := range stringtests {
		if i != v {
			t.Fail()
		}
	}

	for i, v := range inttests {
		if i != v {
			t.Fail()
		}
	}
}

func testMap(t *testing.T) {
	s := new(teststruct)

	config.Map(s)

	stringtests := map[string]string{
		s.Username: "www-data",
		s.Group: "www-data",
		s.Database.Host: "127.0.0.1",
		s.Database.Auth.User: "testuser",
	}

	inttests := map[int]int{
		s.TestInteger: 22,
	}

	for i, v := range stringtests {
		if i != v {
			t.Fail()
		}
	}

	for i, v := range inttests {
		if i != v {
			t.Fail()
		}
	}
}

func TestAll(t *testing.T) {
	createTempFile(t)
	defer deleteTempFile(t)
	readTempFile(t)

	testString(t)
	testMap(t)
}