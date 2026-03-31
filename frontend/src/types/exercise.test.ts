import { describe, expect, it } from "vitest";
import type { CategoryScore, FalsePositive, MissedReview, ScoreResult, ScoringMatch } from "./exercise";
import {
	ALL_CATEGORIES,
	ALL_DIFFICULTIES,
	CATEGORY_LABELS,
	DIFFICULTY_LABELS,
	SEVERITY_COLORS,
	SEVERITY_LABELS,
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
});
