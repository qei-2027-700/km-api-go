# 設計書: AI最適化APIサーバーアーキテクチャ

**プロジェクト**: km-api-go  
**フィーチャーブランチ**: `001-go-echo-gemini`  
**作成日**: 2025-09-07  
**ステータス**: 設計完了  

## 1. アーキテクチャ概要

### 1.1 設計思想
- **AI理解性の最適化**: Gemini、Claude Code等のAIツールが理解しやすいコード構造
- **モジュール結束**: 関連するコンポーネントを1つのディレクトリにまとめる
- **インターフェースと実装の近接**: 従来のクリーンアーキテクチャの分散問題を解決
- **スケーラビリティ**: 小規模から中規模への成長に対応

### 1.2 レイヤーアーキテクチャ
```
Handler Layer (HTTP) 
    ↓
Service Layer (Business Logic)
    ↓  
Repository Layer (Data Access)
    ↓
Database Layer (PostgreSQL)
```

## 2. API仕様

### 2.1 ベースURL
```
http://localhost:8080/api/v1
```

### 2.2 エンドポイント一覧

| Method | エンドポイント | 説明 | リクエスト | レスポンス |
|--------|-------------|------|-----------|-----------|
| GET | `/users` | ユーザー一覧取得 | - | `APIResponse<[]UserResponse>` |
| GET | `/users/:id` | ユーザー詳細取得 | Path: `id` | `APIResponse<UserResponse>` |
| POST | `/users` | ユーザー登録 | Body: `CreateUserRequest` | `APIResponse<UserResponse>` |
| PUT | `/users/:id` | ユーザー更新 | Path: `id`, Body: `UpdateUserRequest` | `APIResponse<UserResponse>` |
| DELETE | `/users/:id` | ユーザー削除 | Path: `id` | `APIResponse<DeleteResponse>` |

### 2.3 共通レスポンス形式
```json
{
  "success": true,
  "data": {...},
  "message": "操作が成功しました",
  "error": null
}
```

## 3. データモデル

### 3.1 エンティティモデル
```go
type User struct {
    ID        uint      `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### 3.2 リクエスト/レスポンスモデル
```go
// リクエストモデル
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=50"`
    Email string `json:"email" validate:"omitempty,email"`
}

// レスポンスモデル
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 4. ディレクトリ構造

### 4.1 AI最適化 + クリーンアーキテクチャ構造
```
├── cmd/
│   └── server/
│       └── main.go              # アプリケーションエントリーポイント
├── internal/
│   ├── user/                    # Userモジュール（AI理解単位）
│   │   ├── domain/              # ドメイン層
│   │   │   ├── model.go         # User エンティティ
│   │   │   ├── model_test.go
│   │   │   ├── repository.go    # リポジトリインターフェース
│   │   │   └── repository_test.go
│   │   ├── usecase/             # ユースケース層
│   │   │   ├── service.go       # サービス実装（インターフェースなし）
│   │   │   └── service_test.go
│   │   ├── infrastructure/      # インフラ層
│   │   │   ├── repository.go    # リポジトリ実装
│   │   │   └── repository_test.go
│   │   └── presentation/        # プレゼンテーション層
│   │       ├── handler.go       # HTTPハンドラー
│   │       ├── handler_test.go
│   │       └── dto.go          # リクエスト/レスポンスDTO
│   ├── common/                  # 共通コンポーネント
│   │   ├── response.go          # APIレスポンス形式
│   │   ├── response_test.go
│   │   ├── database.go          # DB接続設定
│   │   ├── database_test.go
│   │   ├── validator.go         # バリデーション
│   │   └── validator_test.go
│   └── server/                  # サーバー設定
│       ├── middleware.go        # ミドルウェア
│       └── middleware_test.go
├── routes/                      # APIルーティング
│   ├── routes.go                # ルート定義とルーター初期化を統合
│   └── routes_test.go
├── migrations/                  # DBマイグレーション
│   └── 001_create_users.sql
├── docker-compose.yml           # 開発環境
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

### 4.2 AI最適化のポイント
1. **モジュール単位結束**: `internal/user/`内にUser関連の全レイヤーを配置
2. **クリーンアーキテクチャ準拠**: domain → usecase → infrastructure → presentation の依存関係
3. **レイヤー間の近接**: 同一モジュール内で依存関係が明確
4. **インターフェース最小化**: サービス層のインターフェースを省略（小規模最適化）
5. **ルーティング可視性**: `routes/routes.go`でAPIパス一覧を明確化

### 4.3 レイヤー構造の詳細

#### ドメイン層（domain/）
```go
// domain/model.go - エンティティ定義
type User struct {
    ID        uint      `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// domain/repository.go - リポジトリインターフェース
type UserRepository interface {
    GetAll() ([]User, error)
    GetByID(id uint) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}
```

#### ユースケース層（usecase/）
```go
// usecase/service.go - ビジネスロジック実装
type UserService struct {
    repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

#### インフラ層（infrastructure/）
```go
// infrastructure/repository.go - リポジトリ実装
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
    return &userRepository{db: db}
}
```

#### プレゼンテーション層（presentation/）
```go
// presentation/handler.go - HTTPハンドラー
type Handler struct {
    service *usecase.UserService
}

