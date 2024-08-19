package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"test3/package/service"
)

func (api *api) IpCheck(check_handle func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenAccessHash struct {
			Authorization string `json:""Authorization"`
		}
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if err := json.NewDecoder(r.Body).Decode(&tokenAccessHash); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if tokenAccessHash.Authorization == "" {
			http.Error(w, "Нет токена", http.StatusForbidden)
			return
		}
		accessClaims, err := service.ParseAccessToken(tokenAccessHash.Authorization)
		ip := strings.Split(r.RemoteAddr, ":")
		if accessClaims["ip"] != ip[0] {
			guid := accessClaims["sub"].(string)
			ip := accessClaims["ip"].(string)
			api.db.EmailMark(guid, ip)
		}
		if err != nil {
			http.Error(w, "Нет токена", http.StatusForbidden)
			return
		}

		check_handle(w, r)
	}
}
