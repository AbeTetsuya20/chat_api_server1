package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse エラーレスポンスを表す。
type errorResponse struct {
	// Message エラーメッセージ
	Message string `json:"message"`
}

// writeErrorResponse w にエラーレスポンスを書き込む。
func writeErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)

	resp := &errorResponse{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("[ERROR] error response encoding failed: %+v\n", err)
	}
}

// writeHTTPError w にエラーレスポンスを書き込む。
// エラーメッセージは HTTP ステータスコードに対応した文字列になる。
func writeHTTPError(w http.ResponseWriter, code int) {
	writeErrorResponse(w, code, http.StatusText(code))
}
