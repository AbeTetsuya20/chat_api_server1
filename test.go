package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Admin struct {
	ID        string
	Token     sql.NullString
	Password  string
	CreatedAT *time.Time
	UpdatedAt *time.Time
}

func main() {
	fmt.Println("Start!")
	// database に接続
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:tmp@tcp(127.0.0.1:3306)/server?parseTime=true")
	if err != nil {
		log.Fatalf("main sql.Open error err:%v", err)
	}
	defer db.Close()

	query := "SELECT * FROM admin"
	row, err := db.QueryContext(ctx, query)

	for row.Next() {
		a := &Admin{}
		err := row.Scan(&a.ID, &a.Token, &a.Password, &a.CreatedAT, &a.UpdatedAt)
		if err != nil {
			log.Fatalf("main sql.Scan error err:%v", err)
		}
		fmt.Printf("admin: %+v \n", a)
		fmt.Println("ID: ", a.ID)
	}

	fmt.Println("Success!")

	//db.QueryContext(ctx,)

}
