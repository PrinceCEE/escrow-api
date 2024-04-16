package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ITransactionRepository interface {
	Create(t *models.Transaction, tx pgx.Tx) error
	Update(t *models.Transaction, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Transaction, error)
	GetMany(args []any, where string, tx pgx.Tx) ([]*models.Transaction, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type TransactionRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewTransactionRepository(db *pgxpool.Pool, timeout time.Duration) *TransactionRepository {
	return &TransactionRepository{DB: db, Timeout: timeout}
}

func (repo *TransactionRepository) Create(t *models.Transaction, tx pgx.Tx) error {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	args := []any{
		t.Status,
		t.Type,
		t.SellerID,
		t.BuyerID,
		t.CreatedBy,
		t.DeliveryDuration,
		t.Currency,
		t.ChargeConfiguration,
		t.ProductDetails,
		t.TotalAmount,
		t.TotalCost,
		t.Charges,
		t.ReceivableAmount,
		t.CreatedAt,
		t.UpdatedAt,
	}

	query := `INSERT INTO transactions (status, type, seller_id, buyer_id, created_by, delivery_duration, currency, charge_configuration, product_details, total_amount, total_cost, charges, receivable_amount, created_at, updated_at
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, version`

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &t.Version)
		if err != nil {
			return err
		}

		t.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &t.Version)
	if err != nil {
		return err
	}

	t.ID = id.String()
	return nil
}

func (repo *TransactionRepository) Update(t *models.Transaction, tx pgx.Tx) error {
	t.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(t, "transactions")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&t.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&t.Version)
}

func (repo *TransactionRepository) getByKey(key string, value any, tx pgx.Tx) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	t := new(models.Transaction)
	var id, buyerId, sellerId uuid.UUID
	var sellerImgUrl, buyerImgUrl *string
	var buyer *models.User
	var seller *models.Business

	query := fmt.Sprintf(`
		SELECT
			t.id,
			t.status,
			t.type,
			t.created_by,
			t.buyer_id,
			t.seller_id,
			t.delivery_duration,
			t.currency,
			t.charge_configuration,
			t.product_details,
			t.total_amount,
			t.total_cost,
			t.charges,
			t.receivable_amount,
			t.created_at,
			t.updated_at,
			t.deleted_at,
			u.email,
			u.phone_number,
			u.first_name,
			u.last_name,
			u.is_phone_number_verified,
			u.is_email_verified,
			u.reg_stage,
			u.account_type,
			u.image_url,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			b.name,
			b.email,
			b.image_url,
			b.created_at,
			b.updated_at,
			b.deleted_at

		FROM transactions t
		INNER JOIN businesses b ON b.id = t.seller_id
		INNER JOIN users u ON u.id= t.buyer_id
		WHERE %s = $1`,
		key,
	)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&t.Status,
		&t.Type,
		&t.CreatedBy,
		&buyerId,
		&sellerId,
		&t.DeliveryDuration,
		&t.Currency,
		&t.ChargeConfiguration,
		&t.ProductDetails,
		&t.TotalAmount,
		&t.TotalCost,
		&t.Charges,
		&t.ReceivableAmount,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
		&buyer.Email,
		&buyer.PhoneNumber,
		&buyer.FirstName,
		&buyer.LastName,
		&buyer.IsPhoneNumberVerified,
		&buyer.IsEmailVerified,
		&buyer.RegStage,
		&buyer.AccountType,
		&buyerImgUrl,
		&buyer.CreatedAt,
		&buyer.UpdatedAt,
		&buyer.DeletedAt,
		&seller.Name,
		&seller.Email,
		&sellerImgUrl,
		&seller.CreatedAt,
		&seller.UpdatedAt,
		&seller.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	t.ID = id.String()
	t.SellerID = sellerId.String()
	t.BuyerID = buyerId.String()
	seller.ID = sellerId.String()
	buyer.ID = buyerId.String()

	if sellerImgUrl != nil {
		seller.ImageUrl = *sellerImgUrl
	}
	if buyerImgUrl != nil {
		buyer.ImageUrl = *buyerImgUrl
	}

	t.Buyer = buyer
	t.Seller = seller

	return t, nil
}

func (repo *TransactionRepository) GetById(id string, tx pgx.Tx) (*models.Transaction, error) {
	return repo.getByKey("t.id", id, tx)
}

func (repo *TransactionRepository) GetMany(args []any, where string, tx pgx.Tx) ([]*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	argLen := len(args)
	query := fmt.Sprintf(`
		SELECT
			t.id,
			t.status,
			t.type,
			t.created_by,
			t.buyer_id,
			t.seller_id,
			t.delivery_duration,
			t.currency,
			t.charge_configuration,
			t.product_details,
			t.total_amount,
			t.total_cost,
			t.charges,
			t.receivable_amount,
			t.created_at,
			t.updated_at,
			t.deleted_at,
			u.email,
			u.phone_number,
			u.first_name,
			u.last_name,
			u.is_phone_number_verified,
			u.is_email_verified,
			u.reg_stage,
			u.account_type,
			u.image_url,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			b.name,
			b.email,
			b.image_url,
			b.created_at,
			b.updated_at,
			b.deleted_at
		
		FROM transactions t
		INNER JOIN businesses b ON b.id = t.seller_id
		INNER JOIN users u ON u.id= t.buyer_id
		%s
		OFFSET $%d
		LIMIT $%d
	`, where, argLen-1, argLen)

	var rows pgx.Rows
	if tx != nil {
		_rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	} else {
		_rows, err := repo.DB.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	}

	transactions := []*models.Transaction{}

	for rows.Next() {
		var t *models.Transaction
		var id, buyerId, sellerId uuid.UUID
		var sellerImgUrl, buyerImgUrl *string
		var buyer *models.User
		var seller *models.Business

		err := rows.Scan(
			&id,
			&t.Status,
			&t.Type,
			&t.CreatedBy,
			&buyerId,
			&sellerId,
			&t.DeliveryDuration,
			&t.Currency,
			&t.ChargeConfiguration,
			&t.ProductDetails,
			&t.TotalAmount,
			&t.TotalCost,
			&t.Charges,
			&t.ReceivableAmount,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
			&buyer.Email,
			&buyer.PhoneNumber,
			&buyer.FirstName,
			&buyer.LastName,
			&buyer.IsPhoneNumberVerified,
			&buyer.IsEmailVerified,
			&buyer.RegStage,
			&buyer.AccountType,
			&buyerImgUrl,
			&buyer.CreatedAt,
			&buyer.UpdatedAt,
			&buyer.DeletedAt,
			&seller.Name,
			&seller.Email,
			&sellerImgUrl,
			&seller.CreatedAt,
			&seller.UpdatedAt,
			&seller.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		t.ID = id.String()
		t.SellerID = sellerId.String()
		t.BuyerID = buyerId.String()
		seller.ID = sellerId.String()
		buyer.ID = buyerId.String()

		if sellerImgUrl != nil {
			seller.ImageUrl = *sellerImgUrl
		}
		if buyerImgUrl != nil {
			buyer.ImageUrl = *buyerImgUrl
		}

		t.Buyer = buyer
		t.Seller = seller

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *TransactionRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM transactions WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query, id)
	} else {
		_, err = repo.DB.Exec(ctx, query, id)
	}

	return
}

func (repo *TransactionRepository) SoftDelete(id string, tx pgx.Tx) error {
	u, err := repo.GetById(id, tx)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	u.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}
	u.UpdatedAt = now
	return repo.Update(u, tx)
}
