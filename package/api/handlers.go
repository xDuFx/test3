package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"test3/package/sendemail"
	"test3/package/service"

	"github.com/gorilla/mux"
)

func (api *api) auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		guid, ok := vars["guid"]
		if !ok {
			http.Error(w, "No guid parameter", http.StatusInternalServerError)
			return
		}
		check, err := api.db.CheckGuid(guid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !check {
			fmt.Fprintln(w, "Нет данных по пользователю")
			return
		}
		splitIp := strings.Split(r.RemoteAddr, ":")
		token, err := service.CreateToken(guid, splitIp[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = api.db.Create(guid, token.RefreshToken, splitIp[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rtoken := base64.StdEncoding.EncodeToString([]byte(token.RefreshToken))
		atoken := base64.StdEncoding.EncodeToString([]byte(token.AccessToken))
		w.Header().Set("Authorization", fmt.Sprintf("{Access Token: %v}, {refresh Token: %v}", atoken, rtoken))

	}
}

func (api *api) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var tokenAccessHash struct {
			Authorization string `json:""Authorization"`
		}
		if err := json.NewDecoder(r.Body).Decode(&tokenAccessHash); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		vars := mux.Vars(r)
		refreshHash, ok := vars["refresh"]
		if !ok {
			http.Error(w, "Нет параметра", http.StatusInternalServerError)
			return
		}
		checkUnion, err := service.CompareToken(tokenAccessHash.Authorization, refreshHash)
		if err != nil {
			http.Error(w, "Ошибка в сравнении токенов", http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
		if !checkUnion {
			http.Error(w, "Неправильная пара токенов", http.StatusInternalServerError)
		}
		claims, err := service.ParseRefreshToken(refreshHash)
		if err != nil {
			http.Error(w, "Ошибка в проверке токена", http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
		guid, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Ошибка в проверке токена", http.StatusInternalServerError)
			log.Println("Ошибка guid")
			return
		}
		check, err := api.db.CheckGuid(guid)
		if err != nil {
			http.Error(w, "Ошибка в проверке токена", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		if !check {
			fmt.Fprintln(w, "Нет данных по пользователю")
			return
		}
		check, err = api.db.CheckRefresh(refreshHash)
		if err != nil {
			http.Error(w, "Ошибка в проверке токена", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		if !check {
			fmt.Fprintln(w, "Нет данных по токену")
		}

		ipCheck, err := api.db.IpCheck(guid)

		if err != nil {
			http.Error(w, "Ошибка в проверки айпи", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		if ipCheck != "" {
			email, err := api.db.Email(guid)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			err = sendemail.Emailsend(email, ipCheck)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
		}
		splitIp := strings.Split(r.RemoteAddr, ":")
		token, err := service.CreateToken(guid, splitIp[0])
		if err != nil {
			http.Error(w, "Ошибка в создании токена", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		_, err = api.db.Update(guid, token.RefreshToken, splitIp[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rtoken := base64.StdEncoding.EncodeToString([]byte(token.RefreshToken))
		atoken := base64.StdEncoding.EncodeToString([]byte(token.AccessToken))
		w.Header().Set("Authorization", fmt.Sprintf("Access Token: %v, refresh Token: %v", atoken, rtoken))

	}
}
