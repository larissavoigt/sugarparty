package db

import (
	"database/sql"

	"github.com/luizbranco/sugarparty/internal/order"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@/sugarparty?parseTime=true")
	if err != nil {
		panic(err)
	}
}

func Orders() (orders []order.Order, err error) {
	rows, err := db.Query(`
	SELECT id, name, email, status, created_at
	FROM orders
	ORDER BY id DESC`)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		o := order.Order{}
		err = rows.Scan(&o.ID, &o.Name, &o.Email, &o.Status, &o.CreatedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	err = rows.Err()
	return orders, err
}
