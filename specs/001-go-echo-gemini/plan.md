# Go Echo API サーバー実装計画書

**プロジェクト名**: km-api-go (AI最適化Go Echo APIサーバー)  
**フィーチャーブランチ**: `001-go-echo-gemini`  
**作成日**: 2025-09-07  
**アーキテクチャ**: AI理解に最適化されたクリーンアーキテクチャ  
**技術スタック**: Go 1.21+, Echo v4, PostgreSQL, GORM v2, Docker Compose, OpenAPI/Swagger

## 実装フェーズと実行順序

### Phase 1: プロジェクト基盤構築 (1-2日)

**目的**: 開発環境とプロジェクト構造の確立

#### 依存関係
- なし（初期フェーズ）

#### 詳細タスク
1. **プロジェクト初期化**
   - `go mod init km-api-go` でGoモジュール初期化
   - 基本ディレクトリ構造作成
   - `.gitignore` 設定

2. **開発環境設定**
   - `docker-compose.yml` 作成（PostgreSQL 15）
   - `.env.example` と `.env` 作成
   - `air.toml` 設定（Hot Reload）

3. **依存関係追加**
   ```bash
   # Core dependencies
   go get github.com/labstack/echo/v4
   go get gorm.io/gorm
   go get gorm.io/driver/postgres
   go get github.com/go-playground/validator/v10
   go get github.com/joho/godotenv
   go get golang.org/x/crypto/bcrypt
   go get github.com/cosmtrek/air
   
   # OpenAPI/Swagger dependencies
   go get github.com/swaggo/swag/cmd/swag
   go get github.com/swaggo/echo-swagger
   go get github.com/swaggo/files
   ```

4. **基本ディレクトリ構造作成**
   ```
   ├── cmd/api/
   ├── internal/
   │   ├── user/
   │   │   └── repository/
   │   ├── company/
   │   │   └── repository/
   │   ├── helper/
   │   └── infra/
   ├── server/
   │   └── middleware/
   ├── migrations/
   └── docs/                    # OpenAPI ドキュメント
   ```

5. **OpenAPI初期設定**
   - `swag init` コマンド実行環境設定
   - Swagger UI ルート設定

#### 成功基準
- Docker Compose でPostgreSQLが起動する
- `air` でHot Reloadが動作する
- Go modulesが正しく設定されている
- 基本ディレクトリ構造が完成している
- Swagger UI にアクセスできる

#### 時間見積もり
約1-2日

---

### Phase 2: データベース・ドメイン層実装 (2-3日)

**目的**: データ基盤とドメインモデルの実装

#### 依存関係
- Phase 1 完了

#### 詳細タスク
1. **データベース接続設定**
   - `internal/infra/database.go` 実装
   - GORM接続設定とDB初期化
   - 接続プール設定

2. **Userドメインモデル実装**
   - `internal/user/domain.go` 作成
   - User構造体定義（ID, Name, Email, Password, Timestamps）
   - GORMタグとJSONタグ設定
   - OpenAPI/Swagger モデルコメント追加

3. **Companyドメインモデル実装**
   - `internal/company/domain.go` 作成
   - Company構造体定義（ID, Name, Email, Address, Timestamps）
   - User-Company関係定義

4. **リポジトリインターフェース定義**
   - `internal/user/repository/interface.go` 作成
   - `internal/company/repository/interface.go` 作成
   - CRUD操作メソッド定義

5. **データベースマイグレーション**
   - `migrations/001_create_users.sql` 作成
   - `migrations/002_create_companies.sql` 作成
   - テーブル作成SQL、インデックス設定
   - マイグレーション実行スクリプト作成

6. **共通コンポーネント**
   - `internal/helper/types.go` 共通型定義
   - `internal/helper/response.go` API レスポンス形式
   - `internal/helper/validator.go` バリデーション設定

#### 成功基準
- データベース接続が確立できる
- User・Company テーブルが正しく作成される
- ドメインモデルが適切に定義されている
- リポジトリインターフェースが明確に定義されている

#### 時間見積もり
約2-3日

---

### Phase 3: インフラストラクチャ層実装 (2-3日)

**目的**: データアクセス層の実装

