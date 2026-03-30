# review-gym — アーキテクチャ概要

## アーキテクチャ

モノレポ構成。フロントエンド（React SPA）とバックエンド（Go API）を分離。

```
[Browser] → [React SPA (Vite)] → [Go API Server] → [PostgreSQL]
                                                  → [Redis]
```

## 主要な設計判断

- ADR-001: (未作成) モノレポ vs マルチリポ
- ADR-002: (未作成) Go API vs Node.js API の選定理由
- ADR-003: (未作成) スコアリングアルゴリズムの設計

## 外部サービス連携

- **GitHub API**: PR データの取得元（将来の自動収集用）
- **GitHub OAuth**: ユーザー認証（Phase 1）

## データモデル概要

- `exercises`: 匿名化された PR 問題
- `review_comments`: ユーザーのレビューコメント
- `reference_reviews`: 模範レビュー
- `scores`: ユーザーのスコア履歴
- `user_profiles`: ユーザー情報・弱点カテゴリ
