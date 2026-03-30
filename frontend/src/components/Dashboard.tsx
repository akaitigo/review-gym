import { getUserAnalytics, getUserRecommendations } from "@/lib/api";
import type { AnalyticsData, Category, RecommendationItem, RecommendationsData } from "@/types/exercise";
import { CATEGORY_COLORS, CATEGORY_LABELS, DIFFICULTY_LABELS, TREND_ICONS, TREND_LABELS } from "@/types/exercise";
import { useCallback, useEffect, useState } from "react";
import {
	CartesianGrid,
	Legend,
	Line,
	LineChart,
	PolarAngleAxis,
	PolarGrid,
	PolarRadiusAxis,
	Radar,
	RadarChart,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from "recharts";

interface DashboardProps {
	userId: string;
	onSelectExercise: (id: string) => void;
	onBack: () => void;
}

export function Dashboard({ userId, onSelectExercise, onBack }: DashboardProps) {
	const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
	const [recommendations, setRecommendations] = useState<RecommendationsData | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		setLoading(true);
		setError(null);

		Promise.all([getUserAnalytics(userId), getUserRecommendations(userId)])
			.then(([analyticsData, recsData]) => {
				setAnalytics(analyticsData);
				setRecommendations(recsData);
			})
			.catch((err: unknown) => {
				setError(err instanceof Error ? err.message : "Failed to load dashboard data");
			})
			.finally(() => {
				setLoading(false);
			});
	}, [userId]);

	const handleExerciseClick = useCallback(
		(exerciseId: string) => {
			onSelectExercise(exerciseId);
		},
		[onSelectExercise],
	);

	if (loading) {
		return <p className="loading">Loading dashboard...</p>;
	}

	if (error !== null) {
		return (
			<div className="error-container">
				<p className="error">Error: {error}</p>
				<button type="button" onClick={onBack} className="btn-secondary">
					Back to exercises
				</button>
			</div>
		);
	}

	// Show message if not enough data.
	if (analytics?.message !== undefined) {
		return (
			<div className="dashboard">
				<div className="dashboard-header">
					<button type="button" onClick={onBack} className="btn-back" aria-label="Back to exercise list">
						&larr; Back
					</button>
					<h2>Dashboard</h2>
				</div>
				<div className="dashboard-message">
					<p>{analytics.message}</p>
					<p className="dashboard-progress">
						{String(analytics.total_exercises_completed)} / {String(analytics.min_exercises_required ?? 3)} exercises
						completed
					</p>
					<button type="button" onClick={onBack} className="btn-primary">
						Start Practicing
					</button>
				</div>
			</div>
		);
	}

	if (analytics === null) {
		return null;
	}

	// Prepare radar chart data.
	const radarData = analytics.categories.map((cat) => ({
		category: CATEGORY_LABELS[cat.category],
		score: cat.average_score,
		fullMark: 100,
	}));

	// Prepare line chart data.
	const lineData = analytics.score_history.map((point) => ({
		name: `#${String(point.attempt_index)}`,
		score: point.overall_score,
		date: point.date,
	}));

	return (
		<div className="dashboard">
			<div className="dashboard-header">
				<button type="button" onClick={onBack} className="btn-back" aria-label="Back to exercise list">
					&larr; Back
				</button>
				<h2>Dashboard</h2>
			</div>

			{/* Stats Overview */}
			<div className="dashboard-stats">
				<div className="stat-card">
					<span className="stat-value">{String(analytics.total_exercises_completed)}</span>
					<span className="stat-label">Exercises Completed</span>
				</div>
				<div className="stat-card">
					<span className="stat-value">{String(analytics.total_attempts)}</span>
					<span className="stat-label">Total Attempts</span>
				</div>
				<div className="stat-card">
					<span className="stat-value" style={{ color: getScoreColor(analytics.overall_average_score) }}>
						{String(analytics.overall_average_score)}
					</span>
					<span className="stat-label">Average Score</span>
				</div>
				<div className="stat-card">
					<span className="stat-value">{String(analytics.consecutive_days)}</span>
					<span className="stat-label">Day Streak</span>
				</div>
			</div>

			{/* Charts Section */}
			<div className="dashboard-charts">
				{/* Radar Chart */}
				<div className="chart-card">
					<h3>Category Breakdown</h3>
					<div className="chart-container">
						<ResponsiveContainer width="100%" height={300}>
							<RadarChart data={radarData}>
								<PolarGrid stroke="#30363d" />
								<PolarAngleAxis dataKey="category" tick={{ fill: "#8b949e", fontSize: 12 }} />
								<PolarRadiusAxis angle={90} domain={[0, 100]} tick={{ fill: "#8b949e", fontSize: 10 }} />
								<Radar name="Score" dataKey="score" stroke="#58a6ff" fill="#58a6ff" fillOpacity={0.3} />
							</RadarChart>
						</ResponsiveContainer>
					</div>

					{/* Weakness highlights */}
					{analytics.weakness_categories.length > 0 && (
						<div className="weakness-section">
							<h4>Areas for Improvement</h4>
							<div className="weakness-tags">
								{analytics.weakness_categories.map((cat) => (
									<span key={cat} className="weakness-tag" style={{ borderColor: CATEGORY_COLORS[cat] }}>
										{CATEGORY_LABELS[cat]}
									</span>
								))}
							</div>
						</div>
					)}
				</div>

				{/* Line Chart (only when 3+ attempts) */}
				{analytics.score_history.length >= 3 && (
					<div className="chart-card">
						<h3>Score Trend</h3>
						<div className="chart-container">
							<ResponsiveContainer width="100%" height={300}>
								<LineChart data={lineData}>
									<CartesianGrid strokeDasharray="3 3" stroke="#30363d" />
									<XAxis dataKey="name" tick={{ fill: "#8b949e", fontSize: 12 }} />
									<YAxis domain={[0, 100]} tick={{ fill: "#8b949e", fontSize: 12 }} />
									<Tooltip
										contentStyle={{
											backgroundColor: "#161b22",
											border: "1px solid #30363d",
											borderRadius: "6px",
											color: "#e6edf3",
										}}
									/>
									<Legend />
									<Line
										type="monotone"
										dataKey="score"
										stroke="#58a6ff"
										strokeWidth={2}
										dot={{ fill: "#58a6ff", r: 4 }}
										name="Overall Score"
									/>
								</LineChart>
							</ResponsiveContainer>
						</div>
					</div>
				)}
			</div>

			{/* Category Details */}
			<div className="dashboard-categories">
				<h3>Category Details</h3>
				<div className="category-detail-grid">
					{analytics.categories.map((cat) => (
						<CategoryDetailCard key={cat.category} category={cat} />
					))}
				</div>
			</div>

			{/* Recommendations */}
			{recommendations !== null && recommendations.recommendations.length > 0 && (
				<div className="dashboard-recommendations">
					<h3>Recommended Exercises</h3>
					<p className="recommendations-subtitle">Based on your weakness areas, try these exercises to improve:</p>
					<ul className="recommendation-list">
						{recommendations.recommendations.map((rec) => (
							<RecommendationCard key={rec.exercise.id} recommendation={rec} onClick={handleExerciseClick} />
						))}
					</ul>
				</div>
			)}
		</div>
	);
}

