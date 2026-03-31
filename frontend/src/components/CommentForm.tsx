import { type FormEvent, useState } from "react";
import { ALL_CATEGORIES, CATEGORY_LABELS, type Category } from "@/types/exercise";

interface CommentFormProps {
	filePath: string;
	lineNumber: number;
	onSubmit: (data: { content: string; category: Category }) => void;
	onCancel: () => void;
	submitting: boolean;
}

export function CommentForm({ filePath, lineNumber, onSubmit, onCancel, submitting }: CommentFormProps) {
	const [content, setContent] = useState("");
	const [category, setCategory] = useState<Category>("security");

	const handleSubmit = (e: FormEvent) => {
		e.preventDefault();
		if (content.trim() === "") {
			return;
		}
		onSubmit({ content: content.trim(), category });
	};

	const contentLength = content.length;
	const isValid = contentLength >= 1 && contentLength <= 5000;

	return (
		<div className="comment-form-container">
			<form className="comment-form" onSubmit={handleSubmit}>
				<div className="comment-form-header">
					<span className="comment-location">
						{filePath} : Line {String(lineNumber)}
					</span>
					<button type="button" className="btn-close" onClick={onCancel} aria-label="Close comment form">
						&times;
					</button>
				</div>

				<div className="form-field">
					<label htmlFor="comment-category">
						Category:
						<select id="comment-category" value={category} onChange={(e) => setCategory(e.target.value as Category)}>
							{ALL_CATEGORIES.map((cat) => (
								<option key={cat} value={cat}>
									{CATEGORY_LABELS[cat]}
								</option>
							))}
						</select>
					</label>
				</div>

				<div className="form-field">
					<label htmlFor="comment-content">
						Comment:
						<textarea
							id="comment-content"
							value={content}
							onChange={(e) => setContent(e.target.value)}
							placeholder="Describe the issue you found..."
							rows={4}
							maxLength={5000}
						/>
					</label>
					<span className="char-count">{String(contentLength)} / 5,000</span>
				</div>

				<div className="comment-form-actions">
					<button type="button" className="btn-secondary" onClick={onCancel} disabled={submitting}>
						Cancel
					</button>
					<button type="submit" className="btn-primary" disabled={!isValid || submitting}>
						{submitting ? "Saving..." : "Submit Comment"}
					</button>
				</div>
			</form>
		</div>
	);
}
