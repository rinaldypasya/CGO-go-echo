package models

type GetToken struct {
	AccessToken	string 	`json:"access_token"`
	ExpiresIn	int	`json:"expires_in"`
	TokenType	string	`json:"token_type"`
}

type RegisterAndUpdateUser struct {
	Id 			string	`json:"id"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Name 		string	`json:"name"`
	GivenName	string	`json:"givenname"`
	FamilyName	string	`json:"familyname"`
	Email 		string	`json:"email"`
	EmailVerified	bool	`json:"emailverified"`
	Website		string	`json:"website"`
	Address 	string	`json:"address"`
}

type GetUserInfo struct {
	Id 			string	`json:"id"`
	Username	string	`json:"username"`
	Name 		string	`json:"name"`
	GivenName	string	`json:"givenname"`
	FamilyName	string	`json:"familyname"`
	Email 		string	`json:"email"`
	EmailVerified	bool	`json:"emailverified"`
	Website		string	`json:"website"`
	Address 	string	`json:"address"`
}

type Login struct {
	Email				string		`json:"email"`
	Password			string		`json:"password"`
	Type 				string		`json:"type"`
}
