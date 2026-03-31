import { useState } from "react";
import { useExercises } from "@/hooks/useExercises";
import {
	ALL_CATEGORIES,
	ALL_DIFFICULTIES,
	CATEGORY_LABELS,
	type Category,
	DIFFICULTY_LABELS,
	type Difficulty,
} from "@/types/exercise";

interface ExerciseListProps {
	onSelect: (id: string) => void;
}

export function ExerciseList({ onSelect }: ExerciseListProps) {
	const [categoryFilter, setCategoryFilter] = useState<Category | undefined>(undefined);
	const [difficultyFilter, setDifficultyFilter] = useState<Difficulty | undefined>(undefined);

	const { exercises, loading, error } = useExercises({
		category: categoryFilter,
		difficulty: difficultyFilter,
	});

	return (
		<div className="exercise-list">
			<h2>Practice Exercises</h2>

			<div className="filters">
				<label htmlFor="category-filter">
					Category:
					<select
						id="category-filter"
						value={categoryFilter ?? ""}
						onChange={(e) => {
							const val = e.target.value;
							setCategoryFilter(val === "" ? undefined : (val as Category));
						}}
					>
						<option value="">All Categories</option>
						{ALL_CATEGORIES.map((cat) => (
							<option key={cat} value={cat}>
								{CATEGORY_LABELS[cat]}
							</option>
						))}
					</select>
				</label>

				<label htmlFor="difficulty-filter">
					Difficulty:
					<select
						id="difficulty-filter"
						value={difficultyFilter ?? ""}
						onChange={(e) => {
							const val = e.target.value;
							setDifficultyFilter(val === "" ? undefined : (val as Difficulty));
						}}
					>
						<option value="">All Difficulties</option>
						{ALL_DIFFICULTIES.map((diff) => (
							<option key={diff} value={diff}>
								{DIFFICULTY_LABELS[diff]}
							</option>
						))}
					</select>
				</label>
			</div>

			{loading && <p className="loading">Loading exercises...</p>}
			{error !== null && <p className="error">Error: {error}</p>}

			{!loading && exercises.length === 0 && <p className="empty">No exercises found matching your filters.</p>}

			<ul className="exercise-items">
				{exercises.map((ex) => (
					<li key={ex.id} className="exercise-item">
						<button type="button" className="exercise-card" onClick={() => onSelect(ex.id)}>
							<div className="exercise-header">
								<h3>{ex.title}</h3>
								<div className="exercise-meta">
									<span className={`badge difficulty-${ex.difficulty}`}>{DIFFICULTY_LABELS[ex.difficulty]}</span>
									<span className={`badge category-${ex.category}`}>{CATEGORY_LABELS[ex.category]}</span>
									<span className="badge language">{ex.language}</span>
								</div>
							</div>
							<p className="exercise-description">{ex.description}</p>
							{ex.category_tags.length > 1 && (
								<div className="exercise-tags">
									{ex.category_tags
										.filter((tag) => tag !== ex.category)
										.map((tag) => (
											<span key={tag} className="tag">
												{CATEGORY_LABELS[tag]}
											</span>
										))}
								</div>
							)}
						</button>
					</li>
				))}
			</ul>
		</div>
	);
}
