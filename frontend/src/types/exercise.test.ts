import { describe, expect, it } from "vitest";
import {
	ALL_CATEGORIES,
	ALL_DIFFICULTIES,
	CATEGORY_COLORS,
	CATEGORY_LABELS,
	DIFFICULTY_LABELS,
	SEVERITY_COLORS,
	SEVERITY_LABELS,
	TREND_ICONS,
	TREND_LABELS,
} from "./exercise";
import type {
	AnalyticsData,
	CategoryAnalytics,
	CategoryScore,
	FalsePositive,
	MissedReview,
	RecommendationItem,
	RecommendationsData,
	ScoreHistoryPoint,
	ScoreResult,
	ScoringMatch,
} from "./exercise";

describe("exercise types", () => {
	it("ALL_CATEGORIES should have 5 categories", () => {
		expect(ALL_CATEGORIES).toHaveLength(5);
	});

	it("ALL_CATEGORIES should contain all expected categories", () => {
		expect(ALL_CATEGORIES).toContain("security");
		expect(ALL_CATEGORIES).toContain("performance");
		expect(ALL_CATEGORIES).toContain("design");
		expect(ALL_CATEGORIES).toContain("readability");
		expect(ALL_CATEGORIES).toContain("error-handling");
	});

	it("ALL_DIFFICULTIES should have 3 levels", () => {
		expect(ALL_DIFFICULTIES).toHaveLength(3);
	});

	it("ALL_DIFFICULTIES should contain all expected levels", () => {
		expect(ALL_DIFFICULTIES).toContain("beginner");
		expect(ALL_DIFFICULTIES).toContain("intermediate");
		expect(ALL_DIFFICULTIES).toContain("advanced");
	});

	it("CATEGORY_LABELS should have labels for all categories", () => {
		for (const cat of ALL_CATEGORIES) {
			expect(CATEGORY_LABELS[cat]).toBeDefined();
			expect(CATEGORY_LABELS[cat].length).toBeGreaterThan(0);
		}
	});

	it("DIFFICULTY_LABELS should have labels for all difficulties", () => {
		for (const diff of ALL_DIFFICULTIES) {
			expect(DIFFICULTY_LABELS[diff]).toBeDefined();
			expect(DIFFICULTY_LABELS[diff].length).toBeGreaterThan(0);
		}
	});

	it("SEVERITY_LABELS should have labels for all severity levels", () => {
		const severities = ["critical", "major", "minor", "info"] as const;
		for (const sev of severities) {
			expect(SEVERITY_LABELS[sev]).toBeDefined();
			expect(SEVERITY_LABELS[sev].length).toBeGreaterThan(0);
		}
	});

	it("SEVERITY_COLORS should have colors for all severity levels", () => {
		const severities = ["critical", "major", "minor", "info"] as const;
		for (const sev of severities) {
			expect(SEVERITY_COLORS[sev]).toBeDefined();
			expect(SEVERITY_COLORS[sev]).toMatch(/^#[0-9a-f]{6}$/);
		}
	});

	it("ScoreResult type should be constructable", () => {
		const score: ScoreResult = {
			id: "1",
			exercise_id: "ex-1",
			user_id: "user-1",
			precision_score: 100,
			recall_score: 75,
			overall_score: 85,
			category_scores: [],
			attempt_number: 1,
			matches: [],
			missed_reviews: [],
			false_positives: [],
		};
		expect(score.overall_score).toBe(85);
	});

	it("CategoryScore type should be constructable", () => {
		const cs: CategoryScore = {
			category: "security",
			score: 100,
			max_points: 3,
			earned: 3,
		};
		expect(cs.score).toBe(100);
	});

	it("MissedReview type should be constructable", () => {
		const missed: MissedReview = {
			file_path: "main.go",
			line_number: 10,
			content: "SQL injection",
			category: "security",
			severity: "critical",
			explanation: "Use parameterized queries",
		};
		expect(missed.severity).toBe("critical");
	});

	it("FalsePositive type should be constructable", () => {
		const fp: FalsePositive = {
			file_path: "main.go",
			line_number: 5,
			content: "not a real issue",
			category: "design",
		};
		expect(fp.category).toBe("design");
	});

	it("ScoringMatch type should be constructable", () => {
		const match: ScoringMatch = {
			user_comment: {
				id: "c-1",
				exercise_id: "ex-1",
				user_id: "user-1",
				file_path: "main.go",
				line_number: 10,
				content: "Found it",
				category: "security",
				created_at: "2026-01-01T00:00:00Z",
				updated_at: "2026-01-01T00:00:00Z",
			},
			reference_review: {
				id: "ref-1",
				exercise_id: "ex-1",
				file_path: "main.go",
				line_number: 10,
				content: "SQL injection",
				category: "security",
				severity: "critical",
				explanation: "Use parameterized queries",
			},
			line_delta: 0,
		};
		expect(match.line_delta).toBe(0);
	});

	it("CATEGORY_COLORS should have colors for all categories", () => {
		for (const cat of ALL_CATEGORIES) {
			expect(CATEGORY_COLORS[cat]).toBeDefined();
			expect(CATEGORY_COLORS[cat]).toMatch(/^#[0-9a-f]{6}$/);
		}
	});

	it("TREND_LABELS should have labels for all trend types", () => {
		const trends = ["improving", "stagnating", "declining"];
		for (const trend of trends) {
			expect(TREND_LABELS[trend]).toBeDefined();
			expect(TREND_LABELS[trend].length).toBeGreaterThan(0);
		}
	});

	it("TREND_ICONS should have icons for all trend types", () => {
		const trends = ["improving", "stagnating", "declining"];
		for (const trend of trends) {
			expect(TREND_ICONS[trend]).toBeDefined();
		}
	});

	it("CategoryAnalytics type should be constructable", () => {
		const ca: CategoryAnalytics = {
			category: "security",
			average_score: 65.5,
			min_score: 40,
			max_score: 80,
			trend: "improving",
			is_weakness: false,
		};
		expect(ca.trend).toBe("improving");
	});

	it("ScoreHistoryPoint type should be constructable", () => {
		const point: ScoreHistoryPoint = {
			date: "2026-03-28",
			overall_score: 72.5,
			attempt_index: 3,
		};
		expect(point.attempt_index).toBe(3);
	});

	it("AnalyticsData type should be constructable", () => {
		const data: AnalyticsData = {
			user_id: "user-1",
			total_exercises_completed: 5,
			total_attempts: 8,
			overall_average_score: 65.0,
			categories: [],
			weakness_categories: ["security"],
			score_history: [],
			consecutive_days: 3,
		};
		expect(data.consecutive_days).toBe(3);
	});

	it("RecommendationItem type should be constructable", () => {
		const rec: RecommendationItem = {
			exercise: {
				id: "ex-1",
				title: "Test Exercise",
				description: "desc",
				difficulty: "beginner",
				category: "security",
				category_tags: ["security"],
				language: "Go",
				file_paths: [],
				is_published: true,
			},
			recommended_reason: "Targets security weakness",
			target_weakness: "security",
			previously_attempted: false,
		};
		expect(rec.previously_attempted).toBe(false);
	});

	it("RecommendationsData type should be constructable", () => {
		const data: RecommendationsData = {
			user_id: "user-1",
			weakness_categories: ["security"],
			recommendations: [],
		};
		expect(data.weakness_categories).toHaveLength(1);
	});
});
