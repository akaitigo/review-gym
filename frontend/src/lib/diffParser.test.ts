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

	it("should parse multi-file diff into separate DiffFile entries", () => {
		const diff = `--- a/a.go
+++ b/a.go
@@ -1 +1 @@
-a
+b
--- a/b.go
+++ b/b.go
@@ -1 +1 @@
-c
+d`;

		const files = parseDiff(diff);
		expect(files).toHaveLength(2);

		expect(files[0]?.oldPath).toBe("a.go");
		expect(files[0]?.newPath).toBe("a.go");
		expect(files[0]?.hunks).toHaveLength(1);
		expect(files[0]?.hunks[0]?.lines).toHaveLength(2);
		expect(files[0]?.hunks[0]?.lines[0]?.type).toBe("deletion");
		expect(files[0]?.hunks[0]?.lines[0]?.content).toBe("a");
		expect(files[0]?.hunks[0]?.lines[1]?.type).toBe("addition");
		expect(files[0]?.hunks[0]?.lines[1]?.content).toBe("b");

		expect(files[1]?.oldPath).toBe("b.go");
		expect(files[1]?.newPath).toBe("b.go");
		expect(files[1]?.hunks).toHaveLength(1);
		expect(files[1]?.hunks[0]?.lines).toHaveLength(2);
		expect(files[1]?.hunks[0]?.lines[0]?.type).toBe("deletion");
		expect(files[1]?.hunks[0]?.lines[0]?.content).toBe("c");
		expect(files[1]?.hunks[0]?.lines[1]?.type).toBe("addition");
		expect(files[1]?.hunks[0]?.lines[1]?.content).toBe("d");
	});

	it("should parse three-file diff with multiple hunks", () => {
		const diff = `--- a/foo.ts
+++ b/foo.ts
@@ -1,2 +1,2 @@
 const x = 1;
-const y = 2;
+const y = 3;
--- a/bar.ts
+++ b/bar.ts
@@ -1,1 +1,2 @@
 export {};
+export const a = 1;
@@ -10,1 +11,1 @@
-old
+new
--- a/baz.ts
+++ b/baz.ts
@@ -1,1 +1,1 @@
-removed
+added`;

		const files = parseDiff(diff);
		expect(files).toHaveLength(3);

		expect(files[0]?.newPath).toBe("foo.ts");
		expect(files[0]?.hunks).toHaveLength(1);

		expect(files[1]?.newPath).toBe("bar.ts");
		expect(files[1]?.hunks).toHaveLength(2);
		expect(files[1]?.hunks[0]?.lines.filter((l) => l.type === "addition")).toHaveLength(1);
		expect(files[1]?.hunks[1]?.lines.filter((l) => l.type === "addition")).toHaveLength(1);
		expect(files[1]?.hunks[1]?.lines.filter((l) => l.type === "deletion")).toHaveLength(1);

		expect(files[2]?.newPath).toBe("baz.ts");
		expect(files[2]?.hunks).toHaveLength(1);
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

	it("should return max across multiple files", () => {
		const diff = `--- a/small.go
+++ b/small.go
@@ -1 +1 @@
-x
+y
--- a/big.go
+++ b/big.go
@@ -50,1 +50,2 @@
 existing
+added`;

		expect(getMaxNewLineNumber(diff)).toBe(51);
	});
});
