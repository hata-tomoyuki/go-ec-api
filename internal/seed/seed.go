package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Run executes seed data insertion into the database.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("トランザクション開始に失敗: %w", err)
	}
	defer tx.Rollback(ctx)

	// --------------------------------------------------
	// 既存のシードデータをクリア（依存順で削除）
	// --------------------------------------------------
	tables := []string{
		"cart_items", "carts",
		"order_items", "orders",
		"product_categories",
		"addresses",
		"refresh_tokens", "revoked_tokens",
		"products", "categories", "users",
	}
	for _, t := range tables {
		if _, err := tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s", t)); err != nil {
			return fmt.Errorf("%s のクリアに失敗: %w", t, err)
		}
	}

	// --------------------------------------------------
	// Users
	// --------------------------------------------------
	adminHash := hashPassword("admin1234")
	userHash := hashPassword("password1234")

	var adminID, userID int64

	err = tx.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, 'admin') RETURNING id`,
		"管理者", "admin@example.com", adminHash,
	).Scan(&adminID)
	if err != nil {
		return fmt.Errorf("管理者ユーザー作成に失敗: %w", err)
	}
	fmt.Printf("  管理者ユーザー作成 (id=%d)\n", adminID)

	err = tx.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, 'user') RETURNING id`,
		"田中 太郎", "tanaka@example.com", userHash,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("一般ユーザー作成に失敗: %w", err)
	}
	fmt.Printf("  一般ユーザー作成 (id=%d)\n", userID)

	// --------------------------------------------------
	// Categories
	// --------------------------------------------------
	type cat struct {
		name       string
		desc       string
		imageColor string
	}
	categories := []cat{
		{"メンズファッション", "メンズ向けの衣類・アクセサリー", "from-blue-600 to-blue-800"},
		{"レディースファッション", "レディース向けの衣類・アクセサリー", "from-pink-500 to-rose-600"},
		{"家電・ガジェット", "最新の家電製品やガジェット", "from-slate-600 to-slate-800"},
		{"食品・グルメ", "厳選された食品・飲料", "from-amber-500 to-orange-600"},
		{"本・書籍", "話題の書籍・専門書", "from-emerald-600 to-teal-700"},
	}

	catIDs := make([]int64, len(categories))
	for i, c := range categories {
		err = tx.QueryRow(ctx,
			`INSERT INTO categories (name, description, image_color) VALUES ($1, $2, $3) RETURNING id`,
			c.name, c.desc, c.imageColor,
		).Scan(&catIDs[i])
		if err != nil {
			return fmt.Errorf("カテゴリ %q の作成に失敗: %w", c.name, err)
		}
	}
	fmt.Printf("  カテゴリ %d 件作成\n", len(categories))

	// --------------------------------------------------
	// Products
	// --------------------------------------------------
	type prod struct {
		name         string
		description  string
		priceInCents int32
		quantity     int32
		imageColor   string
		categoryIdx  int
	}
	products := []prod{
		// メンズファッション (catIDs[0])
		{"プレミアムコットンTシャツ", "上質なコットン素材を使用した、肌触りの良いベーシックTシャツ。シンプルなデザインで普段使いからジャケットのインナーまで幅広く活躍します。", 4980, 50, "from-sky-400 to-blue-600", 0},
		{"スリムフィットデニムパンツ", "美しいシルエットにこだわったスリムフィットデニム。程よいストレッチ性があり、快適な履き心地と洗練された印象を両立した一本です。", 8980, 30, "from-indigo-500 to-blue-700", 0},
		{"レザーウォレット", "高級感のあるレザーを使用した上品なウォレット。使うほどに風合いが増し、収納力とデザイン性を兼ね備えた長く愛用できるアイテムです。", 12800, 20, "from-amber-700 to-orange-900", 0},

		// レディースファッション (catIDs[1])
		{"フローラルワンピース", "華やかなフローラル柄が魅力のワンピース。軽やかな着心地で、デイリーからお出かけまで女性らしいスタイルを演出します。", 6980, 25, "from-pink-400 to-rose-500", 1},
		{"カシミアニットセーター", "やわらかなカシミア素材を使用した、上品で暖かなニットセーター。シンプルなデザインで秋冬のコーディネートに高級感をプラスします。", 19800, 15, "from-rose-400 to-pink-600", 1},

		// 家電・ガジェット (catIDs[2])
		{"ワイヤレスノイズキャンセリングヘッドホン", "高性能ノイズキャンセリング機能を搭載したワイヤレスヘッドホン。クリアな音質と快適な装着感で、通勤や作業時間をより充実させます。", 34800, 40, "from-slate-500 to-gray-700", 2},
		{"スマートウォッチ Pro", "健康管理や通知機能を備えた多機能スマートウォッチ。スタイリッシュなデザインと高い実用性で、日常生活をスマートにサポートします。", 49800, 35, "from-gray-600 to-slate-800", 2},
		{"ポータブルBluetoothスピーカー", "コンパクトながら迫力あるサウンドを楽しめるBluetoothスピーカー。持ち運びしやすく、自宅でもアウトドアでも活躍します。", 12800, 60, "from-zinc-500 to-slate-700", 2},

		// 食品・グルメ (catIDs[3])
		{"宇治抹茶スイーツセット", "香り高い宇治抹茶を贅沢に使用したスイーツの詰め合わせ。上品な甘さとほろ苦さが楽しめる、ギフトにも人気のセットです。", 3980, 100, "from-amber-400 to-orange-500", 3},
		{"国産黒毛和牛 すき焼きセット", "厳選した国産黒毛和牛を使用した贅沢なすき焼きセット。とろけるような旨みとやわらかさをご家庭で手軽に楽しめます。", 9800, 20, "from-orange-500 to-red-600", 3},

		// 本・書籍 (catIDs[4])
		{"はじめてのプログラミング入門", "プログラミングの基礎をやさしく学べる初心者向け入門書。図解やサンプルコードが豊富で、初めて学ぶ方にもわかりやすい一冊です。", 2980, 200, "from-emerald-500 to-teal-600", 4},
		{"AI時代の働き方改革", "AI時代に求められる仕事の進め方やキャリア形成を解説した実践書。変化の大きい時代を生き抜くためのヒントが詰まっています。", 1980, 150, "from-teal-500 to-emerald-700", 4},
	}

	prodIDs := make([]int64, len(products))
	for i, p := range products {
		err = tx.QueryRow(ctx,
			`INSERT INTO products (name, description, price_in_cents, quantity, image_color) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			p.name, p.description, p.priceInCents, p.quantity, p.imageColor,
		).Scan(&prodIDs[i])
		if err != nil {
			return fmt.Errorf("商品 %q の作成に失敗: %w", p.name, err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)`,
			prodIDs[i], catIDs[p.categoryIdx],
		)
		if err != nil {
			return fmt.Errorf("商品-カテゴリ紐付けに失敗: %w", err)
		}
	}
	fmt.Printf("  商品 %d 件作成（カテゴリ紐付け済み）\n", len(products))

	// --------------------------------------------------
	// Addresses（一般ユーザー用）
	// --------------------------------------------------
	type addr struct {
		street, city, state, zipCode, country string
	}
	addresses := []addr{
		{"神宮前1-2-3 ABCマンション 401号室", "渋谷区", "東京都", "150-0001", "日本"},
		{"梅田4-5-6", "大阪市北区", "大阪府", "530-0001", "日本"},
	}

	addrIDs := make([]int32, len(addresses))
	for i, a := range addresses {
		err = tx.QueryRow(ctx,
			`INSERT INTO addresses (user_id, street, city, state, zip_code, country)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			userID, a.street, a.city, a.state, a.zipCode, a.country,
		).Scan(&addrIDs[i])
		if err != nil {
			return fmt.Errorf("住所作成に失敗: %w", err)
		}
	}
	fmt.Printf("  住所 %d 件作成\n", len(addresses))

	// --------------------------------------------------
	// Cart（一般ユーザー用）
	// --------------------------------------------------
	var cartID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO carts (user_id) VALUES ($1) RETURNING id`,
		userID,
	).Scan(&cartID)
	if err != nil {
		return fmt.Errorf("カート作成に失敗: %w", err)
	}

	cartItems := []struct {
		prodIdx  int
		quantity int32
	}{
		{0, 2}, // プレミアムコットンTシャツ
		{5, 1}, // ワイヤレスノイズキャンセリングヘッドホン
		{8, 3}, // 宇治抹茶スイーツセット
	}
	for _, ci := range cartItems {
		_, err = tx.Exec(ctx,
			`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3)`,
			cartID, prodIDs[ci.prodIdx], ci.quantity,
		)
		if err != nil {
			return fmt.Errorf("カートアイテム追加に失敗: %w", err)
		}
	}
	fmt.Printf("  カートアイテム %d 件作成\n", len(cartItems))

	// --------------------------------------------------
	// Orders（一般ユーザー用）
	// --------------------------------------------------

	// 注文1: 配達完了（Tシャツ×2 + 抹茶スイーツ×1）
	var order1ID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (customer_id, status) VALUES ($1, 'completed') RETURNING id`,
		userID,
	).Scan(&order1ID)
	if err != nil {
		return fmt.Errorf("注文1作成に失敗: %w", err)
	}
	for _, oi := range []struct {
		prodIdx      int
		quantity     int32
		priceInCents int32
	}{
		{0, 2, 4980}, // Tシャツ×2
		{8, 1, 3980}, // 抹茶スイーツ×1
	} {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
			 VALUES ($1, $2, $3, $4)`,
			order1ID, prodIDs[oi.prodIdx], oi.quantity, oi.priceInCents,
		)
		if err != nil {
			return fmt.Errorf("注文1アイテム追加に失敗: %w", err)
		}
	}

	// 注文2: 保留中（スマートウォッチ×1）
	var order2ID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (customer_id, status) VALUES ($1, 'pending') RETURNING id`,
		userID,
	).Scan(&order2ID)
	if err != nil {
		return fmt.Errorf("注文2作成に失敗: %w", err)
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
		 VALUES ($1, $2, $3, $4)`,
		order2ID, prodIDs[6], int32(1), int32(49800),
	)
	if err != nil {
		return fmt.Errorf("注文2アイテム追加に失敗: %w", err)
	}

	// 注文3: キャンセル済み（AI本×1 + ワンピース×1）
	var order3ID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (customer_id, status) VALUES ($1, 'cancelled') RETURNING id`,
		userID,
	).Scan(&order3ID)
	if err != nil {
		return fmt.Errorf("注文3作成に失敗: %w", err)
	}
	for _, oi := range []struct {
		prodIdx      int
		quantity     int32
		priceInCents int32
	}{
		{11, 1, 1980}, // AI本×1
		{3, 1, 6980},  // ワンピース×1
	} {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
			 VALUES ($1, $2, $3, $4)`,
			order3ID, prodIDs[oi.prodIdx], oi.quantity, oi.priceInCents,
		)
		if err != nil {
			return fmt.Errorf("注文3アイテム追加に失敗: %w", err)
		}
	}
	fmt.Printf("  注文 3 件作成（注文アイテム含む）\n")

	// --------------------------------------------------
	// コミット
	// --------------------------------------------------
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("コミットに失敗: %w", err)
	}

	// --------------------------------------------------
	// サマリー出力
	// --------------------------------------------------
	fmt.Println()
	fmt.Println("=== シードデータ サマリー ===")
	fmt.Printf("ユーザー:\n")
	fmt.Printf("  管理者: admin@example.com / admin1234\n")
	fmt.Printf("  一般:   tanaka@example.com / password1234\n")
	fmt.Printf("カテゴリ: %d 件\n", len(categories))
	fmt.Printf("商品:     %d 件\n", len(products))
	fmt.Printf("住所:     %d 件 (一般ユーザー)\n", len(addresses))
	fmt.Printf("カート:   %d アイテム (一般ユーザー)\n", len(cartItems))
	fmt.Printf("注文:     3 件 (completed, pending, cancelled)\n")

	return nil
}

func hashPassword(plain string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("パスワードハッシュに失敗:", err)
	}
	return string(hash)
}