// presentation/dto.go - リクエスト/レスポンスDTO
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 4.4 ルーティング構造

#### routes/routes.go
```go
package routes

import (
    "github.com/labstack/echo/v4"
    "km-api-go/internal/user/presentation"
    "km-api-go/internal/server"
)

// NewRouter はEchoルーターを初期化し、全ルートを設定
func NewRouter(userHandler *presentation.Handler) *echo.Echo {
    e := echo.New()
    
    // ミドルウェア設定
    server.SetupMiddleware(e)
    
    // APIルート設定（Laravel風）
    api := e.Group("/api/v1")
    
    // User API Routes
    api.GET("/users", userHandler.GetUsers)           // ユーザー一覧
    api.GET("/users/:id", userHandler.GetUserByID)    // ユーザー詳細
    api.POST("/users", userHandler.CreateUser)        // ユーザー登録
    api.PUT("/users/:id", userHandler.UpdateUser)     // ユーザー更新
    api.DELETE("/users/:id", userHandler.DeleteUser)  // ユーザー削除
    
    return e
}
```

## 5. データベース設計

### 5.1 テーブル定義
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 5.2 インデックス
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### 5.3 制約
- **name**: 2〜50文字、NOT NULL
- **email**: メール形式、UNIQUE、NOT NULL
- **password**: bcryptハッシュ、NOT NULL
- **timestamps**: 自動更新トリガー設定

## 6. 技術スタック

### 6.1 Backend
- **Language**: Go 1.21+
- **Framework**: Echo v4
- **Database**: PostgreSQL 15+
- **ORM**: GORM v2
- **Validation**: go-playground/validator
- **Password**: bcrypt
- **Environment**: godotenv

### 6.2 Development
- **Container**: Docker Compose
- **Hot Reload**: Air
- **Migration**: カスタムマイグレーション

## 7. 実装順序

### 7.1 Phase 1: 基盤構築
1. プロジェクト初期化 (go mod init)
2. ディレクトリ構造作成
3. Docker Compose設定
4. 基本設定ファイル (.env, air.toml)

### 7.2 Phase 2: データ層
1. データベース接続設定
2. User エンティティ定義
3. マイグレーション実装
4. Repository層実装

### 7.3 Phase 3: ビジネス層
1. Service インターフェース定義
2. Service 実装
3. バリデーション設定
4. エラーハンドリング

### 7.4 Phase 4: プレゼンテーション層
1. Handler実装
2. ルーティング設定
3. ミドルウェア設定
4. レスポンス形式統一

### 7.5 Phase 5: テスト・最適化
1. API動作テスト
2. パフォーマンス確認
3. ドキュメント整備

## 8. 開発コマンド

```bash
# 依存関係インストール
go mod tidy

# 開発サーバー起動（Hot Reload）
air

# Docker環境起動
docker compose up -d

# マイグレーション実行
go run migrations/migrate.go

# 本番ビルド
go build -o bin/server cmd/server/main.go
```

## 9. AI開発支援の考慮事項

### 9.1 コード構造の最適化
- 1ファイル内でインターフェースと実装を配置
- 明確な命名規則でAIが理解しやすく
- モジュール間の依存関係を明確化

### 9.2 ドキュメント連携
- コード内コメントは最小限
- 設計書と仕様書で補完
- 実装例を含めた説明

---

**設計完了**: この設計書に基づいて実装フェーズに移行可能