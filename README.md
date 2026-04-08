# review-gym

実際のOSSプロジェクトの過去のPR（匿名化済み）を題材に、コードレビュースキルを練習するトレーニングプラットフォーム。

## 特徴

- 匿名化されたOSS PRの差分表示とレビューコメント入力
- 模範レビューとの比較によるスコアリング（精度・網羅性）
- 弱点カテゴリの分析と重点練習PRの推薦

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | TypeScript, React, Vite |
| バックエンド | Go |
| データベース | PostgreSQL |
| キャッシュ | Redis |
| E2Eテスト | Playwright |

## アーキテクチャ

```
[Browser] → [React SPA (Vite)] → [Go API Server] → [PostgreSQL]
                                                  → [Redis]
```

フロントエンドはReact SPAで差分表示とレビューコメントUIを提供し、バックエンドのGo APIがスコアリングエンジンとデータ管理を担当します。

## セットアップ

### 前提条件

- Node.js >= 22
- Go >= 1.23

> **Note**: デフォルトはインメモリストアで動作します（PostgreSQL・Redis 不要）。永続化が必要な場合は下記「PostgreSQL + Redis で起動」を参照してください。

### クイックスタート（インメモリモード）

```bash
# リポジトリのクローン
git clone git@github.com:akaitigo/review-gym.git
cd review-gym

# フロントエンド
cd frontend
npm install
npm run dev

# バックエンド（別ターミナル）
cd backend
go run ./cmd/review-gym

# E2Eテスト
npx playwright test
```

### PostgreSQL + Redis で起動

```bash
# .env を作成
cp .env.example .env
# DB_PASSWORD を設定

# Docker Compose で PostgreSQL + Redis を起動
docker compose up -d

# マイグレーション実行
cd backend && make migrate

# シードデータ投入
cd backend && make seed

# バックエンドを PostgreSQL モードで起動
STORE_TYPE=postgres \
DATABASE_URL=postgresql://review_gym:your_password@localhost:5432/review_gym?sslmode=disable \
REDIS_URL=redis://localhost:6379 \
go run ./cmd/review-gym
```

### 環境変数

| 変数 | 説明 | デフォルト |
|------|------|-----------|
| `STORE_TYPE` | ストア種別 (`memory` or `postgres`) | `memory` |
| `DATABASE_URL` | PostgreSQL接続文字列 | - |
| `REDIS_URL` | Redis接続文字列（省略時はキャッシュなし） | - |
| `API_PORT` | APIサーバーポート | `8080` |
| `CORS_ORIGINS` | 許可するCORSオリジン（カンマ区切り） | `http://localhost:3000,http://localhost:5173` |

## 開発

```bash
# 全体チェック
make check

# 品質チェック
make quality
```

## ライセンス

MIT
