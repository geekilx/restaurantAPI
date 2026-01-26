package models

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

type PermissionModel struct {
	DB *sql.DB
}

func (m *PermissionModel) GetForAllUser(userID int64) (Permissions, error) {
	stmt := `SELECT p.code FROM Permissions AS p
	INNER JOIN users_permissions as up on up.Permission_id = p.id
	INNER JOIN users AS u on up.user_id = u.id
	WHERE up.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil

}

func (m *PermissionModel) AddForUser(userID int64, codes ...string) error {
	stmt := `INSERT INTO users_permissions
	SELECT $1, permissions.id from Permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, userID, pq.Array(codes))
	return err

}
