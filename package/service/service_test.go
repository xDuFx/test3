package service_test

import (
	"encoding/base64"
	"test3/package/service"
	"testing"
)

type test_data struct {
	guid string
	ip   string
}

func Test_token(t *testing.T) {
	test1 := test_data{guid: "15asd15ad2", ip: "144.155.120"}
	tokens, err := service.CreateToken(test1.guid, test1.ip)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Были созданы токены, аксесс: %s \n рефреш: %s \n", tokens.AccessToken, tokens.RefreshToken)
	rtoken := base64.StdEncoding.EncodeToString([]byte(tokens.RefreshToken))
	atoken := base64.StdEncoding.EncodeToString([]byte(tokens.AccessToken))
	accessClaims, err := service.ParseAccessToken(atoken)
	if err != nil {
		t.Error(err)
	}
	flag, err := service.CompareToken(atoken, rtoken)
	if err != nil {
		t.Error(err.Error())
	}
	if !flag{
		t.Error("Ошибка сравнения")
	}
	t.Logf("Значения из акцесс токена \n guid: %s, ip: %s \n", accessClaims["sub"], accessClaims["ip"])
	refClaims, err := service.ParseRefreshToken(rtoken)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Значения из рефреш токена \n guid: %s, ip: %s \n", refClaims["sub"], refClaims["ip"])
	// проверка bcrypt
	part1 := tokens.RefreshToken[:70]
	part2 := tokens.RefreshToken[70:140]
	part3 := tokens.RefreshToken[140:]
	hash1, err := service.HashToken(part1)
	if err != nil {
		t.Error("Ошибка создания хэша: ", err.Error())
	}
	hash2, err := service.HashToken(part2)
	if err != nil {
		t.Error("Ошибка создания хэша: ", err.Error())
	}
	hash3, err := service.HashToken(part3)
	if err != nil {
		t.Error("Ошибка создания хэша: ", err.Error())
	}
	allHash := hash1 + hash2 + hash3
	check := service.CheckTokenHash(tokens.RefreshToken[:70], allHash[:60])
	if check {
		t.Log("Проверка хэша пройдена")
	} else {
		t.Error("Проверка хэша не пройдена")
	}
	check = service.CheckTokenHash(tokens.RefreshToken[70:140], allHash[60:120])
	if check {
		t.Log("Проверка хэша пройдена")
	} else {
		t.Error("Проверка хэша не пройдена")
	}
	check = service.CheckTokenHash(tokens.RefreshToken[140:], allHash[120:])
	if check {
		t.Log("Проверка хэша пройдена")
	} else {
		t.Error("Проверка хэша не пройдена")
	}
}
