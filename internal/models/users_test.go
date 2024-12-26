package models

import (
	"testing"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestUserModelExist(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	tests := []struct {
		name   string
		userID int
		exp    bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			exp:    true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			exp:    false,
		},
		{
			name:   "Negative ID",
			userID: -1,
			exp:    false,
		},
		{
			name:   "Non-existent ID",
			userID: 999,
			exp:    false,
		},
		{
			name:   "Large ID",
			userID: 1000000,
			exp:    false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exist, err := m.Exist(tc.userID)
			assert.Equal(t, exist, tc.exp)
			assert.NilError(t, err)
		})
	}
}
