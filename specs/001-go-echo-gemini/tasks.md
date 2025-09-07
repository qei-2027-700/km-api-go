# Go Echo API サーバー開発タスク詳細分解

**プロジェクト名**: km-api-go (AI最適化Go Echo APIサーバー)  
**フィーチャーブランチ**: `001-go-echo-gemini`  
**作成日**: 2025-09-07  

実装計画書を基に、各フェーズを具体的で実行可能なタスクに分解しました。各タスクは2-4時間以内で完了できるように設計されています。

## **Phase 1: プロジェクト基盤構築**

### **Task 1.1: Go モジュール初期化と基本構造作成**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: なし

**具体的なアクション**:
```bash
# Go モジュール初期化
cd /Users/km/web-dev/_github/km-api-go
go mod init km-api-go

# 基本ディレクトリ構造作成
mkdir -p cmd/api
mkdir -p internal/{user,company,helper,infra}/{repository,}
mkdir -p server/middleware
mkdir -p migrations
mkdir -p docs
```

**成果物**:
- `go.mod` ファイル
- プロジェクトディレクトリ構造

**受け入れ基準**:
- [ ] `go mod init` が正常に完了している
- [ ] 計画書通りのディレクトリ構造が作成されている

---

### **Task 1.2: 依存関係インストール**
- **優先度**: High  
- **推定時間**: 1時間
- **依存関係**: Task 1.1

**具体的なアクション**:
```bash
# コア依存関係
go get github.com/labstack/echo/v4
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get github.com/cosmtrek/air

# OpenAPI/Swagger依存関係
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/echo-swagger
go get github.com/swaggo/files
```

**受け入れ基準**:
- [ ] 全ての依存関係が `go.mod` に追加されている
- [ ] `go mod tidy` が正常に実行される

---

### **Task 1.3: 開発環境設定ファイル作成**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 1.2

**具体的なアクション**:
1. `docker-compose.yml` 作成（PostgreSQL 15設定）
2. `.env.example` と `.env` 作成
3. `air.toml` 作成（Hot Reload設定）
4. `.gitignore` 作成

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/docker-compose.yml`
- `/Users/km/web-dev/_github/km-api-go/.env.example`
- `/Users/km/web-dev/_github/km-api-go/.env`
- `/Users/km/web-dev/_github/km-api-go/air.toml`
- `/Users/km/web-dev/_github/km-api-go/.gitignore`

**受け入れ基準**:
- [ ] `docker compose up -d` でPostgreSQLが起動する
- [ ] `.env` ファイルから環境変数が読み込める
- [ ] `air.toml` の設定が適切に構成されている

---

## **Phase 2: データベース・ドメイン層実装**

### **Task 2.1: データベース接続設定実装**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 1.3

**具体的なアクション**:
1. `internal/infra/database.go` 作成
2. GORM接続設定とDB初期化関数実装
3. 接続プール設定

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/infra/database.go`

**受け入れ基準**:
- [ ] データベース接続が確立できる
- [ ] エラーハンドリングが実装されている
- [ ] 接続プールが適切に設定されている

---

### **Task 2.2: Userドメインモデル実装**
- **優先度**: High
- **推定時間**: 1.5時間
- **依存関係**: Task 2.1

**具体的なアクション**:
1. `internal/user/domain.go` 作成
2. User構造体定義（ID, Name, Email, Password, Timestamps）
3. GORMタグとJSONタグ設定
4. OpenAPI/Swagger モデルコメント追加

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/domain.go`

**受け入れ基準**:
- [ ] User構造体が適切に定義されている
- [ ] GORMタグが正しく設定されている
- [ ] JSONタグが設定されている
- [ ] Swaggerコメントが追加されている

---

### **Task 2.3: Companyドメインモデル実装**
- **優先度**: High
- **推定時間**: 1.5時間
- **依存関係**: Task 2.2

**具体的なアクション**:
1. `internal/company/domain.go` 作成
2. Company構造体定義
3. User-Company関係定義

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/company/domain.go`

**受け入れ基準**:
- [ ] Company構造体が定義されている
- [ ] User との関係が適切に定義されている
- [ ] 必要なタグが設定されている

---

### **Task 2.4: リポジトリインターフェース定義**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 2.3

