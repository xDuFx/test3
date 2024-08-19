package models


type Token struct{
	AccessToken string
	RefreshToken string
}


type SessionUser struct{
	Id int
	RefreshToken string
	Ip string
	Guid string
}

type Configuration struct {
    LoginBd  string
    PassBd   string
	PortBd 	string
	NameBd string
	ServerPort string
	Email string
	EmailPass string
}
