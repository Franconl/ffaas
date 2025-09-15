package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/Franconl/ffaas/internal/core"
	"github.com/google/uuid"
)

type Repo struct {
	mu    sync.RWMutex
	byID  map[string]core.FeatureFlag
	byKey map[string]string
}

func New() *Repo {
	return &Repo{
		byID:  make(map[string]core.FeatureFlag),
		byKey: make(map[string]string),
	}
}

var (
	ErrNotFound       = errors.New("not found")
	ErrKeyAlreadyUsed = errors.New("flag key already exists")
	ErrKeyRequired    = errors.New("key is required")
	ErrInvalidPercent = errors.New("invalid percentage")
)

func (r *Repo) Create(f *core.FeatureFlag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if f.Key == "" {
		return ErrNotFound
	}

	if _, exist := r.byKey[f.Key]; exist {
		return ErrKeyAlreadyUsed
	}

	if f.ID == "" {
		f.ID = uuid.NewString()
	}

	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	r.byID[f.ID] = *f
	r.byKey[f.Key] = f.ID

	return nil
}

func (r *Repo) Update(f *core.FeatureFlag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cur, exist := r.byID[f.ID]

	if !exist {
		return ErrNotFound
	}

	if f.Key == "" {
		return ErrKeyRequired
	}

	if cur.Key != f.Key {
		if otherID, exist := r.byKey[f.Key]; exist && otherID != f.ID {
			return ErrKeyAlreadyUsed
		} else {

			delete(r.byKey, cur.Key)
			r.byKey[f.Key] = cur.ID

		}
	}

	cur.Key = f.Key
	cur.Description = f.Description
	cur.Enabled = f.Enabled
	cur.Percentage = f.Percentage
	cur.UpdatedAt = time.Now()

	r.byID[f.ID] = cur

	return nil
}

func (r *Repo) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if flag, exist := r.byID[id]; exist {

		delete(r.byID, id)
		delete(r.byKey, flag.Key)

	} else {
		return ErrNotFound
	}

	return nil
}

func (r *Repo) GetByID(id string) (*core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, exist := r.byID[id]

	if !exist {
		return nil, ErrNotFound
	}

	cur := flag
	return &cur, nil
}

func (r *Repo) GetByKey(key string) (*core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exist := r.byKey[key]

	if !exist || id == "" {
		return nil, ErrNotFound
	}

	cur := r.byID[id]
	return &cur, nil
}

func (r *Repo) List() ([]core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]core.FeatureFlag, 0, len(r.byID))

	for _, f := range r.byID {
		list = append(list, f)
	}

	return list, nil
}