#### 依存関係
- Phase 2 完了

#### 詳細タスク
1. **Userリポジトリ実装**
   - `internal/user/repository/gorm.go` 作成
   - userRepository 構造体実装
   - GORM を使用したCRUD操作実装

2. **Companyリポジトリ実装**
   - `internal/company/repository/gorm.go` 作成
   - companyRepository 構造体実装
   - CRUD操作とUser関連機能実装

3. **リポジトリメソッド実装**
   - `GetAll()` - エンティティ一覧取得
   - `GetByID()` - 単一エンティティ取得
   - `Create()` - エンティティ作成
   - `Update()` - エンティティ更新
   - `Delete()` - エンティティ削除

4. **データベーストランザクション対応**
   - トランザクション処理の実装
   - ロールバック処理

5. **リポジトリテスト作成**
   - `repository_test.go` 作成
   - 各メソッドの単体テスト

#### 成功基準
- リポジトリがインターフェースを実装している
- 全てのCRUD操作が正常に動作する
- エラーハンドリングが適切に実装されている
- 単体テストが通る

#### 時間見積もり
約2-3日

---

### Phase 4: ビジネスロジック層実装 (2-3日)

**目的**: アプリケーションのビジネスロジック実装

#### 依存関係
- Phase 3 完了

#### 詳細タスク
1. **Userユースケース実装**
   - `internal/user/usecase.go` 作成
   - UserUsecase 構造体実装
   - 依存性注入（リポジトリ）

2. **Companyユースケース実装**
   - `internal/company/usecase.go` 作成
   - 複数リポジトリの組み合わせ処理

3. **ビジネスロジック実装**
   - ユーザー作成時のパスワードハッシュ化
   - メール重複チェック
   - バリデーション処理
   - エラーハンドリング

4. **ユースケースメソッド実装**
   - `GetAllUsers()` / `GetAllCompanies()`
   - `GetByID()` - 詳細取得
   - `Create()` - 新規作成
   - `Update()` - 更新
   - `Delete()` - 削除

5. **複合処理実装**
   - User-Company関連処理
   - 複数エンティティの同期処理

6. **パスワード処理**
   - bcrypt を使用したハッシュ化
   - パスワード強度バリデーション

7. **ユースケーステスト作成**
   - `usecase_test.go` 作成
   - ビジネスロジックの単体テスト
   - モックリポジトリを使用

#### 成功基準
- 全てのビジネスロジックが正しく実装されている
- パスワードが安全にハッシュ化される
- バリデーションが適切に動作する
- 単体テストが通る

#### 時間見積もり
約2-3日

---

### Phase 5: プレゼンテーション層・OpenAPI実装 (3-4日)

**目的**: HTTPハンドラーとAPI エンドポイント、OpenAPIドキュメントの実装

#### 依存関係
- Phase 4 完了

#### 詳細タスク
1. **DTOモデル作成**
   - `internal/user/dto.go` 作成
   - `internal/company/dto.go` 作成
   - リクエストDTO（CreateUserRequest, UpdateUserRequest等）
   - レスポンスDTO（UserResponse, CompanyResponse等）
   - バリデーションタグ設定
   - **OpenAPIコメント追加**

2. **User HTTPハンドラー実装**
   - `internal/user/handler.go` 作成
   - Handler 構造体とコンストラクタ
   - **Swaggerアノテーション付きハンドラーメソッド**

3. **User エンドポイント実装（OpenAPI対応）**
   ```go
   // @Summary ユーザー一覧取得
   // @Description 全てのユーザーを取得します
   // @Tags users
   // @Accept json
   // @Produce json
   // @Success 200 {object} helper.APIResponse{data=[]dto.UserResponse}
   // @Failure 500 {object} helper.APIResponse
   // @Router /api/v1/users [get]
   ```
   - `GET /api/v1/users` - ユーザー一覧
   - `GET /api/v1/users/:id` - ユーザー詳細
   - `POST /api/v1/users` - ユーザー登録
   - `PUT /api/v1/users/:id` - ユーザー更新
   - `DELETE /api/v1/users/:id` - ユーザー削除

