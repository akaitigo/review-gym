/** Represents a single line in a parsed unified diff. */
export interface DiffLine {
	/** The type of diff line. */
	type: "addition" | "deletion" | "context" | "header";
	/** The raw content of the line (without the leading +/-/space). */
	content: string;
	/** The line number in the new file (null for deletions and headers). */
	newLineNumber: number | null;
	/** The line number in the old file (null for additions and headers). */
	oldLineNumber: number | null;
}

/** Represents a parsed diff hunk. */
export interface DiffHunk {
	/** The hunk header line (e.g., "@@ -1,5 +1,7 @@"). */
	header: string;
	/** Lines within this hunk. */
	lines: DiffLine[];
}

/** Represents a parsed diff file. */
export interface DiffFile {
	/** Old file path (from "--- a/..." line). */
	oldPath: string;
	/** New file path (from "+++ b/..." line). */
	newPath: string;
	/** Hunks in this file diff. */
	hunks: DiffHunk[];
}

/** Internal mutable state used during parsing. */
interface ParseState {
	files: DiffFile[];
	currentFile: DiffFile | null;
	currentHunk: DiffHunk | null;
	oldLine: number;
	newLine: number;
}

function ensureCurrentFile(state: ParseState): DiffFile {
	if (state.currentFile === null) {
		state.currentFile = { oldPath: "", newPath: "", hunks: [] };
	}
	return state.currentFile;
}

function handleOldFileHeader(state: ParseState, line: string): void {
	const path = line.slice(4).replace(/^a\//, "");
	const file = ensureCurrentFile(state);
	file.oldPath = path;
}

function handleNewFileHeader(state: ParseState, line: string): void {
	const path = line.slice(4).replace(/^b\//, "");
	const file = ensureCurrentFile(state);
	file.newPath = path;
}

function handleHunkHeader(state: ParseState, line: string): void {
	const file = ensureCurrentFile(state);
	const match = /@@ -(\d+)(?:,\d+)? \+(\d+)(?:,\d+)? @@/.exec(line);
	state.oldLine = match ? Number.parseInt(match[1] ?? "0", 10) : 0;
	state.newLine = match ? Number.parseInt(match[2] ?? "0", 10) : 0;
	state.currentHunk = { header: line, lines: [] };
	file.hunks.push(state.currentHunk);
}

function handleContentLine(state: ParseState, line: string): void {
	if (state.currentHunk === null) {
		return;
	}

	if (line.startsWith("+")) {
		state.currentHunk.lines.push({
			type: "addition",
			content: line.slice(1),
			newLineNumber: state.newLine,
			oldLineNumber: null,
		});
		state.newLine++;
		return;
	}

	if (line.startsWith("-")) {
		state.currentHunk.lines.push({
			type: "deletion",
			content: line.slice(1),
			newLineNumber: null,
			oldLineNumber: state.oldLine,
		});
		state.oldLine++;
		return;
	}

	if (line.startsWith(" ") || line === "") {
		state.currentHunk.lines.push({
			type: "context",
			content: line.startsWith(" ") ? line.slice(1) : line,
			newLineNumber: state.newLine,
			oldLineNumber: state.oldLine,
		});
		state.oldLine++;
		state.newLine++;
	}
}

/**
 * Parse a unified diff string into structured file diffs.
 */
export function parseDiff(diffContent: string): DiffFile[] {
	const lines = diffContent.split("\n");
	const state: ParseState = {
		files: [],
		currentFile: null,
		currentHunk: null,
		oldLine: 0,
		newLine: 0,
	};

	for (const line of lines) {
		if (line.startsWith("--- ")) {
			handleOldFileHeader(state, line);
		} else if (line.startsWith("+++ ")) {
			handleNewFileHeader(state, line);
		} else if (line.startsWith("@@")) {
			handleHunkHeader(state, line);
		} else {
			handleContentLine(state, line);
		}
	}

	if (state.currentFile !== null) {
		state.files.push(state.currentFile);
	}

	return state.files;
}

/**
 * Get the maximum new line number from a diff for validation purposes.
 */
export function getMaxNewLineNumber(diffContent: string): number {
	const files = parseDiff(diffContent);
	let max = 0;
	for (const file of files) {
		for (const hunk of file.hunks) {
			for (const line of hunk.lines) {
				if (line.newLineNumber !== null && line.newLineNumber > max) {
					max = line.newLineNumber;
				}
			}
		}
	}
	return max;
}