**具体的なアクション**:
1. `internal/user/repository/interface.go` 作成
2. `internal/company/repository/interface.go` 作成
3. CRUD操作メソッド定義

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/repository/interface.go`
- `/Users/km/web-dev/_github/km-api-go/internal/company/repository/interface.go`

**受け入れ基準**:
- [ ] 全てのCRUD操作メソッドが定義されている
- [ ] インターフェースが適切に設計されている
- [ ] エラーハンドリングが考慮されている

---

### **Task 2.5: データベースマイグレーション作成**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 2.4

**具体的なアクション**:
1. `migrations/001_create_users.sql` 作成
2. `migrations/002_create_companies.sql` 作成
3. テーブル作成SQL、インデックス設定
4. マイグレーション実行スクリプト作成

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/migrations/001_create_users.sql`
- `/Users/km/web-dev/_github/km-api-go/migrations/002_create_companies.sql`
- マイグレーション実行スクリプト

**受け入れ基準**:
- [ ] テーブルが正しく作成される
- [ ] 適切なインデックスが設定されている
- [ ] マイグレーションが実行できる

---

### **Task 2.6: 共通コンポーネント実装**
- **優先度**: Medium
- **推定時間**: 2時間
- **依存関係**: Task 2.4

**具体的なアクション**:
1. `internal/helper/types.go` 共通型定義
2. `internal/helper/response.go` API レスポンス形式
3. `internal/helper/validator.go` バリデーション設定

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/helper/types.go`
- `/Users/km/web-dev/_github/km-api-go/internal/helper/response.go`
- `/Users/km/web-dev/_github/km-api-go/internal/helper/validator.go`

**受け入れ基準**:
- [ ] 共通型が適切に定義されている
- [ ] APIレスポンス形式が統一されている
- [ ] バリデーション設定が動作する

---

## **Phase 3: インフラストラクチャ層実装**

### **Task 3.1: Userリポジトリ実装**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 2.6

**具体的なアクション**:
1. `internal/user/repository/gorm.go` 作成
2. userRepository構造体実装
3. CRUD操作実装（GetAll, GetByID, Create, Update, Delete）

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/repository/gorm.go`

**受け入れ基準**:
- [ ] インターフェースが実装されている
- [ ] 全てのCRUD操作が正常に動作する
- [ ] エラーハンドリングが実装されている

---

### **Task 3.2: Companyリポジトリ実装**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 3.1

**具体的なアクション**:
1. `internal/company/repository/gorm.go` 作成
2. companyRepository構造体実装
3. CRUD操作とUser関連機能実装

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/company/repository/gorm.go`

**受け入れ基準**:
- [ ] CRUD操作が正常に動作する
- [ ] User関連機能が実装されている
- [ ] リレーション処理が適切に実装されている

---

### **Task 3.3: リポジトリ単体テスト作成**
- **優先度**: Medium
- **推定時間**: 4時間
- **依存関係**: Task 3.2

**具体的なアクション**:
1. `internal/user/repository/gorm_test.go` 作成
2. `internal/company/repository/gorm_test.go` 作成
3. 各メソッドの単体テスト実装

**成果物**:
- テストファイル2つ

**受け入れ基準**:
- [ ] 全てのリポジトリメソッドがテストされている
- [ ] テストが正常に実行される
- [ ] カバレッジが80%以上

---

## **Phase 4: ビジネスロジック層実装**

### **Task 4.1: Userユースケース実装**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 3.2

**具体的なアクション**:
1. `internal/user/usecase.go` 作成
2. UserUsecase構造体実装
3. ビジネスロジック実装（パスワードハッシュ化、バリデーション等）

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/usecase.go`

**受け入れ基準**:
- [ ] 全てのビジネスロジックが実装されている
- [ ] パスワードが安全にハッシュ化される
- [ ] バリデーションが適切に動作する

---

### **Task 4.2: Companyユースケース実装**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 4.1

**具体的なアクション**:
1. `internal/company/usecase.go` 作成
2. 複数リポジトリの組み合わせ処理実装

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/company/usecase.go`

**受け入れ基準**:
- [ ] ビジネスロジックが適切に実装されている
- [ ] リポジトリ間の依存関係が管理されている

---

### **Task 4.3: ユースケース単体テスト作成**
- **優先度**: Medium
- **推定時間**: 4時間
- **依存関係**: Task 4.2

**具体的なアクション**:
1. `internal/user/usecase_test.go` 作成
2. `internal/company/usecase_test.go` 作成
3. モックリポジトリを使用したテスト実装

**成果物**:
- テストファイル2つ

**受け入れ基準**:
- [ ] ビジネスロジックがテストされている
- [ ] モックが適切に使用されている
- [ ] テストが正常に実行される

---

## **Phase 5: プレゼンテーション層・OpenAPI実装**

### **Task 5.1: DTO（Data Transfer Object）作成**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 4.2

**具体的なアクション**:
1. `internal/user/dto.go` 作成
2. `internal/company/dto.go` 作成
3. リクエスト・レスポンスDTO定義
4. バリデーションタグ設定
5. OpenAPIコメント追加

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/dto.go`
- `/Users/km/web-dev/_github/km-api-go/internal/company/dto.go`

