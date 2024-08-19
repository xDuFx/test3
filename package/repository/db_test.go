package repository_test

import (
	"encoding/base64"
	"test3/package/repository"
	"testing"
)

func Test_CRUD(t *testing.T) {
	test_table := []struct {
		guid         string
		ip           string
		refreshToken string
		valid        bool
	}{{
		guid:         "1452asd145",
		ip:           "127.0.0.1:64103",
		refreshToken: "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBjQ0k2SWpFeU55NHdMakF1TVNJc0luTjFZaUk2SWpFME5USmhjMlF4TkRVaWZRLjRyQzVtUi1yc0Y4OFJWaG1LbDRWSllqZXg4X19jb2x2YWJsd3p5cFdMbVhSQkxWc3BxSHc2SG00aHBBSjRUaFZuMjA5VGp0enh5SUQ0RXB0elNGMjdn",
		valid:        true,
	}, {
		guid:         "23123asddq",
		ip:           "",
		refreshToken: "",
		valid:        false,
	},
	}
	db, err := repository.New("postgres://postgres:123@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range test_table {
		flag, err := db.CheckGuid(tc.guid)
		if err != nil && err.Error() != "no rows in result set" {
			t.Fatal(err)
		}
		if flag != tc.valid {
			t.Fatal("Неправильно отработал")
		}
		flag, err = db.CheckRefresh(tc.refreshToken)
		if err != nil && err.Error() != "token contains an invalid number of segments" {
			t.Fatal(err)
		}
		if flag != tc.valid {
			t.Fatal("Неправильно отработал")
		}
		rtoken, err := base64.StdEncoding.DecodeString(tc.refreshToken)
		if err != nil {
			t.Fatal(err)
		}
		id, err := db.Create(tc.guid, string(rtoken), tc.ip)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(id, tc.guid)
		email, err := db.Email(tc.guid)
		if err != nil && err.Error() != "no rows in result set" {
			t.Fatal(err)
		}
		t.Log(email)
		ip, err := db.IpCheck(tc.guid)
		if err != nil && err.Error() != "no rows in result set" {
			t.Fatal(err)
		}
		if ip == "" && tc.valid {
			t.Fatal("Неправильный ip")
		}
		id, err = db.EmailMark(tc.guid, tc.ip)
		if err != nil && err.Error() != "no rows in result set"  {
			t.Fatal(err)
		}
		t.Log(id)
	}
}
