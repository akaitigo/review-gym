import { describe, expect, it } from "vitest";
import { Dashboard } from "./Dashboard";

describe("Dashboard", () => {
	it("should be a function component", () => {
		expect(typeof Dashboard).toBe("function");
	});
});
