import { describe, expect, it } from "vitest";
import { getMaxNewLineNumber, parseDiff } from "./diffParser";

const sampleDiff = `--- a/internal/handler/user.go
+++ b/internal/handler/user.go
@@ -0,0 +1,10 @@
+package handler
+
+import (
+	"database/sql"
+	"fmt"
+	"net/http"
+)
+
+func main() {
+}`;

describe("parseDiff", () => {
	it("should parse a simple unified diff", () => {
		const files = parseDiff(sampleDiff);
		expect(files).toHaveLength(1);

		const file = files[0];
		expect(file).toBeDefined();
		expect(file?.oldPath).toBe("internal/handler/user.go");
		expect(file?.newPath).toBe("internal/handler/user.go");
		expect(file?.hunks).toHaveLength(1);
	});

	it("should correctly identify addition lines", () => {
		const files = parseDiff(sampleDiff);
		const hunk = files[0]?.hunks[0];
		expect(hunk).toBeDefined();

		const additions = hunk?.lines.filter((l) => l.type === "addition") ?? [];
		expect(additions.length).toBe(10);
	});

	it("should assign correct line numbers to additions", () => {
		const files = parseDiff(sampleDiff);
		const hunk = files[0]?.hunks[0];
		expect(hunk).toBeDefined();

		const firstLine = hunk?.lines[0];
		expect(firstLine?.newLineNumber).toBe(1);
		expect(firstLine?.content).toBe("package handler");

		const lastLine = hunk?.lines[9];
		expect(lastLine?.newLineNumber).toBe(10);
	});

	it("should handle empty diff", () => {
		const files = parseDiff("");
		expect(files).toHaveLength(0);
	});

	it("should handle diff with context lines", () => {
		const diff = `--- a/test.go
+++ b/test.go
@@ -1,3 +1,4 @@
 package main

+import "fmt"
 func main() {}`;

		const files = parseDiff(diff);
		expect(files).toHaveLength(1);

		const hunk = files[0]?.hunks[0];
		expect(hunk).toBeDefined();

		const contextLines = hunk?.lines.filter((l) => l.type === "context") ?? [];
		expect(contextLines.length).toBe(3);

		const additionLines = hunk?.lines.filter((l) => l.type === "addition") ?? [];
		expect(additionLines.length).toBe(1);
		expect(additionLines[0]?.newLineNumber).toBe(3);
	});

	it("should handle diff with deletions", () => {
		const diff = `--- a/test.go
+++ b/test.go
@@ -1,3 +1,2 @@
 package main
-import "os"
 func main() {}`;

		const files = parseDiff(diff);
		const hunk = files[0]?.hunks[0];
		expect(hunk).toBeDefined();

		const deletionLines = hunk?.lines.filter((l) => l.type === "deletion") ?? [];
		expect(deletionLines.length).toBe(1);
		expect(deletionLines[0]?.oldLineNumber).toBe(2);
		expect(deletionLines[0]?.newLineNumber).toBeNull();
	});
});

describe("getMaxNewLineNumber", () => {
	it("should return the max new line number", () => {
		expect(getMaxNewLineNumber(sampleDiff)).toBe(10);
	});

	it("should return 0 for empty diff", () => {
		expect(getMaxNewLineNumber("")).toBe(0);
	});

	it("should handle multiple hunks", () => {
		const diff = `--- a/test.go
+++ b/test.go
@@ -1,2 +1,2 @@
 package main
+import "fmt"
@@ -10,2 +11,3 @@
 func main() {}
+func helper() {}
+func util() {}`;

		expect(getMaxNewLineNumber(diff)).toBe(13);
	});
});
