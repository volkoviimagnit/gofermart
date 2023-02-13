package repository

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserBalanceRepositoryPG struct {
	conn *db.ConnectionPostgres
}

func (u *UserBalanceRepositoryPG) Insert(row model.UserBalance) error {
	return u.Upset(row)
}

func (u *UserBalanceRepositoryPG) FinOneByUserID(userID string) (*model.UserBalance, error) {
	sqlRequest := `SELECT user_id, current, withdrawn FROM public."user_balance" WHERE user_id = $1 LIMIT 1`

	row, _ := u.conn.QueryRow(sqlRequest, userID)
	userBalance, err := u.prepareModel(row)
	if err != nil {
		return nil, err
	}
	if userBalance == nil {
		return nil, nil
	}
	return userBalance, nil
}

func (u *UserBalanceRepositoryPG) prepareModel(row pgx.Row) (*model.UserBalance, error) {
	userBalance := model.UserBalance{}

	err := row.Scan(&userBalance.UserID, &userBalance.Current, &userBalance.Withdrawn)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("ошибка сканирования - " + err.Error())
	}

	return &userBalance, nil
}

func (u *UserBalanceRepositoryPG) Update(row model.UserBalance) error {
	return u.Upset(row)
}

func (u *UserBalanceRepositoryPG) Upset(row model.UserBalance) error {
	errConnecting := u.conn.TryConnect()
	if errConnecting != nil {
		return errConnecting
	}

	sqlRequest := `INSERT INTO public.user_balance (user_id, current, withdrawn)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id)
			DO UPDATE
			SET current = excluded.current, withdrawn = excluded.withdrawn;`

	errExecuting := u.conn.Exec(sqlRequest, row.UserID, row.Current, row.Withdrawn)
	return errExecuting
}

func NewUserBalanceRepositoryPG(conn *db.ConnectionPostgres) IUserBalanceRepository {
	return &UserBalanceRepositoryPG{conn: conn}
}
