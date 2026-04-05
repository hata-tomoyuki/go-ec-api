package middleware

import "net/http"

// SecurityHeaders は全レスポンスにセキュリティ関連の HTTP ヘッダーを付与するミドルウェア。
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// MIME スニッフィング防止: ブラウザが Content-Type を無視して内容を推測するのを防ぐ
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// クリックジャッキング防止: iframe での埋め込みを全面禁止
		w.Header().Set("X-Frame-Options", "DENY")
		// レガシー XSS フィルター無効化: 誤検知による脆弱性を防ぐ（CSP での対策が推奨）
		w.Header().Set("X-XSS-Protection", "0")
		// リファラー制御: 同一オリジンではフルURL、クロスオリジンではオリジンのみ送信
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
