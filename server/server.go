package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// DB 設計

// User
// ID, Name, Status, ChatNumber, Token, Password

// Admin
// ID, Name, Token, Password

// User_Profile
// ID, Comment, Friend

// API API を表す構造体。
type API struct {
	// now 現在時刻を取得するための関数
	now func() time.Time
	// db データベースハンドラ
	db *sql.DB
}

func NewAPI(now func() time.Time, db *sql.DB) *API {
	return &API{now: now, db: db}
}

func (s *API) Handler() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/debug", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"message": "hello",
		}
		render.JSON(w, r, data)
	})

	r.Route("/api", func(r chi.Router) {

		// GET: /api/users
		r.Get("/users", s.GetUsers)

		// GET: /api/signup
		r.Get("/signup", s.GetSignUp)

		r.Route("/login", func(r chi.Router) {
			// GET: /api/user/:userID
			r.Get("/user", s.GetLoginUser)

			// GET: /api/admin/:adminID
			r.Get("/admin", s.GetLoginAdmin)
		})

		// POST: /api/user/profile
		r.Post("/user/profile", s.PostUserProfile)

		// POST: /api/admin/ban
		r.Post("/admin/ban", s.PostAdminBan)
	})

	addr := os.Getenv("Addr")
	if addr == "" {
		addr = ":1001"
	}

	log.Printf("listen: %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

type ResponseGetUsers struct {
	Users []ResponceGetUser `json:"users"`
}

type ResponceGetUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	ChatNumber int    `json:"chat_number"`
	CreatedAT  time.Time
}

type User struct {
	ID         string
	Name       string
	Address    string
	Status     string
	ChatNumber int
	Token      string
	Password   string
	CreatedAT  time.Time
	UpdatedAt  time.Time
}

func (s *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := "SELECT * FROM user"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[ERROR] not found User: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	users := make([]ResponceGetUser, 0)
	for rows.Next() {
		var v User
		if err := rows.Scan(&v.ID, &v.Name, &v.Address, &v.Status, &v.Password, &v.ChatNumber, &v.Token, &v.CreatedAT, &v.UpdatedAt); err != nil {
			log.Printf("[ERROR] scan user: %+v", err)
			writeHTTPError(w, http.StatusInternalServerError)
			return
		}

		user := ResponceGetUser{
			ID:         v.ID,
			Name:       v.Name,
			Status:     v.Status,
			ChatNumber: v.ChatNumber,
			CreatedAT:  v.CreatedAT,
		}

		users = append(users, user)
	}

	resp := &ResponseGetUsers{
		Users: users,
	}

	// レスポンスを返す。
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("[ERROR] response encoding failed: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}
}

type ResponseGetSignUp struct {
	Success bool `json:"success"`
}

func (s *API) GetSignUp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// header を取得
	headerName := r.Header.Get("name")
	headerAddress := r.Header.Get("address")
	headerPassword := r.Header.Get("password")

	// debug 用
	fmt.Println("header_name:", headerName)
	fmt.Println("header_address:", headerAddress)
	fmt.Println("header_password:", headerPassword)

	// address が以前登録されたものと一致しないか確認
	query := "select count(*) from user where address = ?"
	rows, err := s.db.QueryContext(ctx, query, headerAddress)

	if err != nil {
		log.Printf("[ERROR] not found User: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// count が 0 でない場合 Error
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Printf("[ERROR] scan user: %+v", err)
			writeHTTPError(w, http.StatusInternalServerError)
			return
		}
	}
	if count != 0 {
		log.Printf("[ERROR] address is already registered")
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// id を生成
	userID := "user_" + randomWithCharset(3)

	// ユーザー登録
	query2 := "INSERT INTO user (id, name, address, status, password, chat_number, token, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?) "
	_, err = s.db.ExecContext(ctx, query2, userID, headerName, headerAddress, "online", headerPassword, 0, "", s.now(), s.now())
	if err != nil {
		log.Printf("[ERROR] Insert: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	responseGetSignUp := &ResponseGetSignUp{
		Success: true,
	}
	if err := json.NewEncoder(w).Encode(&responseGetSignUp); err != nil {
		log.Printf("[ERROR] response encoding failed: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}
}

type ResponseGetLoginUser struct {
	Success string `json:"success"`
	ID      string `json:"id"`
	Token   string `json:"token"`
}

func randomWithCharset(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}

		b[i] = charset[n.Int64()]
	}

	return string(b)
}

func (s *API) GetLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// header を取得
	headerName := r.Header.Get("name")
	headerAddress := r.Header.Get("address")
	headerPassword := r.Header.Get("password")

	// debug 用
	fmt.Println("header_name:", headerName)
	fmt.Println("header_address:", headerAddress)
	fmt.Println("header_password:", headerPassword)

	// データベースから値を持ってくる
	query := "SELECT * FROM user WHERE name = ? AND address = ? AND password = ?"
	rows, err := s.db.QueryContext(ctx, query, headerName, headerAddress, headerPassword)
	if err != nil {
		log.Printf("[ERROR] can't login: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	var v User
	for rows.Next() {
		if err := rows.Scan(&v.ID, &v.Name, &v.Address, &v.Status, &v.Password, &v.ChatNumber, &v.Token, &v.CreatedAT, &v.UpdatedAt); err != nil {
			log.Printf("[ERROR] can't scan user: %+v", err)
			writeHTTPError(w, http.StatusInternalServerError)
			return
		}
	}

	// token 生成
	var token string
	token = v.ID + randomWithCharset(3)

	// token を登録
	query2 := "UPDATE user SET token = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, query2, token, v.ID)
	if err != nil {
		log.Printf("[ERROR] can't update token: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	responseGetSignUp := &ResponseGetLoginUser{
		Success: "true",
		ID:      v.ID,
		Token:   token,
	}

	if err := json.NewEncoder(w).Encode(&responseGetSignUp); err != nil {
		log.Printf("[ERROR] response encoding failed: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}
}

type Admin struct {
	ID       string `json:"id"`
	Token    string
	Password string
}

type ResponseGetLoginAdmin struct {
	Success string `json:"success"`
	ID      string `json:"id"`
	Token   string `json:"token"`
}

func (s *API) GetLoginAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// header を取得
	headerID := r.Header.Get("id")
	headerPassword := r.Header.Get("password")

	// debug 用
	fmt.Println("header_id:", headerID)
	fmt.Println("header_password:", headerPassword)

	// データベースから値を持ってくる
	query := ""
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[ERROR] can't login: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	var v Admin
	for rows.Next() {
		if err := rows.Scan(&v.ID, &v.Token, &v.Password); err != nil {
			log.Printf("[ERROR] can't scan admin: %+v", err)
			writeHTTPError(w, http.StatusInternalServerError)
			return
		}
	}

	// token 生成
	var token string
	token = v.ID + randomWithCharset(10)

	// token を登録
	query2 := ""
	_, err = s.db.ExecContext(ctx, query2)
	if err != nil {
		log.Printf("[ERROR] can't update token: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	responseGetLoginAdmin := &ResponseGetLoginAdmin{
		Success: "true",
		ID:      v.ID,
		Token:   token,
	}

	if err := json.NewEncoder(w).Encode(&responseGetLoginAdmin); err != nil {
		log.Printf("[ERROR] response encoding failed: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}
}

func (s *API) PostUserProfile(w http.ResponseWriter, r *http.Request) {

}

func (s *API) PostAdminBan(w http.ResponseWriter, r *http.Request) {

}
