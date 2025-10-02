package cached

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Franconl/ffaas/internal/core"
	"github.com/redis/go-redis/v9"
)

// Errores (mismos contratos que otras capas)
var (
	ErrNotFound       = errors.New("not found")
	ErrKeyAlreadyUsed = errors.New("flag key already exists")
	ErrKeyRequired    = errors.New("key is required")
	ErrInvalidPercent = errors.New("invalid percentage")
)

// Contrato que debe cumplir el backend (memory, postgres, etc.)
type BaseRepo interface {
	Create(ctx context.Context, f *core.FeatureFlag) error
	Update(ctx context.Context, f *core.FeatureFlag) error
	DeleteByID(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*core.FeatureFlag, error)
	GetByKey(ctx context.Context, key string) (*core.FeatureFlag, error)
	List(ctx context.Context) ([]core.FeatureFlag, error)
}

type Repo struct {
	base BaseRepo
	rdb  *redis.Client
	ttl  time.Duration
}

func New(base BaseRepo, rdb *redis.Client, ttl time.Duration) *Repo {
	return &Repo{base: base, rdb: rdb, ttl: ttl}
}

// --- keys ---

func keyByID(id string) string { return fmt.Sprintf("ff:id:%s", id) }
func keyByKey(k string) string { return fmt.Sprintf("ff:key:%s", k) }
func setJSON(ctx context.Context, rdb *redis.Client, k string, v any, ttl time.Duration) {
	b, _ := json.Marshal(v)
	_ = rdb.Set(ctx, k, b, ttl).Err()
}

func getJSON[T any](ctx context.Context, rdb *redis.Client, k string, dst *T) (bool, error) {
	s, err := rdb.Get(ctx, k).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(s, dst); err != nil {
		return false, nil
	}
	return true, nil
}

// --- CRUD con caché ---

func (r *Repo) Create(ctx context.Context, f *core.FeatureFlag) error {
	if err := r.base.Create(ctx, f); err != nil {
		return err
	}
	// popular caché
	setJSON(ctx, r.rdb, keyByID(f.ID), f, r.ttl)
	setJSON(ctx, r.rdb, keyByKey(f.Key), f, r.ttl)
	return nil
}

func (r *Repo) Update(ctx context.Context, f *core.FeatureFlag) error {
	// Para invalidar correctamente, obtener el estado actual (clave vieja)
	old, _ := r.base.GetByID(ctx, f.ID)

	if err := r.base.Update(ctx, f); err != nil {
		return err
	}

	// invalidar y reescribir caché
	if old != nil {
		_ = r.rdb.Del(ctx, keyByID(old.ID)).Err()
		_ = r.rdb.Del(ctx, keyByKey(old.Key)).Err()
	}
	setJSON(ctx, r.rdb, keyByID(f.ID), f, r.ttl)
	setJSON(ctx, r.rdb, keyByKey(f.Key), f, r.ttl)
	return nil
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	// obtener para invalidar por key
	old, _ := r.base.GetByID(ctx, id)
	if err := r.base.DeleteByID(ctx, id); err != nil {
		return err
	}
	_ = r.rdb.Del(ctx, keyByID(id)).Err()
	if old != nil {
		_ = r.rdb.Del(ctx, keyByKey(old.Key)).Err()
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*core.FeatureFlag, error) {
	var ff core.FeatureFlag
	if ok, err := getJSON(ctx, r.rdb, keyByID(id), &ff); err == nil && ok {
		return &ff, nil
	}
	v, err := r.base.GetByID(ctx, id)
	if err != nil || v == nil {
		return v, err
	}
	setJSON(ctx, r.rdb, keyByID(v.ID), v, r.ttl)
	setJSON(ctx, r.rdb, keyByKey(v.Key), v, r.ttl)
	return v, nil
}

func (r *Repo) GetByKey(ctx context.Context, k string) (*core.FeatureFlag, error) {
	var ff core.FeatureFlag
	if ok, err := getJSON(ctx, r.rdb, keyByKey(k), &ff); err == nil && ok {
		return &ff, nil
	}
	v, err := r.base.GetByKey(ctx, k)
	if err != nil || v == nil {
		return v, err
	}
	setJSON(ctx, r.rdb, keyByID(v.ID), v, r.ttl)
	setJSON(ctx, r.rdb, keyByKey(v.Key), v, r.ttl)
	return v, nil
}

func (r *Repo) List(ctx context.Context) ([]core.FeatureFlag, error) {
	// Podrías cachear páginas, pero para mantenerlo simple, vamos directo a base.
	return r.base.List(ctx)
}
