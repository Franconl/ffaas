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

func (r *Repo) Create(f *core.FeatureFlag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if f.Key == "" {
		return errors.New("key is required")
	}

	if _, exist := r.byKey[f.Key]; exist {
		return errors.New("flag key already exist")
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
		return errors.New("id not found")
	}

	if f.Key == "" {
		return errors.New("key is required")
	}

	if cur.Key != f.Key {
		if otherID, exist := r.byKey[f.Key]; exist && otherID != f.ID {
			return errors.New("flag key already exists")
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
		return errors.New("flag id not found")
	}

	return nil
}

func (r *Repo) GetByID(id string) (*core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, exist := r.byID[id]

	if !exist {
		return nil, errors.New("id not found")
	}

	cur := flag
	return &cur, nil
}

func (r *Repo) GetByKey(key string) (*core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, exist := r.byID[key]

	if !exist {
		return nil, errors.New("key not found")
	}

	cur := flag
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
