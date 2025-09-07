# km-api-go
- Go (Echoフレームワーク) を使用して構築された、練習用のバックエンドAPIサーバーです。
- 生成AIに認識させやすくするために、敢えてクリーンアーキテクチャではない構成にしています。
- githubのspec kitを利用して作成しています。

## 期間
2025年9月7日-約3hを利用し、作成。

## ✨ 技術スタック
- Go
- Echo (Webフレームワーク)
- GORM (ORM)
- PostgreSQL (データベース)

## 🚀 セットアップと実行
### 前提条件
- Go (バージョン 1.21 以降)
- PostgreSQL データベースがローカル環境で起動していること

### 環境変数の設定
```bash
cp .env.example .env
```

### 依存関係のインストール
```bash
go mod tidy
```

### データベースのマイグレーション
```bash
go run migrations/migrate.go
```

> **Note:** 現在の実装では、このコマンドはSQLクエリを画面に表示するだけです。表示されたSQLを手動でデータベースで実行するか、`migrate.go` を修正してGORMの `Exec()` メソッドで実行するようにしてください。

### アプリケーションの起動
```bash
go run cmd/api/main.go
```
サーバーが正常に起動すると、デフォルトでは `localhost:8080` でリクエストを待ち受けます。

### データベースへの接続
ローカルで起動しているPostgreSQLデータベースには、`psql`コマンドを使用して接続できます。
`.env`ファイルで設定したユーザー名（例: `km`）を使用してください。

```bash
psql --host=localhost --port=5432 --username=km --dbname=km_api
```

##  APIエンドポイント例
- **ヘルスチェック:** `GET /api/v1/health`
- **ユーザー作成:** `POST /api/v1/users`
