package core

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEvalBasics(t *testing.T) {
	flag := FeatureFlag{
		ID:          "123",
		Key:         "new_checkin",
		Description: "nuevo flujo checkin",
		Enabled:     true,
		Percentage:  100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	prueba := flag.Eval("1234")

	if prueba != true {
		t.Errorf("expected true, got %v", prueba)
	}

	flag.Percentage = 0

	prueba2 := flag.Eval("1234")

	if prueba2 != false {
		t.Errorf("expected false, got %v", prueba2)
	}
}

func TestEval_Monontonicity_Approx(t *testing.T) {
	users := GenerateUserID(1000)

	if len(users) > 1000 {
		t.Fatalf("expected 1000, got %d", len(users))
	}

	flag := FeatureFlag{
		ID:          "123",
		Key:         "new_checkin",
		Description: "nuevo flujo checkin",
		Enabled:     true,
		Percentage:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	in10 := 0
	out10 := 0

	for _, u := range users {
		eval := flag.Eval(u)
		if eval {
			in10 += 1
		} else {
			out10 += 1
		}
	}

	flag.Percentage = 60

	in60 := 0
	out60 := 0

	for _, u := range users {
		eval := flag.Eval(u)
		if eval {
			in60 += 1
		} else {
			out60 += 1
		}
	}

	if in60 < in10 {
		t.Errorf("expected more users in percentage 60")
	}
}

func GenerateUserID(n int) []string {
	users := make([]string, 0, n-1)

	for i := 0; i < n; i++ {
		users = append(users, uuid.NewString())
	}

	return users
}
