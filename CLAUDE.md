# review-gym

## 概要

実際のOSSプロジェクトの過去のPR（匿名化済み）を題材に、コードレビュースキルを練習するトレーニングプラットフォーム。セキュリティ脆弱性・パフォーマンス問題・設計上の懸念等を見つけ出す「目利き力」を鍛える。

## 技術スタック

- Go（API サーバー）
- TypeScript / React（フロントエンド）
- PostgreSQL（データベース）
- Redis（キャッシュ）
- Vite（ビルドツール）
- Docker Compose（開発環境）

## コーディングルール

- TypeScript: `~/.claude/rules/typescript.md` のルールに従うこと
  - `any` 型の使用禁止。`unknown` + 型ガードまたはジェネリクスを使う
  - `as` 型アサーションは最小限に。型推論またはジェネリクスで解決
- Go: 標準の Go スタイルガイドに従うこと
  - `golangci-lint` でリント
  - `gofumpt` でフォーマット

## ビルド & テスト

```bash
# フロントエンド
cd frontend && npm run check     # format + lint + typecheck + test + build

# バックエンド
cd backend && make check         # format + tidy + lint + test + build

# E2E テスト
npx playwright test

# 全体チェック
make check                       # frontend + backend + E2E
```

## ディレクトリ構造

```
review-gym/
├── frontend/          # TypeScript/React SPA
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── lib/
│   │   └── types/
│   ├── public/
│   ├── package.json
│   ├── tsconfig.json
│   └── biome.json
├── backend/           # Go API サーバー
│   ├── cmd/review-gym/
│   ├── internal/
│   ├── pkg/
│   ├── go.mod
│   └── Makefile
├── test/e2e/          # Playwright E2E テスト
├── docs/              # ドキュメント
├── .claude/           # Claude Code 設定
├── .github/           # GitHub Actions, テンプレート
├── PRD.md
├── CLAUDE.md
├── Makefile
└── docker-compose.yml
```

## 環境変数

```bash
# PostgreSQL
DATABASE_URL=postgresql://localhost:5432/review_gym
DB_PASSWORD=${DB_PASSWORD:review_gym_dev}

# Redis
REDIS_URL=redis://localhost:6379

# GitHub OAuth（将来）
GITHUB_CLIENT_ID=${GITHUB_CLIENT_ID:}
GITHUB_CLIENT_SECRET=${GITHUB_CLIENT_SECRET:}

# Server
API_PORT=8080
FRONTEND_PORT=3000
```
