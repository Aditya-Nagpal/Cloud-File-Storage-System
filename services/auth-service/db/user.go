package db

import (
	"context"
)

func DoesEmailExist(ctx context.Context, email string) (bool, error) {
	var count int

	err := DB.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