interface CategoryDetailCardProps {
	category: {
		category: Category;
		average_score: number;
		min_score: number;
		max_score: number;
		trend: string;
		is_weakness: boolean;
	};
}

function CategoryDetailCard({ category: cat }: CategoryDetailCardProps) {
	const trendIcon = TREND_ICONS[cat.trend] ?? "";
	const trendLabel = TREND_LABELS[cat.trend] ?? cat.trend;
	const color = CATEGORY_COLORS[cat.category];

	return (
		<div className={`category-detail-card ${cat.is_weakness ? "is-weakness" : ""}`} style={{ borderLeftColor: color }}>
			<div className="category-detail-header">
				<span className="category-detail-name">{CATEGORY_LABELS[cat.category]}</span>
				{cat.is_weakness && <span className="weakness-indicator">Weakness</span>}
			</div>
			<div className="category-detail-score" style={{ color }}>
				{String(cat.average_score)}%
			</div>
			<div className="category-detail-meta">
				<span className="category-detail-range">
					Range: {String(cat.min_score)} - {String(cat.max_score)}
				</span>
				<span className={`category-detail-trend trend-${cat.trend}`}>
					{trendIcon} {trendLabel}
				</span>
			</div>
		</div>
	);
}

interface RecommendationCardProps {
	recommendation: RecommendationItem;
	onClick: (id: string) => void;
}

function RecommendationCard({ recommendation: rec, onClick }: RecommendationCardProps) {
	return (
		<li className="recommendation-item">
			<button type="button" className="recommendation-card" onClick={() => onClick(rec.exercise.id)}>
				<div className="recommendation-header">
					<h4>{rec.exercise.title}</h4>
					<div className="exercise-meta">
						<span className={`badge difficulty-${rec.exercise.difficulty}`}>
							{DIFFICULTY_LABELS[rec.exercise.difficulty]}
						</span>
						<span className={`badge category-${rec.exercise.category}`}>{CATEGORY_LABELS[rec.exercise.category]}</span>
						{rec.previously_attempted && <span className="badge attempted">Retake</span>}
					</div>
				</div>
				<p className="recommendation-reason">{rec.recommended_reason}</p>
				<span className="recommendation-weakness-tag" style={{ borderColor: CATEGORY_COLORS[rec.target_weakness] }}>
					Target: {CATEGORY_LABELS[rec.target_weakness]}
				</span>
			</button>
		</li>
	);
}

function getScoreColor(score: number): string {
	if (score >= 80) return "#3fb950";
	if (score >= 60) return "#d29922";
	if (score >= 40) return "#e67e22";
	return "#f85149";
}
