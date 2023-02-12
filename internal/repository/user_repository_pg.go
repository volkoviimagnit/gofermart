package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserRepositoryPG struct {
	conn *db.ConnectionPostgres
}

func (u *UserRepositoryPG) Insert(user model.User) error {
	sqlRequest := `INSERT INTO public."user" (id, login, password, token) VALUES (DEFAULT, $1, $2, $3);`

	errConnecting := u.conn.TryConnect()
	if errConnecting != nil {
		return errConnecting
	}

	var userToken sql.NullString
	errExecuting := u.conn.Exec(sqlRequest, user.GetLogin(), user.GetPassword(), userToken)
	return errExecuting
}

func (u *UserRepositoryPG) FindOneByCredentials(login string, password string) (*model.User, error) {
	sqlRequest := `SELECT id, login, password, token FROM public."user" 
                                  WHERE login = $1 AND password = $2 LIMIT 1`

	row := u.conn.Query(sqlRequest, login, password)
	return u.prepareModel(row)
}

func (u *UserRepositoryPG) prepareModel(row pgx.Row) (*model.User, error) {
	var userID, userLogin, userPass string
	var userToken sql.NullString

	user := model.User{}

	err := row.Scan(&userID, &userLogin, &userPass, &userToken)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("ошибка сканирования - " + err.Error())
	}

	user.SetID(userID)
	user.SetLogin(userLogin)
	user.SetPassword(userPass)
	if userToken.Valid {
		user.SetToken(userToken.String)
	}

	return &user, nil
}

func (u *UserRepositoryPG) FindOneByLogin(login string) (*model.User, error) {
	sqlRequest := `SELECT id, login, password, token FROM public."user" WHERE login = $1 LIMIT 1`

	row := u.conn.Query(sqlRequest, login)
	return u.prepareModel(row)
}

func (u *UserRepositoryPG) FindOneByToken(token string) (*model.User, error) {
	sqlRequest := `SELECT id, login, password, token FROM public."user" WHERE token = $1 LIMIT 1`

	row := u.conn.Query(sqlRequest, token)
	return u.prepareModel(row)
}

func (u *UserRepositoryPG) Update(user model.User) error {
	sqlRequest := `UPDATE public."user"
SET login    = $2,
    password = $3,
    token    = $4
WHERE id = $1;`

	userToken := sql.NullString{}
	if len(user.Token) > 0 {
		userToken.Valid = true
		userToken.String = user.Token
	} else {
		userToken.Valid = false
	}
	errExecuting := u.conn.Exec(sqlRequest, user.GetID(), user.GetLogin(), user.GetPassword(), userToken)
	return errExecuting
}

func NewUserRepositoryPG(connection *db.ConnectionPostgres) IUserRepository {
	return &UserRepositoryPG{conn: connection}
}
