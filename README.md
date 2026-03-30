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

## セットアップ

### 前提条件

- Node.js >= 22
- Go >= 1.23
- PostgreSQL
- Redis

### クイックスタート

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

## 開発

```bash
# 全体チェック
make check

# 品質チェック
make quality
```

## ライセンス

MIT
