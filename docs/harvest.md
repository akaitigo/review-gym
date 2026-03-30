# Harvest: review-gym

## 使えたもの
- [x] Makefile (make check / make quality)
- [x] lint設定 (golangci-lint + oxlint + biome)
- [x] CI YAML
- [x] CLAUDE.md (27行、50行以下達成)
- [x] ADR テンプレート (1件: scoring-algorithm)
- [x] 品質チェックリスト (make quality)
- [x] E2Eテスト雛形 (test/e2e/smoke.spec.ts)
- [x] Hooks（PostToolUse golangci-lint + oxlint）
- [x] lefthook
- [x] startup.sh

## 使えなかったもの（理由付き）
- 特になし

## テンプレート改善提案

| 対象ファイル | 変更内容 | 根拠 |
|-------------|---------|------|
| idea-eval SKILL.md | tech_stack_optimal の変更理由を idea-launch にも引き継ぐ | Kotlin→Goの変更理由がPRDに反映されなかった |
| seed パッケージパターン | シードデータをテンプレート化 | 演習問題12問+模範レビュー42件のパターンは他PJでも使える |

## メトリクス

| 項目 | 値 |
|------|-----|
| Issue (closed/total) | 5/5 |
| PR merged | 5 |
| テスト数 | 100+ |
| CI失敗数 | 0 |
| ADR数 | 1 |
| テンプレート実装率 | 95% |
| CLAUDE.md行数 | 27 |

## 次のPJへの申し送り
- スコアリングエンジンのgreedy matchingパターンはコードレビュー以外のドメインでも応用可能
- recharts でのレーダーチャート+折れ線グラフの組み合わせはダッシュボードの定番構成
