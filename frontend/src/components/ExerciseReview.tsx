import { useExercise } from "@/hooks/useExercise";
import { createReview, listReviews } from "@/lib/api";
import type { Category, ReviewComment } from "@/types/exercise";
import { CATEGORY_LABELS, DIFFICULTY_LABELS } from "@/types/exercise";
import { useCallback, useEffect, useState } from "react";
import { CommentForm } from "./CommentForm";
import { DiffViewer } from "./DiffViewer";

interface ExerciseReviewProps {
	exerciseId: string;
	onBack: () => void;
	onComplete: () => void;
}

export function ExerciseReview({ exerciseId, onBack, onComplete }: ExerciseReviewProps) {
	const { exercise, loading, error } = useExercise(exerciseId);
	const [comments, setComments] = useState<ReviewComment[]>([]);
	const [activeComment, setActiveComment] = useState<{ filePath: string; lineNumber: number } | null>(null);
	const [submitting, setSubmitting] = useState(false);
	const [submitError, setSubmitError] = useState<string | null>(null);

	useEffect(() => {
		listReviews(exerciseId)
			.then((data) => setComments(data))
			.catch(() => {
				/* ignore load error for comments */
			});
	}, [exerciseId]);

	const handleLineClick = useCallback(
		(filePath: string, lineNumber: number) => {
			if (activeComment?.filePath === filePath && activeComment.lineNumber === lineNumber) {
				setActiveComment(null);
			} else {
				setActiveComment({ filePath, lineNumber });
			}
			setSubmitError(null);
		},
		[activeComment],
	);

	const handleSubmitComment = useCallback(
		(data: { content: string; category: Category }) => {
			if (activeComment === null) {
				return;
			}

			setSubmitting(true);
			setSubmitError(null);

			createReview(exerciseId, {
				file_path: activeComment.filePath,
				line_number: activeComment.lineNumber,
				content: data.content,
				category: data.category,
			})
				.then((newComment) => {
					setComments((prev) => [...prev, newComment]);
					setActiveComment(null);
				})
				.catch((err: unknown) => {
					setSubmitError(err instanceof Error ? err.message : "Failed to save comment");
				})
				.finally(() => {
					setSubmitting(false);
				});
		},
		[exerciseId, activeComment],
	);

	const handleCancelComment = useCallback(() => {
		setActiveComment(null);
		setSubmitError(null);
	}, []);

	if (loading) {
		return <p className="loading">Loading exercise...</p>;
	}

	if (error !== null) {
		return (
			<div className="error-container">
				<p className="error">Error: {error}</p>
				<button type="button" onClick={onBack} className="btn-secondary">
					Back to exercises
				</button>
			</div>
		);
	}

	if (exercise === null) {
		return (
			<div className="error-container">
				<p className="error">Exercise not found</p>
				<button type="button" onClick={onBack} className="btn-secondary">
					Back to exercises
				</button>
			</div>
		);
	}

	return (
		<div className="exercise-review">
			<div className="exercise-review-header">
				<button type="button" onClick={onBack} className="btn-back" aria-label="Back to exercise list">
					&larr; Back
				</button>
				<div className="exercise-info">
					<h2>{exercise.title}</h2>
					<div className="exercise-meta">
						<span className={`badge difficulty-${exercise.difficulty}`}>{DIFFICULTY_LABELS[exercise.difficulty]}</span>
						<span className={`badge category-${exercise.category}`}>{CATEGORY_LABELS[exercise.category]}</span>
						<span className="badge language">{exercise.language}</span>
					</div>
					<p className="exercise-description">{exercise.description}</p>
				</div>
			</div>

			<div className="review-instructions">
				<p>
					Click on any line to add a review comment. Identify issues related to security, performance, design,
					readability, or error handling.
				</p>
			</div>

			<div className="diff-section">
				<DiffViewer
					diffContent={exercise.diff_content}
					comments={comments}
					onLineClick={handleLineClick}
					activeCommentLine={activeComment?.lineNumber ?? null}
				/>

				{activeComment !== null && (
					<CommentForm
						filePath={activeComment.filePath}
						lineNumber={activeComment.lineNumber}
						onSubmit={handleSubmitComment}
						onCancel={handleCancelComment}
						submitting={submitting}
					/>
				)}

				{submitError !== null && <p className="error submit-error">Error: {submitError}</p>}
			</div>

			<div className="review-footer">
				<div className="comment-count">
					{String(comments.length)} comment{comments.length !== 1 ? "s" : ""} submitted
				</div>
				<button
					type="button"
					className="btn-primary btn-complete"
					onClick={onComplete}
					disabled={comments.length === 0}
				>
					Complete Review
				</button>
			</div>
		</div>
	);
}
