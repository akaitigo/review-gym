import { describe, expect, it } from "vitest";
import { ApiError } from "./api";

describe("ApiError", () => {
	it("should have the correct name", () => {
		const err = new ApiError(404, "Not found");
		expect(err.name).toBe("ApiError");
	});

	it("should have the correct status", () => {
		const err = new ApiError(400, "Bad request");
		expect(err.status).toBe(400);
	});

	it("should have the correct message", () => {
		const err = new ApiError(500, "Internal server error");
		expect(err.message).toBe("Internal server error");
	});

	it("should be an instance of Error", () => {
		const err = new ApiError(404, "Not found");
		expect(err).toBeInstanceOf(Error);
	});
});
