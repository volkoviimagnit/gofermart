package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/volkoviimagnit/gofermart/internal/db"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserOrderRepositoryPG struct {
	conn *db.ConnectionPostgres
}

func NewUserOrderRepositoryPG(connection *db.ConnectionPostgres) IUserOrderRepository {
	return &UserOrderRepositoryPG{conn: connection}
}

func (u *UserOrderRepositoryPG) Insert(row model.UserOrder) error {
	return u.Upsert(row)
}

func (u *UserOrderRepositoryPG) Update(row model.UserOrder) error {
	return u.Upsert(row)
}

func (u *UserOrderRepositoryPG) Upsert(row model.UserOrder) error {
	sqlRequest := `INSERT INTO public."user_order" (user_id, number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (number)
			DO UPDATE
			SET user_id = excluded.user_id, status = excluded.status, accrual = excluded.accrual, uploaded_at = excluded.uploaded_at;
`

	accrual := sql.NullFloat64{
		Valid: false,
	}
	if row.Accrual != nil {
		accrual = sql.NullFloat64{
			Float64: *(row.Accrual),
			Valid:   true,
		}
	}
	errExecuting := u.conn.Exec(sqlRequest, row.UserID, row.Number, row.Status.String(), accrual, row.UploadedAt)
	return errExecuting
}

func (u *UserOrderRepositoryPG) FindByUserID(userID string) ([]model.UserOrder, error) {
	entities := make([]model.UserOrder, 0)

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

func (u *UserOrderRepositoryPG) FindOneByNumber(number string) (*model.UserOrder, error) {
	sqlRequest := `SELECT user_id, number, status, accrual, uploaded_at FROM public."user_order" WHERE number = $1 LIMIT 1`

	row, err := u.conn.QueryRow(sqlRequest, number)
	if err != nil {
		return nil, err
	}
	return u.prepareModel(row)
}

func (u *UserOrderRepositoryPG) SumAccrualByUserID(userID string) (float64, error) {
	sqlRequest := `
		SELECT COALESCE(sum(accrual), 0) as sum
		FROM public.user_order 
		WHERE user_id = $1 AND status = $2`

	row, err := u.conn.QueryRow(sqlRequest, userID, model.UserOrderStatusProcessed)
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

func (u *UserOrderRepositoryPG) prepareModel(row pgx.Row) (*model.UserOrder, error) {
	var userID, orderNumber, status string
	var uploadedAt time.Time
	var accrual sql.NullFloat64

	err := row.Scan(&userID, &orderNumber, &status, &accrual, &uploadedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("ошибка сканирования - " + err.Error())
	}

	userOrder := model.UserOrder{
		UserID:     userID,
		Number:     orderNumber,
		Status:     model.UserOrderStatus(status),
		UploadedAt: uploadedAt,
	}

	if accrual.Valid {
		userOrder.Accrual = &accrual.Float64
	}

	return &userOrder, nil
}
