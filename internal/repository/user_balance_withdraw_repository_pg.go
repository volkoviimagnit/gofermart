package repository

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserBalanceWithdrawRepositoryPG struct {
	conn *db.ConnectionPostgres
}

func NewUserBalanceWithdrawRepositoryPG(connection *db.ConnectionPostgres) IUserBalanceWithdrawRepository {
	return &UserBalanceWithdrawRepositoryPG{conn: connection}
}

func (u *UserBalanceWithdrawRepositoryPG) Insert(row model.UserBalanceWithdraw) error {
	sqlRequest := `INSERT INTO public."user_balance_withdraw" (user_id, order_number, sum, processed_at)
	VALUES ($1, $2, $3, $4);`

	errExecuting := u.conn.Exec(sqlRequest, row.GetUserId(), row.GetOrderNumber(), row.GetSum(), row.GetProcessedAt())
	return errExecuting
}

func (u *UserBalanceWithdrawRepositoryPG) FindByUserID(userID string) ([]model.UserBalanceWithdraw, error) {
	entities := make([]model.UserBalanceWithdraw, 0)

	request := `SELECT user_id, number, status, accrual, uploaded_at FROM public.user_order WHERE user_id = $1 ORDER BY uploaded_at ASC`
	rows, errConnection := u.conn.Query(request, userID)
	if errConnection != nil {
		return nil, errConnection
	}
	defer rows.Close()
	for rows.Next() {
		entity, errPreparing := u.prepareModel(rows)
		if errPreparing != nil {
			return nil, errPreparing
		}
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (u *UserBalanceWithdrawRepositoryPG) SumWithdrawByUserID(userID string) (float64, error) {
	sqlRequest := `
		SELECT COALESCE(sum(user_balance_withdraw.sum), 0) as sum
		FROM public.user_balance_withdraw
		WHERE user_id = $1`

	row, err := u.conn.QueryRow(sqlRequest, userID)
	if err != nil {
		return 0, err
	}
	var sum float64
	errScan := row.Scan(&sum)
	if errScan != nil {
		return 0, errScan
	}

	return sum, nil
}

func (u *UserBalanceWithdrawRepositoryPG) prepareModel(row pgx.Row) (*model.UserBalanceWithdraw, error) {
	userBalanceWithdraw := model.UserBalanceWithdraw{}

	err := row.Scan(&userBalanceWithdraw.UserId, &userBalanceWithdraw.OrderNumber, &userBalanceWithdraw.Sum, &userBalanceWithdraw.ProcessedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("ошибка сканирования - " + err.Error())
	}

	return &userBalanceWithdraw, nil
}
