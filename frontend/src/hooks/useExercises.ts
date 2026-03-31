import { useCallback, useEffect, useState } from "react";
import { listExercises } from "@/lib/api";
import type { Category, Difficulty, ExerciseListItem } from "@/types/exercise";

interface UseExercisesResult {
	exercises: ExerciseListItem[];
	loading: boolean;
	error: string | null;
	refetch: () => void;
}

/**
 * Hook to fetch and filter exercise list.
 */
export function useExercises(filters?: {
	category?: Category | undefined;
	difficulty?: Difficulty | undefined;
}): UseExercisesResult {
	const [exercises, setExercises] = useState<ExerciseListItem[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const fetchData = useCallback(() => {
		setLoading(true);
		setError(null);
		listExercises({
			category: filters?.category,
			difficulty: filters?.difficulty,
		})
			.then((data) => {
				setExercises(data);
			})
			.catch((err: unknown) => {
				setError(err instanceof Error ? err.message : "Failed to load exercises");
			})
			.finally(() => {
				setLoading(false);
			});
	}, [filters?.category, filters?.difficulty]);

	useEffect(() => {
		fetchData();
	}, [fetchData]);

	return { exercises, loading, error, refetch: fetchData };
}
