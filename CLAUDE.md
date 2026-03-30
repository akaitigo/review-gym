# review-gym

OSSの過去PRを題材にコードレビュースキルを練習するプラットフォーム。セキュリティ・パフォーマンス問題の「目利き力」を鍛える。

## 技術スタック
- Go: REST API + スコアリングエンジン
- TypeScript/React: diff表示+コメントUI
- PostgreSQL + Redis

## ルール
- Go: golangci-lint + gofumpt
- TypeScript: `~/.claude/rules/typescript.md`

## コマンド
```
make check     # lint → test → build
make quality   # 品質ゲート
```

## 構造
```
backend/cmd/ internal/handler/ store/ scoring/ seed/ model/
frontend/src/ components/ hooks/ lib/ types/
```

## 環境変数
`.env.example` を参照
