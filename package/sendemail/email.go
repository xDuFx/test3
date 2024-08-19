package sendemail

import (
	"encoding/json"
	"net/smtp"
	"os"
	"test3/package/models"
)

func Emailsend(email, ip string) error {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := models.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return err
	}
	from := configuration.Email
	password := configuration.EmailPass

	// Информация о получателе
	to := []string{
		email,
	}

	// smtp сервер конфигурация
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Сообщение.

	message := []byte("На ваш аккаунт зашли с другого ip " + ip)
	
	// Авторизация.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка почты.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return err
}