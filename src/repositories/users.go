package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type Users struct {
	db *sql.DB
}

func NewUserRep(db *sql.DB) *Users {
	return &Users{db}
}

func (u Users) Create(user models.User) (uint64, error) {
	sql, err := u.db.Prepare("INSERT INTO users (name, nick, email, password) VALUES(?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer sql.Close()

	result, err := sql.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

func (u Users) Search(value string) ([]models.User, error) {
	newValue := fmt.Sprintf("%%%s%%", value)

	sql, err := u.db.Query("select id, name, nick, email, created_at from users where name LIKE ? or nick LIKE ?", newValue, newValue)

	if err != nil {
		return nil, err
	}
	defer sql.Close()

	var users []models.User

	for sql.Next(){
		var user models.User
		if err = sql.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u Users) GetById(id uint64) (models.User, error) {
	sql, err := u.db.Query("select id, name, nick, email, created_at from users where id = ?", id)
	if err != nil {
		return models.User{}, err
	}
	defer sql.Close()

	var user models.User

	if sql.Next() {
		if err := sql.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (u Users) Update(id uint64, user models.User) error{
	sql, err := u.db.Prepare("UPDATE users SET name = ?, email = ?, nick = ? where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _ , err := sql.Exec(user.Name, user.Nick, user.Email, id); err != nil {
		return err
	}
	return nil
}

func (u Users) Delete(id uint64) error{
	sql, err := u.db.Prepare("DELETE from users where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _ , err := sql.Exec(id); err != nil {
		return err
	}
	return nil
}

func (u Users) SearchByEmail(email string) (models.User, error) {
	sql, err := u.db.Query("SELECT id, password from users where email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer sql.Close()

	var user models.User

	if sql.Next() {
		if err := sql.Scan(
			&user.ID,
			&user.Password,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (u Users) Follow(userID, followerID uint64) error{
	sql, err := u.db.Prepare("INSERT ignore INTO followers(user_id, follower_id) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _ , err := sql.Exec(userID, followerID); err != nil {
		return err
	}

	return nil	
}

func (u Users) StopFollowing(userID, followerID uint64) error{
	sql, err := u.db.Prepare("DELETE from followers where user_id = ? and follower_id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _ , err := sql.Exec(userID, followerID); err != nil {
		return err
	}

	return nil	
}

func (u Users) GetFollowersById(userID uint64) ([]models.User, error){
	sql, err := u.db.Query("select u.id, u.name, u.nick, u.email, u.created_at from users u inner join followers f on u.id = f.follower_id where f.user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer sql.Close()

	var users []models.User

	for sql.Next(){
		var user models.User
		if err = sql.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u Users) GetFollowing(userID uint64) ([]models.User, error){
	sql, err := u.db.Query("select u.id, u.name, u.nick, u.email, u.created_at from users u inner join followers f on u.id = f.user_id where f.follower_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer sql.Close()

	var users []models.User

	for sql.Next(){
		var user models.User
		if err = sql.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u Users) GetCurrentPassword(userID uint64) (string, error){
	sql, err := u.db.Query("select password from users where id = ?", userID)
	if err != nil {
		return "", err
	}
	defer sql.Close()

	var user models.User
	if sql.Next() {
		if err = sql.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (u Users) UpdatePassword(userID uint64, newPassword string) error{
	sql, err := u.db.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _ , err := sql.Exec(newPassword, userID); err != nil {
		return err
	}
	return nil
}