4. **Company HTTPハンドラー・エンドポイント実装**
   - `internal/company/handler.go` 作成
   - 同様のCRUD エンドポイント
   - **Swaggerアノテーション**

5. **リクエスト/レスポンス処理**
   - JSONバインド
   - バリデーション
   - エラーレスポンス
   - 統一されたAPIレスポンス形式

6. **OpenAPIドキュメント生成**
   - `swag init -g cmd/api/main.go` 実行
   - `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml` 生成
   - Swagger UI エンドポイント設定 (`/swagger/*`)

7. **ハンドラーテスト作成**
   - `handler_test.go` 作成
   - Echo のテストヘルパー使用
   - 各エンドポイントのテスト

#### 成功基準
- 全てのAPIエンドポイントが正常に動作する
- OpenAPIドキュメントが自動生成される
- Swagger UI でAPI仕様が確認できる
- リクエスト/レスポンスが仕様通りに処理される
- バリデーションが適切に動作する
- HTTPテストが通る

#### 時間見積もり
約3-4日

---

### Phase 6: ルーティング・サーバー設定実装 (1-2日)

**目的**: APIルーティングとサーバー設定の実装

#### 依存関係
- Phase 5 完了

#### 詳細タスク
1. **ルーティング実装**
   - `server/route.go` 作成
   - NewRouter 関数実装
   - APIルート設定
   - ハンドラーの依存性注入
   - **Swagger UI ルート追加** (`/swagger/*`)

2. **ミドルウェア設定**
   - `server/middleware/cors.go` - CORS設定
   - `server/middleware/logger.go` - ログミドルウェア
   - `server/middleware/auth.go` - 認証（将来用）
   - リカバリーミドルウェア
   - リクエストID生成

3. **メインサーバー実装**
   - `cmd/api/main.go` 作成
   - **Swaggerメインコメント追加**
   ```go
   // @title km-api-go API
   // @version 1.0
   // @description AI最適化Go Echo APIサーバー
   // @host localhost:8080
   // @BasePath /api/v1
   ```
   - 依存性注入とワイヤリング
   - サーバー起動処理
   - グレースフルシャットダウン

4. **設定管理**
   - 環境変数の読み込み
   - サーバー設定（ポート、タイムアウト等）
   - ログレベル設定

#### 成功基準
- サーバーが正常に起動する
- 全てのAPIエンドポイントにアクセスできる
- Swagger UI が `/swagger/` でアクセス可能
- ミドルウェアが適切に動作する
- 環境変数が正しく読み込まれる

#### 時間見積もり
約1-2日

---

### Phase 7: 統合テスト・API動作確認 (2-3日)

**目的**: システム全体の統合テストと動作確認

#### 依存関係
- Phase 6 完了

#### 詳細タスク
1. **APIテスト実装**
   - `server/route_test.go` 作成
   - 統合テスト実装
   - テストデータベース設定

2. **エンドツーエンドテスト**
   - 各エンドポイントの正常系テスト
   - 異常系テスト（バリデーションエラー、404等）
   - データベース連携テスト

3. **OpenAPI仕様確認**
   - **Swagger UI でのAPI仕様確認**
   - **生成されたJSONスキーマの妥当性確認**
   - レスポンス形式とドキュメントの整合性チェック

4. **API動作確認**
   - Postman/curl でのAPI確認
   - **Swagger UI からのAPI実行テスト**
   - レスポンス形式チェック
   - エラーハンドリング確認

5. **パフォーマンステスト**
   - 基本的な負荷テスト
   - メモリ使用量確認
   - データベース接続数確認

6. **セキュリティチェック**
   - SQLインジェクション対策確認
   - パスワードハッシュ化確認
   - バリデーション抜け確認

7. **ドキュメント整備**
   - `README.md` 更新
   - API使用方法の記載
   - OpenAPI仕様書の配布準備

#### 成功基準
- 全ての統合テストが通る
- APIが仕様通りに動作する
- OpenAPIドキュメントが正確に生成される
- Swagger UI でAPI操作が可能
- パフォーマンスが許容範囲内
- セキュリティ要件を満たしている

#### 時間見積もり
約2-3日

---

### Phase 8: ドキュメント・デプロイ準備 (1-2日)

