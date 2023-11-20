package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (u *User) Prepare(stage string) error {
	if err := u.validate(stage); err != nil {
		return err
	}
	if err := u.format(stage); err != nil {
		return err
	}
	return nil
}

func (u *User) validate(stage string) error {
	if u.Name == "" {
		return errors.New("campo nome é obrigatório")
	}

	if u.Nick == "" {
		return errors.New("o campo nick é obrigatório")
	}

	if u.Email == "" {
		return errors.New("o campo email é obrigatório")
	}

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("o email inserido é inválido")
	}

	if stage == "register" && u.Password == "" {
		return errors.New("o campo senha é obrigatório")
	}

	return nil
}

func (u *User) format(stage string) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Nick = strings.TrimSpace(u.Nick)
	u.Email = strings.TrimSpace(u.Email)

	if stage == "register" {
		hashedPassword, err := security.Hash(u.Password)
		if err != nil{
			return err
		}

		u.Password = string(hashedPassword)
	}

	return nil
}