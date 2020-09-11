package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var conn *sql.DB

// InitDB initializes database connection with provided connection
func InitDBWithConn(newConn *sql.DB) {
	conn = newConn
}

// InitDB initializes database connection
func InitDB() (db *sql.DB) {
	connStr := `postgres://postgres:postgres@192.168.1.15/clp?sslmode=require`

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		// logger.Panicw(ctx, fmt.Sprintf("there's error during make connection to db. %v", err))
	}

	if err = db.Ping(); err != nil {
		// logger.Panicw(ctx, fmt.Sprintf("there's error during ping connection to db. %v", err))
	}

	conn = db

	fmt.Println("Initialized Database successfully")
	// logger.Infow(ctx, fmt.Sprint("Initialized Database successfully"))
	// MigrateDB(ctx)
	return db
}

func getSliceFromRow(ctx context.Context, rows *sql.Rows) (cols []string, vals [][]string) {

	cols, err := rows.Columns()
	if err != nil {
		// logger.Errorw(ctx, fmt.Sprintf("Failed to get columns, %v", err))
		return
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	var resultArray [][]string

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			// logger.Errorw(ctx, fmt.Sprintf("Failed to scan row, %v", err))
			return
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		resultArray = append(resultArray, result)
	}

	return cols, resultArray
}