**目的**: ドキュメント整備とデプロイメント準備

#### 依存関係
- Phase 7 完了

#### 詳細タスク
1. **OpenAPIドキュメント最終調整**
   - スキーマ検証・修正
   - `docs/swagger.yaml` の手動調整（必要に応じて）
   - API使用例の追加

2. **プロジェクトドキュメント作成**
   - `README.md` 完成版作成
   - 環境構築手順
   - API使用方法
   - OpenAPI/Swagger アクセス方法

3. **デプロイメント準備**
   - `Dockerfile` 作成（必要に応じて）
   - 本番環境用設定
   - ヘルスチェック エンドポイント追加

4. **CI/CDパイプライン設定**
   - GitHub Actions 設定（オプション）
   - テスト自動実行
   - OpenAPIドキュメント自動生成

#### 成功基準
- 完全なAPIドキュメントが生成されている
- デプロイメントが容易に実行できる
- 開発者が容易にプロジェクトを理解・継続できる

#### 時間見積もり
約1-2日

---

## OpenAPI/Swagger統合の詳細

### 1. **依存関係**
```bash
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/echo-swagger
go get github.com/swaggo/files
```

### 2. **アノテーション例**
```go
// @Summary ユーザー作成
// @Description 新しいユーザーを作成します
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "ユーザー情報"
// @Success 201 {object} helper.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} helper.APIResponse
// @Failure 500 {object} helper.APIResponse
// @Router /api/v1/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
    // 実装
}
```

### 3. **生成されるファイル**
- `docs/docs.go` - Go 用の埋め込み仕様書
- `docs/swagger.json` - JSON形式
- `docs/swagger.yaml` - YAML形式

### 4. **Swagger UI アクセス**
- URL: `http://localhost:8080/swagger/`
- リアルタイムでAPI操作とテストが可能

---

## 成功基準まとめ

### 技術的成功基準
1. **アーキテクチャ準拠**: クリーンアーキテクチャの依存関係が正しく実装されている
2. **AI理解性**: 各モジュール内でコンポーネントの関係性が明確
3. **OpenAPI準拠**: 完全なAPI仕様書が自動生成される
4. **テスト品質**: 単体テスト・統合テストのカバレッジ80%以上
5. **パフォーマンス**: 1000 req/sec の負荷で安定動作
6. **セキュリティ**: OWASP推奨事項に準拠

### ビジネス成功基準
1. **API仕様準拠**: 設計書通りのエンドポイントとレスポンス
2. **拡張性**: 新しいドメインモジュール追加が容易
3. **保守性**: AI支援による効率的な開発・保守が可能
4. **開発効率**: 後続機能開発が50%以上効率化
5. **ドキュメント品質**: OpenAPIによる自動ドキュメント生成

## リスク管理

### 高リスク要因
1. **GORM学習コスト**: GORM の使用経験が少ない場合
2. **テストデータベース**: テスト環境でのDB分離
3. **OpenAPI統合**: Swaggerアノテーションの学習コスト
4. **マイグレーション**: 本番環境でのスキーマ変更

### 軽減策
1. **早期プロトタイピング**: Phase 2 でGORM動作確認
2. **テスト環境分離**: Docker Compose でテスト用DB準備
3. **段階的OpenAPI導入**: 基本実装後にアノテーション追加
4. **段階的マイグレーション**: 破壊的変更を避ける設計

## 総実装期間見積もり

**合計期間**: 16-22日（約3-4週間）

- **Phase 1**: 1-2日（OpenAPI環境設定含む）
- **Phase 2**: 2-3日  
- **Phase 3**: 2-3日
- **Phase 4**: 2-3日
- **Phase 5**: 3-4日（OpenAPI統合含む）
- **Phase 6**: 1-2日
- **Phase 7**: 2-3日（OpenAPI検証含む）
- **Phase 8**: 1-2日

この実装計画により、AI最適化されたクリーンアーキテクチャで、かつOpenAPI/Swaggerによる完全なAPIドキュメントを持つスケーラブルなGo Echo APIサーバーを段階的に構築できます。各フェーズでの成功基準を満たすことで、最終的に保守性・拡張性・可読性に優れたシステムが完成します。