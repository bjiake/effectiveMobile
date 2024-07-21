package db

import (
	"database/sql"
	"effectiveMobile/pkg/config"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// ConnectToBD Подключение к PostgresSql по app.env
func ConnectToBD(cfg config.Config) (*sql.DB, error) {
	// Формирование строки подключения из конфига
	addr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PsqlUser, cfg.PsqlPass, cfg.PsqlHost, cfg.PsqlPort, cfg.PsqlDBName)
	fmt.Printf("\nConnecting to %s\n", addr)
	psqlInfo := addr

	// Подключение к БД
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
