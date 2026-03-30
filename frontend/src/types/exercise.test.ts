import { describe, expect, it } from "vitest";
import { ALL_CATEGORIES, ALL_DIFFICULTIES, CATEGORY_LABELS, DIFFICULTY_LABELS } from "./exercise";

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
});
