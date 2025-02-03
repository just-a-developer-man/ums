package models

type User struct {
	ID       int    `json:"id,omitempty" validate="numeric"`
	Name     string `json:"name,omitempty" validate="required,alphanum"`
	Email    string `json:"email,omitempty" validate="required,email`
	Password string `json:"password,omitempty" validate="required,min=6,max=64,printascii"`
}
