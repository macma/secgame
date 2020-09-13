package db

import (
	"time"
)

type Click struct {
	Color      int
	Created_at time.Time
}

func InsertClicks(color int) {
	conn.Exec(`INSERT INTO clicks (color, created_at) VALUES ($1, $2);`, color, time.Now())
}

func SelectClicks() []Click {
	rows, _ := conn.Query("SELECT color, created_at FROM clicks ORDER BY created_at;")

	clicks := make([]Click, 0)
	for rows.Next() {
		var click Click
		rows.Scan(&click.Color, &click.Created_at)

		clicks = append(clicks, click)
	}
	return clicks
}

func DropClicks() {
	conn.Exec(`DELETE FROM clicks`)
}
