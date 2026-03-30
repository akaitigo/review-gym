import type { ScoreResult as ScoreResultType } from "@/types/exercise";
import { CATEGORY_LABELS, SEVERITY_COLORS, SEVERITY_LABELS } from "@/types/exercise";

interface ScoreResultProps {
	result: ScoreResultType;
	onRetry: () => void;
	onBack: () => void;
}

export function ScoreResult({ result, onRetry, onBack }: ScoreResultProps) {
	const overallColor = getScoreColor(result.overall_score);

	return (
		<div className="score-result">
			<h2>Review Score</h2>

			<div className="score-overview">
				<div className="score-overall" style={{ borderColor: overallColor }}>
					<span className="score-value" style={{ color: overallColor }}>
						{String(result.overall_score)}
					</span>
					<span className="score-label">Overall</span>
				</div>

				<div className="score-breakdown">
					<div className="score-metric">
						<span className="metric-label">Precision</span>
						<div className="score-bar-container">
							<div
								className="score-bar"
								style={{
									width: `${String(result.precision_score)}%`,
									backgroundColor: getScoreColor(result.precision_score),
								}}
							/>
						</div>
						<span className="metric-value">{String(result.precision_score)}%</span>
					</div>
					<div className="score-metric">
						<span className="metric-label">Recall</span>
						<div className="score-bar-container">
							<div
								className="score-bar"
								style={{
									width: `${String(result.recall_score)}%`,
									backgroundColor: getScoreColor(result.recall_score),
								}}
							/>
						</div>
						<span className="metric-value">{String(result.recall_score)}%</span>
					</div>
				</div>
			</div>

			<div className="score-attempt">Attempt #{String(result.attempt_number)}</div>

			<div className="category-scores">
				<h3>Category Breakdown</h3>
				<div className="category-grid">
					{result.category_scores.map((cs) => (
						<div key={cs.category} className="category-score-item">
							<span className="category-name">{CATEGORY_LABELS[cs.category]}</span>
							<div className="score-bar-container">
								<div
									className="score-bar"
									style={{
										width: `${String(cs.score)}%`,
										backgroundColor: getScoreColor(cs.score),
									}}
								/>
							</div>
							<span className="category-value">{String(cs.score)}%</span>
						</div>
					))}
				</div>
			</div>

			{result.matches.length > 0 && (
				<div className="score-section">
					<h3>Matched Findings ({String(result.matches.length)})</h3>
					<ul className="finding-list matched-list">
						{result.matches.map((match, idx) => (
							<li key={`match-${String(idx)}`} className="finding-item matched">
								<div className="finding-header">
									<span className="finding-badge matched-badge">Matched</span>
									<span className="finding-location">
										{match.reference_review.file_path}:{String(match.reference_review.line_number)}
									</span>
									<span className={`finding-category cat-${match.reference_review.category}`}>
										{CATEGORY_LABELS[match.reference_review.category]}
									</span>
								</div>
								<p className="finding-content">{match.reference_review.content}</p>
								{match.line_delta > 0 && (
									<p className="finding-note">
										Your comment was {String(match.line_delta)} line{match.line_delta > 1 ? "s" : ""} away
									</p>
								)}
							</li>
						))}
					</ul>
				</div>
			)}

			{result.missed_reviews.length > 0 && (
				<div className="score-section">
					<h3>Missed Review Points ({String(result.missed_reviews.length)})</h3>
					<ul className="finding-list missed-list">
						{result.missed_reviews.map((missed, idx) => (
							<li key={`missed-${String(idx)}`} className="finding-item missed">
								<div className="finding-header">
									<span
										className="finding-badge severity-badge"
										style={{ backgroundColor: SEVERITY_COLORS[missed.severity] }}
									>
										{SEVERITY_LABELS[missed.severity]}
									</span>
									<span className="finding-location">
										{missed.file_path}:{String(missed.line_number)}
									</span>
									<span className={`finding-category cat-${missed.category}`}>{CATEGORY_LABELS[missed.category]}</span>
								</div>
								<p className="finding-content">{missed.content}</p>
								<details className="finding-explanation">
									<summary>Explanation</summary>
									<p>{missed.explanation}</p>
								</details>
							</li>
						))}
					</ul>
				</div>
			)}

			{result.false_positives.length > 0 && (
				<div className="score-section">
					<h3>Unmatched Comments ({String(result.false_positives.length)})</h3>
					<ul className="finding-list fp-list">
						{result.false_positives.map((fp, idx) => (
							<li key={`fp-${String(idx)}`} className="finding-item false-positive">
								<div className="finding-header">
									<span className="finding-badge fp-badge">Unmatched</span>
									<span className="finding-location">
										{fp.file_path}:{String(fp.line_number)}
									</span>
									<span className={`finding-category cat-${fp.category}`}>{CATEGORY_LABELS[fp.category]}</span>
								</div>
								<p className="finding-content">{fp.content}</p>
							</li>
						))}
					</ul>
				</div>
			)}

			<div className="score-actions">
				<button type="button" className="btn-secondary" onClick={onBack}>
					Back to Exercises
				</button>
				<button type="button" className="btn-primary" onClick={onRetry}>
					Try Again
				</button>
			</div>
		</div>
	);
}

function getScoreColor(score: number): string {
	if (score >= 80) return "#2ecc71";
	if (score >= 60) return "#f39c12";
	if (score >= 40) return "#e67e22";
	return "#e74c3c";
}
