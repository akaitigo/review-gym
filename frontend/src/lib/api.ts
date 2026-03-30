import type {
	AnalyticsData,
	CreateReviewRequest,
	Exercise,
	ExerciseListItem,
	RecommendationsData,
	ReviewComment,
	ScoreHistoryEntry,
	ScoreResult,
} from "@/types/exercise";

const API_BASE = "/api";

/** API error with status code and message. */
export class ApiError extends Error {
	constructor(
		public readonly status: number,
		message: string,
	) {
		super(message);
		this.name = "ApiError";
	}
}

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
	const res = await fetch(url, init);
	if (!res.ok) {
		const body = await res.json().catch(() => ({ error: res.statusText }));
		const message = typeof body === "object" && body !== null && "error" in body ? String(body.error) : res.statusText;
		throw new ApiError(res.status, message);
	}
	return res.json() as Promise<T>;
}

/** Fetch exercise list with optional filters. */
export async function listExercises(params?: {
	category?: string | undefined;
	difficulty?: string | undefined;
}): Promise<ExerciseListItem[]> {
	const url = new URL(`${API_BASE}/exercises`, window.location.origin);
	if (params?.category) {
		url.searchParams.set("category", params.category);
	}
	if (params?.difficulty) {
		url.searchParams.set("difficulty", params.difficulty);
	}
	return fetchJSON<ExerciseListItem[]>(url.toString());
}

/** Fetch a single exercise by ID (includes diff_content). */
export async function getExercise(id: string): Promise<Exercise> {
	return fetchJSON<Exercise>(`${API_BASE}/exercises/${encodeURIComponent(id)}`);
}

/** Create a review comment on an exercise. */
export async function createReview(exerciseId: string, review: CreateReviewRequest): Promise<ReviewComment> {
	return fetchJSON<ReviewComment>(`${API_BASE}/exercises/${encodeURIComponent(exerciseId)}/reviews`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(review),
	});
}

/** Fetch review comments for an exercise (current user). */
export async function listReviews(exerciseId: string): Promise<ReviewComment[]> {
	return fetchJSON<ReviewComment[]>(`${API_BASE}/exercises/${encodeURIComponent(exerciseId)}/reviews`);
}

/** Score the current user's review against reference reviews. */
export async function scoreExercise(exerciseId: string): Promise<ScoreResult> {
	return fetchJSON<ScoreResult>(`${API_BASE}/exercises/${encodeURIComponent(exerciseId)}/score`, {
		method: "POST",
	});
}

/** Fetch score history for an exercise (current user). */
export async function listScores(exerciseId: string): Promise<ScoreHistoryEntry[]> {
	return fetchJSON<ScoreHistoryEntry[]>(`${API_BASE}/exercises/${encodeURIComponent(exerciseId)}/scores`);
}

/** Fetch user analytics (category scores, weaknesses, trends). */
export async function getUserAnalytics(userId: string): Promise<AnalyticsData> {
	return fetchJSON<AnalyticsData>(`${API_BASE}/users/${encodeURIComponent(userId)}/analytics`);
}

/** Fetch recommended exercises based on user weaknesses. */
export async function getUserRecommendations(userId: string): Promise<RecommendationsData> {
	return fetchJSON<RecommendationsData>(`${API_BASE}/users/${encodeURIComponent(userId)}/recommendations`);
}
