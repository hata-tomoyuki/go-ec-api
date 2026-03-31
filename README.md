# BECAUSE - E-Commerce API

Go で構築された EC サイトのバックエンド REST API。

## 技術スタック

- **言語**: Go 1.25
- **ルーター**: chi v5
- **データベース**: PostgreSQL 16
- **マイグレーション**: Goose
- **SQL コード生成**: sqlc
- **認証**: JWT (jwtauth + jwx)
- **テスト**: testcontainers-go
- **トレーシング**: OpenTelemetry

## セットアップ

```bash
# PostgreSQL を起動
docker compose up -d

# マイグレーション実行（Goose で自動適用）
# サーバー起動時に自動実行される

# シードデータ投入
go run ./cmd/seed

# サーバー起動（デフォルト :8080）
go run ./cmd
```

### シードデータ

| 種別 | 内容 |
|------|------|
| ユーザー | 管理者 (`admin@example.com` / `admin1234`)、一般 (`tanaka@example.com` / `password1234`) |
| カテゴリ | メンズファッション、レディースファッション、家電・ガジェット、食品・グルメ、本・書籍 |
| 商品 | 12 商品（各カテゴリに分配） |
| 住所 | 一般ユーザーに 2 件 |
| カート | 3 アイテム |
| 注文 | 3 件（完了・保留・キャンセル） |

## API エンドポイント

### 公開

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/health` | ヘルスチェック |
| GET | `/products` | 商品一覧 |
| GET | `/products/{id}` | 商品詳細 |
| GET | `/categories` | カテゴリ一覧 |
| GET | `/categories/{id}` | カテゴリ詳細 |
| GET | `/categories/{id}/products` | カテゴリ別商品一覧 |

### 認証

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/auth/register` | ユーザー登録 |
| POST | `/auth/login` | ログイン |
| POST | `/auth/refresh` | トークンリフレッシュ |
| POST | `/auth/logout` | ログアウト (要認証) |

### ユーザー (要認証)

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/users/me` | プロフィール取得 |
| PUT | `/users/me` | プロフィール更新 |
| PUT | `/users/me/password` | パスワード変更 |

### カート (要認証)

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/cart` | カート作成 |
| GET | `/cart` | カート取得 |
| POST | `/cart/items` | アイテム追加 |
| PUT | `/cart/items/{id}` | 数量変更 |
| DELETE | `/cart/items/{id}` | アイテム削除 |
| DELETE | `/cart` | カートクリア |

### 注文 (要認証)

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/orders` | 注文作成 |
| GET | `/orders` | 注文一覧 |
| GET | `/orders/{id}` | 注文詳細 |
| PUT | `/orders/{id}/cancel` | 注文キャンセル |

### 住所 (要認証)

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/addresses` | 住所一覧 |
| POST | `/addresses` | 住所作成 |
| GET | `/addresses/{id}` | 住所詳細 |
| PUT | `/addresses/{id}` | 住所更新 |
| DELETE | `/addresses/{id}` | 住所削除 |

### 管理者 (要認証 + admin ロール)

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/products` | 商品作成 |
| PUT | `/products/{id}` | 商品更新 |
| DELETE | `/products/{id}` | 商品削除 |
| POST | `/categories` | カテゴリ作成 |
| PUT | `/categories/{id}` | カテゴリ更新 |
| DELETE | `/categories/{id}` | カテゴリ削除 |
| POST | `/categories/{id}/products` | カテゴリに商品追加 |
| DELETE | `/categories/{id}/products/{productId}` | カテゴリから商品削除 |
| GET | `/admin/orders` | 全注文一覧 |
| PUT | `/admin/orders/{id}/status` | 注文ステータス更新 |

## プロジェクト構成

```
cmd/
  main.go          # エントリーポイント
  api.go           # アプリケーション構造体・サーバー設定
  routes.go        # ルーティング定義
  seed/main.go     # シードデータ投入
internal/
  auth/            # 認証・JWT
  products/        # 商品
  categories/      # カテゴリ
  orders/          # 注文
  carts/           # カート
  address/         # 住所
  adapters/postgresql/
    migrations/    # SQL マイグレーション
    sqlc/          # sqlc 生成コード
  env/             # 環境変数ヘルパー
  json/            # JSON ユーティリティ
  testhelper/      # テストヘルパー
```

## テスト

```bash
make test            # 全テスト（Docker 必須）
make test-unit       # ユニットテストのみ
make test-integration # 統合テストのみ（Docker 必須）
```
