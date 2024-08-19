package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"test3/package/service"
)

func (repo *PGRepo) CheckGuid(guid string) (bool bool, err error) {
	var checkguid string
	err = repo.pool.QueryRow(context.TODO(), `
	SELECT guid FROM Users
	WHERE guid = $1;`,
		guid,
	).Scan(&checkguid)
	if checkguid == "" || err != nil {
		return false, err
	}
	return true, nil
}

func (repo *PGRepo) Create(guid, refresh, ip string) (id int, err error) {
	if len(refresh) < 150 {
		return id, nil
	}
	part1 := refresh[:70]
	part2 := refresh[70:140]
	part3 := refresh[140:]
	hash1, err := service.HashToken(part1)
	if err != nil {
		return id, err
	}
	hash2, err := service.HashToken(part2)
	if err != nil {
		return id, err
	}
	hash3, err := service.HashToken(part3)
	if err != nil {
		return id, err
	}
	err = repo.pool.QueryRow(context.TODO(), `
	INSERT INTO SessionUsers(refreshToken, guid) VALUES
	($1, $2)
	returning id;`,
		hash1 + hash2 + hash3,
		guid,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
func (repo *PGRepo) Update(guid, refresh, ip string) (id int, err error) {
	if len(refresh) < 150 {
		return id, nil
	}
	part1 := refresh[:70]
	part2 := refresh[70:140]
	part3 := refresh[140:]
	hash1, err := service.HashToken(part1)
	if err != nil {
		return id, err
	}
	hash2, err := service.HashToken(part2)
	if err != nil {
		return id, err
	}
	hash3, err := service.HashToken(part3)
	if err != nil {
		return id, err
	}
	err = repo.pool.QueryRow(context.TODO(), `
	UPDATE SessionUsers
	SET refreshToken = $1
	WHERE guid = $2
	returning id;`,
		hash1 + hash2 + hash3,
		guid,
	).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (repo *PGRepo) CheckRefresh(brefresh string) (bool, error) {
	ref, err := base64.StdEncoding.DecodeString(brefresh)
	if err != nil {
		return false, err
	}
	claims, err :=service.ParseRefreshToken(brefresh)
	if err != nil {
		return false, err
	}
	guid := claims["sub"].(string)
	rows, err := repo.pool.Query(context.TODO(), `
	SELECT refreshToken FROM SessionUsers
	WHERE guid = $1;`,
		guid,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next(){
		var hashToken string
		rows.Scan(
			&hashToken,
		)
		if service.CheckTokenHash(string(ref)[:70], hashToken[:60]) && service.CheckTokenHash(string(ref)[70:140], hashToken[60:120]) && service.CheckTokenHash(string(ref)[140:], hashToken[120:]){
			return true, nil
		}
	}
	return false, nil
}

func (repo *PGRepo) EmailMark(guid, ip string) (id int, err error){
	err = repo.pool.QueryRow(context.TODO(), `
	UPDATE SessionUsers
	SET ip = $2
	WHERE guid = $1
	returning id;`,
		guid,
		ip,
	).Scan(&id)
	return id, err
}

func (repo *PGRepo) IpCheck(guid string) (string, error){
	var flag sql.NullString
	err := repo.pool.QueryRow(context.TODO(), `
	SELECT ip FROM SessionUsers
	WHERE guid = $1;`,
		guid,
	).Scan(&flag)
	if flag.Valid {
		return flag.String, err
	}
	return "", err
}


func (repo *PGRepo) Email(guid string) (string,  error) {
	var email sql.NullString
	err := repo.pool.QueryRow(context.TODO(), `
	SELECT email FROM Users
	WHERE guid = $1;`,
		guid,
	).Scan(&email)
	if email.Valid {
		return email.String, err
	}
	return "", err
}