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
