package models

// Account contains all details for a account
type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Secret   string `json:"secret"`
}