**受け入れ基準**:
- [ ] 全てのDTO構造体が定義されている
- [ ] バリデーションタグが適切に設定されている
- [ ] Swaggerコメントが追加されている

---

### **Task 5.2: User HTTPハンドラー実装**
- **優先度**: High
- **推定時間**: 4時間
- **依存関係**: Task 5.1

**具体的なアクション**:
1. `internal/user/handler.go` 作成
2. Handler構造体とコンストラクタ実装
3. Swaggerアノテーション付きハンドラーメソッド実装
4. 5つのCRUDエンドポイント実装

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/user/handler.go`

**受け入れ基準**:
- [ ] 全てのエンドポイントが実装されている
- [ ] Swaggerアノテーションが適切に設定されている
- [ ] エラーハンドリングが実装されている
- [ ] レスポンス形式が統一されている

---

### **Task 5.3: Company HTTPハンドラー実装**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 5.2

**具体的なアクション**:
1. `internal/company/handler.go` 作成
2. 同様のCRUD エンドポイント実装
3. Swaggerアノテーション設定

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/internal/company/handler.go`

**受け入れ基準**:
- [ ] 全てのエンドポイントが実装されている
- [ ] User ハンドラーと一貫性がある

---

### **Task 5.4: OpenAPIドキュメント生成設定**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 5.3

**具体的なアクション**:
1. swag コマンド実行環境設定
2. メインファイルにSwaggerコメント追加
3. `swag init` 実行とドキュメント生成確認

**成果物**:
- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

**受け入れ基準**:
- [ ] OpenAPIドキュメントが自動生成される
- [ ] 全てのエンドポイントがドキュメントに含まれている
- [ ] スキーマが正しく定義されている

---

### **Task 5.5: ハンドラー単体テスト作成**
- **優先度**: Medium
- **推定時間**: 4時間
- **依存関係**: Task 5.3

**具体的なアクション**:
1. `internal/user/handler_test.go` 作成
2. `internal/company/handler_test.go` 作成
3. Echoテストヘルパー使用

**成果物**:
- テストファイル2つ

**受け入れ基準**:
- [ ] 全てのエンドポイントがテストされている
- [ ] HTTPステータスコードが適切にテストされている

---

## **Phase 6: ルーティング・サーバー設定実装**

### **Task 6.1: ミドルウェア実装**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 5.3

**具体的なアクション**:
1. `server/middleware/cors.go` 作成
2. `server/middleware/logger.go` 作成
3. `server/middleware/auth.go` 作成（将来用）
4. リカバリーミドルウェア設定

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/server/middleware/` 以下のファイル

**受け入れ基準**:
- [ ] CORS設定が適切に動作する
- [ ] ログミドルウェアが動作する
- [ ] リカバリーミドルウェアが動作する

---

### **Task 6.2: ルーティング実装**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 6.1

**具体的なアクション**:
1. `server/route.go` 作成
2. NewRouter関数実装
3. APIルート設定
4. ハンドラーの依存性注入
5. Swagger UIルート追加

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/server/route.go`

**受け入れ基準**:
- [ ] 全てのAPIルートが設定されている
- [ ] Swagger UIにアクセスできる
- [ ] 依存性注入が適切に実装されている

---

### **Task 6.3: メインサーバー実装**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 6.2

**具体的なアクション**:
1. `cmd/api/main.go` 作成
2. Swaggerメインコメント追加
3. 依存性注入とワイヤリング実装
4. サーバー起動処理実装
5. グレースフルシャットダウン実装

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/cmd/api/main.go`

**受け入れ基準**:
- [ ] サーバーが正常に起動する
- [ ] 環境変数が正しく読み込まれる
- [ ] グレースフルシャットダウンが動作する

---

## **Phase 7: 統合テスト・API動作確認**

### **Task 7.1: 統合テスト実装**
- **優先度**: High
- **推定時間**: 4時間
- **依存関係**: Task 6.3

**具体的なアクション**:
1. `server/route_test.go` 作成
2. テストデータベース設定
3. エンドツーエンドテスト実装

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/server/route_test.go`

**受け入れ基準**:
- [ ] 全ての統合テストが通る
- [ ] 正常系・異常系テストが実装されている

---

### **Task 7.2: OpenAPI仕様確認・API動作確認**
- **優先度**: High
- **推定時間**: 3時間
- **依存関係**: Task 7.1

