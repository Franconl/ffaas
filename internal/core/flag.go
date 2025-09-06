package core

import (
	"crypto/sha1"
	"encoding/binary"
	"time"
)

type FeatureFlag struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Percentage  int       `json:"percentage"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (f FeatureFlag) Eval(userID string) bool {
	if !f.Enabled {
		return false
	}

	if f.Percentage >= 100 {
		return true
	}

	if f.Percentage <= 0 {
		return false
	}

	data := []byte(f.Key + ":" + userID)
	h := sha1.Sum(data)
	val := binary.BigEndian.Uint32(h[:4])
	bucket := val % 100
	return int(bucket) < f.Percentage
}
