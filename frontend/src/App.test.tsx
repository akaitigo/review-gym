import { describe, expect, it } from "vitest";
import { App } from "./App";

describe("App", () => {
	it("should be a function component", () => {
		expect(typeof App).toBe("function");
	});

	it("should return a valid React element", () => {
		const result = App();
		expect(result).toBeDefined();
		expect(result.type).toBe("main");
	});

	it("should contain the app title", () => {
		const result = App();
		const h1 = result.props.children[0];
		expect(h1.type).toBe("h1");
		expect(h1.props.children).toBe("Review Gym");
	});

	it("should contain the description", () => {
		const result = App();
		const p = result.props.children[1];
		expect(p.type).toBe("p");
		expect(p.props.children).toContain("コードレビュースキル");
	});
});
