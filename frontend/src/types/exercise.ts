/** Review comment categories matching the backend model. */
export type Category = "security" | "performance" | "design" | "readability" | "error-handling";

/** Exercise difficulty levels. */
export type Difficulty = "beginner" | "intermediate" | "advanced";

/** Exercise list item (without diff_content). */
export interface ExerciseListItem {
	id: string;
	title: string;
	description: string;
	difficulty: Difficulty;
	category: Category;
	category_tags: Category[];
	language: string;
	file_paths: string[];
	is_published: boolean;
}

/** Full exercise detail (with diff_content). */
export interface Exercise extends ExerciseListItem {
	diff_content: string;
	metadata: Record<string, unknown>;
	created_at: string;
	updated_at: string;
}

/** Review comment submitted by a user. */
export interface ReviewComment {
	id: string;
	exercise_id: string;
	user_id: string;
	file_path: string;
	line_number: number;
	content: string;
	category: Category;
	created_at: string;
	updated_at: string;
}

/** Request body for creating a review comment. */
export interface CreateReviewRequest {
	file_path: string;
	line_number: number;
	content: string;
	category: string;
}

/** All available categories for display. */
export const ALL_CATEGORIES: readonly Category[] = [
	"security",
	"performance",
	"design",
	"readability",
	"error-handling",
] as const;

/** All available difficulties for display. */
export const ALL_DIFFICULTIES: readonly Difficulty[] = ["beginner", "intermediate", "advanced"] as const;

/** Human-readable labels for categories. */
export const CATEGORY_LABELS: Record<Category, string> = {
	security: "Security",
	performance: "Performance",
	design: "Design",
	readability: "Readability",
	"error-handling": "Error Handling",
};

/** Human-readable labels for difficulties. */
export const DIFFICULTY_LABELS: Record<Difficulty, string> = {
	beginner: "Beginner",
	intermediate: "Intermediate",
	advanced: "Advanced",
};

/** Severity levels matching the backend model. */
export type Severity = "critical" | "major" | "minor" | "info";

/** Category score breakdown. */
export interface CategoryScore {
	category: Category;
	score: number;
	max_points: number;
	earned: number;
}

/** A matched pair of user comment and reference review. */
export interface ScoringMatch {
	user_comment: ReviewComment;
	reference_review: ReferenceReview;
	line_delta: number;
}

/** A reference review point. */
export interface ReferenceReview {
	id: string;
	exercise_id: string;
	file_path: string;
	line_number: number;
	content: string;
	category: Category;
	severity: Severity;
	explanation: string;
}

/** A missed reference review in scoring results. */
export interface MissedReview {
	file_path: string;
	line_number: number;
	content: string;
	category: Category;
	severity: Severity;
	explanation: string;
}

/** A false positive user comment in scoring results. */
export interface FalsePositive {
	file_path: string;
	line_number: number;
	content: string;
	category: Category;
}

/** Full scoring result from the API. */
export interface ScoreResult {
	id: string;
	exercise_id: string;
	user_id: string;
	precision_score: number;
	recall_score: number;
	overall_score: number;
	category_scores: CategoryScore[];
	attempt_number: number;
	matches: ScoringMatch[];
	missed_reviews: MissedReview[];
	false_positives: FalsePositive[];
}

/** Score history entry. */
export interface ScoreHistoryEntry {
	id: string;
	user_id: string;
	exercise_id: string;
	precision_score: number;
	recall_score: number;
	overall_score: number;
	category_scores: string;
	attempt_number: number;
	created_at: string;
}

/** Human-readable labels for severity. */
export const SEVERITY_LABELS: Record<Severity, string> = {
	critical: "Critical",
	major: "Major",
	minor: "Minor",
	info: "Info",
};

/** Color mapping for severity levels. */
export const SEVERITY_COLORS: Record<Severity, string> = {
	critical: "#e74c3c",
	major: "#e67e22",
	minor: "#f1c40f",
	info: "#3498db",
};

/** Category analytics data from the analytics endpoint. */
export interface CategoryAnalytics {
	category: Category;
	average_score: number;
	min_score: number;
	max_score: number;
	trend: "improving" | "stagnating" | "declining";
	is_weakness: boolean;
}

/** Score history point for the trend chart. */
export interface ScoreHistoryPoint {
	date: string;
	overall_score: number;
	attempt_index: number;
}

/** Analytics response from GET /api/users/:id/analytics. */
export interface AnalyticsData {
	user_id: string;
	total_exercises_completed: number;
	total_attempts: number;
	overall_average_score: number;
	categories: CategoryAnalytics[];
	weakness_categories: Category[];
	score_history: ScoreHistoryPoint[];
	consecutive_days: number;
	message?: string;
	min_exercises_required?: number;
}

/** Recommendation item from the recommendations endpoint. */
export interface RecommendationItem {
	exercise: ExerciseListItem;
	recommended_reason: string;
	target_weakness: Category;
	previously_attempted: boolean;
}

/** Recommendations response from GET /api/users/:id/recommendations. */
export interface RecommendationsData {
	user_id: string;
	weakness_categories: Category[];
	recommendations: RecommendationItem[];
	message?: string;
	min_exercises_required?: number;
}

/** Color mapping for categories (for charts). */
export const CATEGORY_COLORS: Record<Category, string> = {
	security: "#f85149",
	performance: "#d29922",
	design: "#58a6ff",
	readability: "#3fb950",
	"error-handling": "#bc8cff",
};

/** Trend labels for display. */
export const TREND_LABELS: Record<string, string> = {
	improving: "Improving",
	stagnating: "Stable",
	declining: "Declining",
};

/** Trend icons for display. */
export const TREND_ICONS: Record<string, string> = {
	improving: "\u2191",
	stagnating: "\u2192",
	declining: "\u2193",
};
