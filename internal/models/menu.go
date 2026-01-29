package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Menu struct {
	ID             int64     `json:"id"`
	CategoryID     int64     `json:"category_id,omitempty"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	RestaurantName string    `json:"restaurant_name"`
	PriceCent      float32   `json:"price_cent"`
	IsAvaiable     bool      `json:"is_available"`
	CreatedAt      time.Time `json:"-"`
}

type MenuWithCategoryName struct {
	Menu
	CategoryName string `json:"category_name"`
}

type MenuModel struct {
	DB *sql.DB
}

func (m *MenuModel) Insert(menu *Menu) error {
	stmt := `INSERT INTO menu (category_id, name, description, price_cent) VALUES($1, $2, $3, $4) RETURNING id, is_available, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{menu.CategoryID, menu.Name, menu.Description, menu.PriceCent}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&menu.ID, &menu.IsAvaiable, &menu.CreatedAt)
	return err

}

func (m *MenuModel) GetAll(name string, f Filters) ([]*Menu, error) {

	sortColumn := f.sortColumn()
	safeSortColumn := "m.id" // Default fallback

	switch sortColumn {
	case "name":
		safeSortColumn = "m.name"
	case "price_cent":
		safeSortColumn = "m.price_cent"
	case "restaurant_name": // Use the alias 'r' for the joined table!
		safeSortColumn = "r.name"
	case "category_id":
		safeSortColumn = "m.category_id"
		// Add other cases here
	}

	stmt := fmt.Sprintf(`SELECT m.id, m.category_id, r.name, m.name, m.description, m.price_cent, m.is_available, m.created_at FROM menu m
	INNER JOIN categories c on c.id = m.category_id
	INNER JOIN restaurant r on r.id = c.restaurant_id
	WHERE (to_tsvector('simple', m.name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, m.id ASC LIMIT %d OFFSET %d`, safeSortColumn, f.sortDirection(), f.Limit(), f.Offset())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, name)
	if err != nil {
		return nil, err
	}

	var menus []*Menu
	for rows.Next() {
		var menu Menu

		err := rows.Scan(&menu.ID, &menu.CategoryID, &menu.RestaurantName, &menu.Name, &menu.Description, &menu.PriceCent, &menu.IsAvaiable, &menu.CreatedAt)
		if err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, nil

}

func (m *MenuModel) GetRestaurantMenus(id int64) ([]*MenuWithCategoryName, error) {
	stmt := `SELECT m.id, m.name, c.name, r.name, m.description, m.price_cent, m.is_available from menu m
	INNER JOIN categories c on c.id = m.category_id
	INNER JOIN restaurant r on r.id = c.restaurant_id
	WHERE c.restaurant_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var menus []*MenuWithCategoryName
	for rows.Next() {
		var menu MenuWithCategoryName
		err := rows.Scan(&menu.ID, &menu.Name, &menu.CategoryName, &menu.RestaurantName, &menu.Description, &menu.PriceCent, &menu.IsAvaiable)
		if err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, nil

}
func (m *MenuModel) GetAllMenuForCategory(id int64) ([]*Menu, error) {
	stmt := `SELECT id, category_id, name, description, price_cent, is_available from menu WHERE category_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var menus []*Menu
	for rows.Next() {
		var menu Menu
		err := rows.Scan(&menu.ID, &menu.CategoryID, &menu.Name, &menu.Description, &menu.PriceCent, &menu.IsAvaiable)
		if err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, nil

}
