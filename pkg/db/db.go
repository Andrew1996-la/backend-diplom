package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(128) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	);

	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);	
`

func Init(dbFile string) error {
	var install bool

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("база данных не доступна: %w", err)
	}

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			DB.Close()
			return fmt.Errorf("не удалось выполнить инициализацию базы данных, %w", err)
		}
		fmt.Println("база данных успешно создана")
	} else {
		fmt.Println("подключение к существующей базе данных успешно выполнено")
	}

	return nil
}
