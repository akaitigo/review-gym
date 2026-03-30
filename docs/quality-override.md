# 品質チェックリスト追加分 — Webアプリ

> Layer-0（共通）+ Layer-1（言語別）のチェックリストに**追加**する項目のみ。

## Webアプリ固有の品質基準

### アクセシビリティ（a11y）
- [ ] WCAG 2.1 AA 準拠をチェック（axe-core / Lighthouse）
- [ ] キーボードのみで全機能が操作可能
- [ ] 適切な ARIA ラベルが設定されている
- [ ] カラーコントラスト比が基準を満たしている

### パフォーマンス
- [ ] Lighthouse Performance スコアが 90+
- [ ] Core Web Vitals（LCP, FID, CLS）が基準内
- [ ] バンドルサイズが適切（不要な依存がない）
- [ ] 画像の最適化（next/image, lazy loading 等）

### SEO・メタデータ（公開アプリの場合）
- [ ] ページタイトル・meta description が設定されている
- [ ] OGP タグが設定されている
- [ ] canonical URL が設定されている

### セキュリティ
- [ ] XSS 対策: ユーザー入力がエスケープされている
- [ ] CSRF 対策: フォーム送信にCSRFトークンがある
- [ ] CSP ヘッダが設定されている
- [ ] `rel="noopener noreferrer"` が外部リンクに付与されている

### E2Eテスト
- [ ] Playwright テストが主要ユーザーフローをカバー
- [ ] ログイン→主要操作→ログアウトのフローテスト
- [ ] モバイルビューポートでのテスト
- [ ] エラー画面（404, 500）のテスト

### ビジュアルリグレッションテスト（オプション）

> UI変更の視覚的な回帰を検出する。CSSの予期しない変更を防ぐ。
> 原則5参照: テスト生成とCI実行を分離する。

| ツール | 特徴 | 導入コスト |
|--------|------|-----------|
| **Chromatic** | Storybook統合。コンポーネント単位のスナップショット比較 | 低（Storybook使用時） |
| **Percy** (BrowserStack) | Playwrightと統合可能。クロスブラウザ比較 | 中 |
| **Argos** | OSS。GitHub Actions統合。セルフホスト可能 | 低 |

**選定基準:**
- Storybookを使っている → **Chromatic**
- Playwrightベースの既存E2Eがある → **Percy** or **Argos**
- セルフホスト・OSS重視 → **Argos**

**具体設定値:**
| 項目 | 設定値 | 説明 |
|------|--------|------|
| `visual_regression_threshold` | 0.1 | Argos-CI差分許容値（10%以下の差分は許容） |
| Playwright `slowMo` | 100ms | DOM変化追従用。アニメーション完了待ち |
| スクリーンショット比較 | pixel-by-pixel | 構造変化は即座に検出 |

**設定例（Playwright + Argos）:**
```bash
# CI設定に追加
npx @argos-ci/playwright test/e2e/visual/
npx @argos-ci/cli upload ./test-results/
```

**playwright.config.ts 追加設定:**
```typescript
// ビジュアルリグレッション用
use: {
  launchOptions: {
    slowMo: 100, // DOM変化追従
  },
},
expect: {
  toHaveScreenshot: {
    maxDiffPixelRatio: 0.1, // 差分許容値
  },
},
```

### 国際化（i18n）（対応する場合）
- [ ] ハードコードされた文字列がない
- [ ] RTL レイアウトの考慮（該当する場合）
- [ ] 日時フォーマットがロケール対応

### レスポンシブ
- [ ] モバイル / タブレット / デスクトップで表示確認
- [ ] ブレークポイントが一貫している
