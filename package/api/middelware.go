package api

import (

	"net/http"
	"strings"
	"test3/package/service"

)

func (api *api) IpCheck(check_handle func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		token := r.FormValue("Authorization")

		if token == "" {
			http.Error(w, "Нет токена", http.StatusForbidden)
			return
		}
		accessClaims, err := service.ParseAccessToken(token)
		ip := strings.Split(r.RemoteAddr, ":")
		if accessClaims["ip"] != ip[0]{
			guid := accessClaims["sub"].(string)
			ip := accessClaims["ip"].(string)
			api.db.EmailMark(guid, ip)
		}
		if err != nil{
			http.Error(w, "Нет токена", http.StatusForbidden)
			return
		}


		check_handle(w,r)
	}
}