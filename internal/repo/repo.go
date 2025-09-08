package repo

import "github.com/Franconl/ffaas/internal/core"

type Flags interface {
	Create(f *core.FeatureFlag) error
	Update(f *core.FeatureFlag) error
	DeleteByID(id string) error
	GetByID(id string) error
	GetByKey(id string) error
	List() ([]core.FeatureFlag, error)
}
