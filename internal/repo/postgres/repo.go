package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Franconl/ffaas/internal/core"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrKeyAlreadyUsed = errors.New("flag key already exists")
	ErrKeyRequired    = errors.New("key is required")
	ErrInvalidPercent = errors.New("invalid percentage")
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo { return &Repo{db: db} }

// --- helpers ---

func validate(f *core.FeatureFlag) error {
	if f.Key == "" {
		return ErrKeyRequired
	}
	if f.Percentage < 0 || f.Percentage > 100 {
		return ErrInvalidPercent
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

// --- CRUD ---

func (r *Repo) Create(ctx context.Context, f *core.FeatureFlag) error {
	if err := validate(f); err != nil {
		return err
	}
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	now := time.Now().UTC()

	const q = `
		INSERT INTO feature_flags
			(id, key, description, enabled, percentage, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := r.db.ExecContext(ctx, q, f.ID, f.Key, f.Description, f.Enabled, f.Percentage, now, now)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrKeyAlreadyUsed
		}
		return err
	}

	f.CreatedAt = now
	f.UpdatedAt = now
	return nil
}

func (r *Repo) Update(ctx context.Context, f *core.FeatureFlag) error {
	if err := validate(f); err != nil {
		return err
	}
	// Asegurar existencia y obtener key anterior para feedback/consistencia si querés
	const sel = `SELECT key FROM feature_flags WHERE id = $1`
	var oldKey string
	if err := r.db.QueryRowContext(ctx, sel, f.ID).Scan(&oldKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	const q = `
		UPDATE feature_flags
		   SET key = $1,
		       description = $2,
		       enabled = $3,
		       percentage = $4,
		       updated_at = NOW()
		 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, q, f.Key, f.Description, f.Enabled, f.Percentage, f.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrKeyAlreadyUsed
		}
		return err
	}
	return nil
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	const q = `DELETE FROM feature_flags WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*core.FeatureFlag, error) {
	const q = `
		SELECT id, key, description, enabled, percentage, created_at, updated_at
		  FROM feature_flags
		 WHERE id = $1`
	var ff core.FeatureFlag
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&ff.ID, &ff.Key, &ff.Description, &ff.Enabled, &ff.Percentage, &ff.CreatedAt, &ff.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// copia defensiva (retorno por puntero a copia local)
	return &ff, nil
}

func (r *Repo) GetByKey(ctx context.Context, key string) (*core.FeatureFlag, error) {
	const q = `
		SELECT id, key, description, enabled, percentage, created_at, updated_at
		  FROM feature_flags
		 WHERE key = $1`
	var ff core.FeatureFlag
	err := r.db.QueryRowContext(ctx, q, key).Scan(
		&ff.ID, &ff.Key, &ff.Description, &ff.Enabled, &ff.Percentage, &ff.CreatedAt, &ff.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &ff, nil
}

func (r *Repo) List(ctx context.Context) ([]core.FeatureFlag, error) {
	const q = `
		SELECT id, key, description, enabled, percentage, created_at, updated_at
		  FROM feature_flags
		 ORDER BY created_at ASC, key ASC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []core.FeatureFlag
	for rows.Next() {
		var ff core.FeatureFlag
		if err := rows.Scan(
			&ff.ID, &ff.Key, &ff.Description, &ff.Enabled, &ff.Percentage, &ff.CreatedAt, &ff.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, ff)
	}
	return out, rows.Err()
}
