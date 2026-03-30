import { type DiffFile, type DiffLine, parseDiff } from "@/lib/diffParser";
import type { Category, ReviewComment } from "@/types/exercise";
import { useMemo } from "react";

interface DiffViewerProps {
	diffContent: string;
	comments: ReviewComment[];
	onLineClick: (filePath: string, lineNumber: number) => void;
	activeCommentLine: number | null;
}

export function DiffViewer({ diffContent, comments, onLineClick, activeCommentLine }: DiffViewerProps) {
	const files = useMemo(() => parseDiff(diffContent), [diffContent]);

	const commentsByLine = useMemo(() => {
		const map = new Map<string, ReviewComment[]>();
		for (const comment of comments) {
			const key = `${comment.file_path}:${String(comment.line_number)}`;
			const existing = map.get(key);
			if (existing !== undefined) {
				existing.push(comment);
			} else {
				map.set(key, [comment]);
			}
		}
		return map;
	}, [comments]);

	return (
		<div className="diff-viewer">
			{files.map((file) => (
				<DiffFileView
					key={file.newPath || file.oldPath}
					file={file}
					commentsByLine={commentsByLine}
					onLineClick={onLineClick}
					activeCommentLine={activeCommentLine}
				/>
			))}
		</div>
	);
}

interface DiffFileViewProps {
	file: DiffFile;
	commentsByLine: Map<string, ReviewComment[]>;
	onLineClick: (filePath: string, lineNumber: number) => void;
	activeCommentLine: number | null;
}

function DiffFileView({ file, commentsByLine, onLineClick, activeCommentLine }: DiffFileViewProps) {
	const filePath = file.newPath || file.oldPath;

	return (
		<div className="diff-file">
			<div className="diff-file-header">{filePath}</div>
			<table className="diff-table">
				<tbody>
					{file.hunks.map((hunk, hunkIdx) => (
						<HunkView
							key={`hunk-${String(hunkIdx)}`}
							header={hunk.header}
							lines={hunk.lines}
							filePath={filePath}
							commentsByLine={commentsByLine}
							onLineClick={onLineClick}
							activeCommentLine={activeCommentLine}
						/>
					))}
				</tbody>
			</table>
		</div>
	);
}

interface HunkViewProps {
	header: string;
	lines: DiffLine[];
	filePath: string;
	commentsByLine: Map<string, ReviewComment[]>;
	onLineClick: (filePath: string, lineNumber: number) => void;
	activeCommentLine: number | null;
}

const CATEGORY_COLORS: Record<Category, string> = {
	security: "#e74c3c",
	performance: "#e67e22",
	design: "#3498db",
	readability: "#2ecc71",
	"error-handling": "#9b59b6",
};

function HunkView({ header, lines, filePath, commentsByLine, onLineClick, activeCommentLine }: HunkViewProps) {
	return (
		<>
			<tr className="diff-hunk-header">
				<td className="diff-line-number" />
				<td className="diff-line-number" />
				<td className="diff-line-content hunk-header">{header}</td>
			</tr>
			{lines.map((line, lineIdx) => {
				const lineKey = `${filePath}:${String(line.newLineNumber)}`;
				const lineComments = line.newLineNumber !== null ? (commentsByLine.get(lineKey) ?? []) : [];
				const isActive = line.newLineNumber === activeCommentLine;
				const isClickable = line.type === "addition" || line.type === "context";

				return (
					<DiffLineView
						key={`line-${String(lineIdx)}`}
						line={line}
						filePath={filePath}
						lineComments={lineComments}
						isActive={isActive}
						isClickable={isClickable}
						onLineClick={onLineClick}
					/>
				);
			})}
		</>
	);
}

interface DiffLineViewProps {
	line: DiffLine;
	filePath: string;
	lineComments: ReviewComment[];
	isActive: boolean;
	isClickable: boolean;
	onLineClick: (filePath: string, lineNumber: number) => void;
}

function DiffLineView({ line, filePath, lineComments, isActive, isClickable, onLineClick }: DiffLineViewProps) {
	const lineClass = `diff-line diff-line-${line.type}${isActive ? " active" : ""}${isClickable ? " clickable" : ""}`;

	const handleClick = () => {
		if (isClickable && line.newLineNumber !== null) {
			onLineClick(filePath, line.newLineNumber);
		}
	};

	const handleKeyDown = (e: React.KeyboardEvent) => {
		if ((e.key === "Enter" || e.key === " ") && isClickable && line.newLineNumber !== null) {
			e.preventDefault();
			onLineClick(filePath, line.newLineNumber);
		}
	};

	const prefix = line.type === "addition" ? "+" : line.type === "deletion" ? "-" : " ";

	return (
		<>
			<tr
				className={lineClass}
				onClick={handleClick}
				onKeyDown={handleKeyDown}
				tabIndex={isClickable ? 0 : undefined}
				role={isClickable ? "button" : undefined}
				aria-label={
					isClickable && line.newLineNumber !== null ? `Add comment on line ${String(line.newLineNumber)}` : undefined
				}
			>
				<td className="diff-line-number old">{line.oldLineNumber ?? ""}</td>
				<td className="diff-line-number new">{line.newLineNumber ?? ""}</td>
				<td className="diff-line-content">
					<code>
						<span className="diff-prefix">{prefix}</span>
						{line.content}
					</code>
				</td>
			</tr>
			{lineComments.map((comment) => (
				<tr key={comment.id} className="diff-comment-row">
					<td className="diff-line-number" />
					<td className="diff-line-number" />
					<td className="diff-comment-cell">
						<div className="diff-comment" style={{ borderLeftColor: CATEGORY_COLORS[comment.category] }}>
							<span className="comment-category">{comment.category}</span>
							<p className="comment-content">{comment.content}</p>
						</div>
					</td>
				</tr>
			))}
		</>
	);
}
