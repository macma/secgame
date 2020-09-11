package db

import (
	"time"
)

func InsertClicks(color int) {
	conn.Exec(`INSERT INTO clicks (color, created_at) VALUES ($1, $2);`, color, time.Now())
}

func DropClicks() {
	conn.Exec(`DELETE FROM clicks`)
}
