package register_repository

import (
	"amikom-pedia-api/model/domain"
	"context"
	"database/sql"
)

type RegisterRepository interface {
	Create(ctx context.Context, tx *sql.Tx, register domain.Register) domain.Register
}