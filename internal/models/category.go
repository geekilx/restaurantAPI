package models

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID             int64     `json:"id"`
	RestaurantID   int64     `json:"restaurant_id"`
	Name           string    `json:"name"`
	RestaurantName string    `json:"restaurant_name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Insert(category *Category) error {
	stmt := `INSERT INTO categories (restaurant_id, name) VALUES($1, $2) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, category.RestaurantID, category.Name).Scan(&category.ID, &category.CreatedAt)
	if err != nil {
		return err
	}

	return nil

}

func (m *CategoryModel) CategoryExists(name string, restaurantID int64) bool {
	stmt := `SELECT EXISTS(SELECT FROM categories where name = $1 AND restaurant_id = $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ok bool
	err := m.DB.QueryRowContext(ctx, stmt, name, restaurantID).Scan(&ok)
	if err != nil || ok {
		return true
	}

	return false

}

func (m *CategoryModel) GetAll() ([]*Category, error) {
	stmt := `SELECT c.id, r.name, c.restaurant_id, c.name, c.created_at FROM categories c inner join restaurant r on r.id = c.restaurant_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var categories []*Category
	for rows.Next() {
		var category Category

		err := rows.Scan(&category.ID, &category.RestaurantName, &category.RestaurantID, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil

}

func (m *CategoryModel) GetAllForRestaurant(id int64) ([]*Category, error) {
	stmt := `SELECT * from categories where restaurant_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	var categories []*Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.RestaurantID, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (m *CategoryModel) CheckIfExists(id int64) bool {
	stmt := `SELECT EXISTS(SELECT FROM categories WHERE id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ok bool
	_ = m.DB.QueryRowContext(ctx, stmt, id).Scan(&ok)
	if !ok {
		return false
	}
	return true

}
