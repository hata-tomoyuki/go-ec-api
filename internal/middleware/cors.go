package middleware

import (
	"net/http"
	"os"

	"github.com/go-chi/cors"
)

// CSRF 対策が追加で不要な理由:
//
//  1. 認証は Authorization: Bearer <JWT> ヘッダーで行われる。
//     ブラウザはクロスオリジンリクエストに自動で Authorization ヘッダーを付与しないため、
//     攻撃者のサイトから被害者のトークンを使ったリクエストは送信できない。
//  2. Cookie (リフレッシュトークン) は httpOnly + SameSite=Lax で設定されており、
//     クロスサイトの POST リクエストには Cookie が送信されない。
//  3. CORS で許可オリジンを明示的に制限しているため、
//     悪意のあるオリジンからのプリフライト付きリクエストはブロックされる。

// CORS は環境変数 CORS_ALLOWED_ORIGIN で許可オリジンを設定可能な CORS ミドルウェアを返す。
// 未設定の場合は http://localhost:3000（開発用）をデフォルトとする。
func CORS() func(http.Handler) http.Handler {
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{origin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // プリフライトキャッシュ 5 分
	})
}
