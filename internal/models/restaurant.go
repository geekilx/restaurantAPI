package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/geekilx/restaurantAPI/internal/validator"
)

type Restaurant struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	FullAddress string    `json:"full_address"`
	Cuisine     string    `json:"cuisine"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RestaurantModel struct {
	DB *sql.DB
}

func (m *RestaurantModel) Insert(restaurant *Restaurant) (int64, error) {
	stmt := `INSERT INTO restaurant (name, country, full_address, cuisine, status) VALUES($1, $2, $3, $4, $5) 
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{restaurant.Name, restaurant.Country, restaurant.FullAddress, restaurant.Cuisine, restaurant.Status}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&restaurant.ID, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "restaurant_name_key"):
			return 0, ErrDuplicateRestaurantName
		default:
			return 0, err
		}
	}
	return restaurant.ID, nil

}

func (m *RestaurantModel) GetAll(name string, f Filters) ([]*Restaurant, error) {
	stmt := fmt.Sprintf(`SELECT * FROM restaurant WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC LIMIT %d OFFSET %d`, f.sortColumn(), f.sortDirection(), f.Limit(), f.Offset())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var restaurants []*Restaurant

	for rows.Next() {
		var restaurant Restaurant

		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Country, &restaurant.FullAddress, &restaurant.Cuisine, &restaurant.Status, &restaurant.CreatedAt, &restaurant.UpdatedAt)
		if err != nil {
			return nil, err
		}

		restaurants = append(restaurants, &restaurant)

	}

	return restaurants, nil

}

func (m *RestaurantModel) Update(id int64, restaurant Restaurant) error {

	stmt := `UPDATE restaurant SET name = $1, country = $2, full_address = $3, cuisine = $4, status = $5, updated_at = NOW() WHERE id = $6`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{restaurant.Name, restaurant.Country, restaurant.FullAddress, restaurant.Cuisine, restaurant.Status, id}

	rows, err := m.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "restaurant_name_key"):
			return ErrDuplicateRestaurantName
		default:
			return err
		}
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRestaurantNotFound
	}

	return nil

}

func (m *RestaurantModel) Get(id int64) (*Restaurant, error) {
	stmt := `SELECT * FROM restaurant WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restaurant Restaurant
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Country, &restaurant.FullAddress, &restaurant.Cuisine, &restaurant.Status, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRestaurantNotFound
		default:
			return nil, err
		}
	}

	return &restaurant, nil
}

func (m *RestaurantModel) Delete(id int64) error {
	stmt := `DELETE FROM restaurant where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRestaurantNotFound
	}

	return nil
}

func (m *RestaurantModel) CheckIfRestaurantExists(id int64) bool {
	stmt := `SELECT EXISTS(SELECT FROM restaurant WHERE id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ok bool
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&ok)
	if err != nil {
		ok = false
	}
	return ok
}

func ValidateRestaurant(v *validator.Validator, res Restaurant) {
	v.Check(v.Empty(res.Name), "name", "restaurant name must be provided")
	v.Check(v.Empty(res.Country), "country", "country must be provided")
	v.Check(v.Empty(res.FullAddress), "full Address", "full address  must be provided")
	v.Check(v.Empty(res.Cuisine), "cuisine", "cuisine must be provided")
	v.Check(v.Empty(res.Status), "status", "status must be provided")

	v.Check(len(res.Name) < 3 || len(res.Name) > 50, "name", "restaurant name must be greater than 3 and less than 50 characters")
	v.Check(len(res.Country) < 3 || len(res.Country) > 50, "country", "country must be greater than 3 and less than 50 characters")
	v.Check(len(res.FullAddress) < 10 || len(res.FullAddress) > 200, "full Address", "full address must be greater than 10 and less than 200 characters")
	v.Check(len(res.Cuisine) < 3 || len(res.Cuisine) > 50, "cuisine", "cuisine must be greater than 3 and less than 50 characters")

	v.Check(!validator.PermittedValue(res.Status, "open", "closed"), "status", "you have to provide valid status (open,closed)")
}
