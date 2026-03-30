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
