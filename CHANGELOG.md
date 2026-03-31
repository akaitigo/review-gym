# Changelog

## v1.0.0 (2026-04-01)

### Initial Release

- chore: upgrade all frontend dependencies (#21)
- harden: fix quality gate false positive on test diff literals (#20)
- fix: parseDiff が複数ファイル diff を正しく分割するよう修正 (#19)
- fix: add CJK character support to content similarity tokenizer
- fix: add content similarity to scoring (prevent gaming with empty comments)
- docs: harvest retrospective
- docs: CHANGELOG + CLAUDE.md 50-line + PRD checked for v1.0.0
- feat(#9): scoring engine (#17)
- feat(#8): diff review UI with comment form (#16)
- feat(#3): プロジェクト基盤セットアップ（CI/CD、リンター、テストフレームワーク、DBスキーマ） (#11)
- Initial project scaffold from idea #693