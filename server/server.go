package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

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
			r.Get("/user/", s.GetLoginUser)

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
	Users []User `json:"users"`
}

type User struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Status     string `json:"status"`
	ChatNumber int    `json:"chat_number"`
}

func (s *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := ""
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[ERROR] not found User: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	users := make([]User, 0)
	for rows.Next() {
		var v User
		if err := rows.Scan(&v.Name, &v.ID, &v.Status, &v.ChatNumber); err != nil {
			log.Printf("[ERROR] scan user: %+v", err)
			writeHTTPError(w, http.StatusInternalServerError)
			return
		}

		users = append(users, v)
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
	Success string `json:"success"`
}

func (s *API) GetSignUp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// header を取得
	headerName := r.Header.Get("name")
	headerAddress := r.Header.Get("address")
	headerPassword := r.Header.Get("password")

	fmt.Println("header_name:", headerName)
	fmt.Println("header_address:", headerAddress)
	fmt.Println("header_password:", headerPassword)

	// ユーザー登録
	query := ""
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[ERROR] Insert: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	responseGetSignUp := &ResponseGetSignUp{
		Success: "true",
	}
	if err := json.NewEncoder(w).Encode(&responseGetSignUp); err != nil {
		log.Printf("[ERROR] response encoding failed: %+v", err)
		writeHTTPError(w, http.StatusInternalServerError)
		return
	}
}

func (s *API) GetLoginUser(w http.ResponseWriter, r *http.Request) {

}

func (s *API) GetLoginAdmin(w http.ResponseWriter, r *http.Request) {

}

func (s *API) PostUserProfile(w http.ResponseWriter, r *http.Request) {

}

func (s *API) PostAdminBan(w http.ResponseWriter, r *http.Request) {

}
