# ADR-003: Scoring Algorithm Design

## Status

Accepted

## Context

review-gym needs to evaluate how well a user's code review matches the curated reference reviews for each exercise. The scoring system must quantify both **precision** (how many of the user's comments are valid) and **recall** (how many reference review points the user found).

Key requirements:
- Users submit review comments with file path, line number, content, and category
- Reference reviews have file path, line number, content, category, and severity
- Security-critical findings should be weighted higher than informational ones
- Scoring must be deterministic and reproducible

## Decision

### Matching Algorithm

User comments are matched to reference reviews using a **proximity + category** heuristic:

1. **Same file path**: user comment and reference review must target the same file
2. **Line proximity**: line numbers must be within +/- 3 lines of each other
3. **Category match**: category must be identical for a match to count
4. **One-to-one**: each reference review can be matched to at most one user comment (greedy, closest line first)

### Scoring Formula

**Precision** = (matched user comments / total user comments) * 100

**Recall** = (weighted matched references / weighted total references) * 100

**Overall** = (Precision * 0.4 + Recall * 0.6)

Recall is weighted higher because finding real issues is more important than avoiding false positives.

### Severity Weights

| Severity | Weight |
|----------|--------|
| critical | 3.0 |
| major | 2.0 |
| minor | 1.0 |
| info | 0.5 |

### Category-Level Scoring

For each category, a sub-score (0-100) is computed using the same recall formula but scoped to that category's reference reviews only.

## Consequences

### Positive
- Simple, explainable algorithm that users can understand
- Deterministic: same inputs always produce the same score
- Severity weighting incentivizes finding critical issues first
- Category sub-scores enable weakness analysis (Issue #10)

### Negative
- Line-proximity matching may miss semantically correct comments placed on different lines
- No NLP/semantic similarity -- comments must be near the reference line to match
- Greedy matching may not find optimal assignment in edge cases

### Future Improvements
- Semantic matching using NLP to compare comment content
- Optimal assignment using Hungarian algorithm instead of greedy matching
- User-adjustable scoring parameters
