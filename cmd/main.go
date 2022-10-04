package main

import (
	"database/sql"
	"diarkis-server/server"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Open データベースを開く。
func Open(host string, port uint16, dbname, username, password string) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Addr = net.JoinHostPort(host, strconv.Itoa(int(port)))
	cfg.DBName = dbname
	cfg.User = username
	cfg.Passwd = password
	cfg.ParseTime = true

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, fmt.Errorf("new connector: %w", err)
	}

	return sql.OpenDB(connector), nil
}

func main() {

	fmt.Println("Service Start!")

	// TODO: 環境変数から取得する
	cfg := struct {
		DBHost     string
		DBPort     uint16
		DBName     string
		DBUsername string
		DBPassword string
	}{
		"localhost",
		3306,
		"server",
		"root",
		"tmp",
	}

	db, err := Open(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := server.NewAPI(time.Now, db)

	api.Handler()
}