**具体的なアクション**:
1. Swagger UIでのAPI仕様確認
2. 生成されたJSONスキーマの妥当性確認
3. Postman/curlでのAPI確認
4. Swagger UIからのAPI実行テスト

**受け入れ基準**:
- [ ] Swagger UIでAPI操作が可能
- [ ] 全てのエンドポイントが正常動作する
- [ ] レスポンス形式が仕様通りである

---

### **Task 7.3: パフォーマンステスト・セキュリティチェック**
- **優先度**: Medium
- **推定時間**: 3時間
- **依存関係**: Task 7.2

**具体的なアクション**:
1. 基本的な負荷テスト実行
2. メモリ使用量確認
3. SQLインジェクション対策確認
4. パスワードハッシュ化確認

**受け入れ基準**:
- [ ] パフォーマンスが許容範囲内
- [ ] セキュリティ要件を満たしている

---

## **Phase 8: ドキュメント・デプロイ準備**

### **Task 8.1: README.md作成**
- **優先度**: High
- **推定時間**: 2時間
- **依存関係**: Task 7.3

**具体的なアクション**:
1. `README.md` 完成版作成
2. 環境構築手順記載
3. API使用方法記載
4. OpenAPI/Swagger アクセス方法記載

**成果物**:
- `/Users/km/web-dev/_github/km-api-go/README.md`

**受け入れ基準**:
- [ ] 開発者が容易にプロジェクトを理解できる
- [ ] セットアップ手順が明確に記載されている

---

### **Task 8.2: デプロイメント準備**
- **優先度**: Low
- **推定時間**: 2時間
- **依存関係**: Task 8.1

**具体的なアクション**:
1. `Dockerfile` 作成（必要に応じて）
2. 本番環境用設定
3. ヘルスチェック エンドポイント追加

**成果物**:
- `Dockerfile`（必要に応じて）
- ヘルスチェックエンドポイント

**受け入れ基準**:
- [ ] デプロイメントが容易に実行できる
- [ ] ヘルスチェックが動作する

---

## **実行順序と優先度まとめ**

### **必須実行順序（High優先度）**:
1. **Task 1.1 → Task 1.2 → Task 1.3**（Phase 1）
2. **Task 2.1 → Task 2.2 → Task 2.3 → Task 2.4 → Task 2.5**（Phase 2）
3. **Task 3.1 → Task 3.2**（Phase 3）
4. **Task 4.1 → Task 4.2**（Phase 4）
5. **Task 5.1 → Task 5.2 → Task 5.3 → Task 5.4**（Phase 5）
6. **Task 6.1 → Task 6.2 → Task 6.3**（Phase 6）
7. **Task 7.1 → Task 7.2**（Phase 7）
8. **Task 8.1**（Phase 8）

### **並行実行可能なタスク**:
- **Task 2.6**（共通コンポーネント）はPhase 3と並行実行可能
- **Task 3.3**（テスト）、**Task 4.3**（テスト）、**Task 5.5**（テスト）は並行実行可能
- **Task 7.3**（パフォーマンステスト）と**Task 8.2**（デプロイ準備）は並行実行可能

### **総見積もり時間**: 
- **High優先度タスク**: 約45時間（約6日間）
- **Medium優先度タスク**: 約17時間（約2日間）
- **Low優先度タスク**: 約2時間

### **推奨実行スケジュール**:

#### **Week 1（6日間）**:
- Day 1: Task 1.1, 1.2, 1.3（基盤構築完了）
- Day 2: Task 2.1, 2.2, 2.3（ドメインモデル完了）
- Day 3: Task 2.4, 2.5, 2.6（インターフェース・マイグレーション完了）
- Day 4: Task 3.1, 3.2（リポジトリ実装完了）
- Day 5: Task 4.1, 4.2（ユースケース実装完了）
- Day 6: Task 5.1, 5.2（DTO・Userハンドラー完了）

#### **Week 2（4日間）**:
- Day 7: Task 5.3, 5.4（Companyハンドラー・OpenAPI完了）
- Day 8: Task 6.1, 6.2, 6.3（サーバー設定完了）
- Day 9: Task 7.1, 7.2（統合テスト・動作確認完了）
- Day 10: Task 8.1, 8.2（ドキュメント・デプロイ準備完了）

#### **テストフェーズ（2-3日間）**:
- 並行してTask 3.3, 4.3, 5.5, 7.3の実行

この詳細タスク分解により、systematic な開発進行が可能になり、各段階で成果物の品質を確保しながらAPI サーバーを構築できます。