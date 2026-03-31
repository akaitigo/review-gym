import { useCallback, useEffect, useState } from "react";
import { getExercise } from "@/lib/api";
import type { Exercise } from "@/types/exercise";

interface UseExerciseResult {
	exercise: Exercise | null;
	loading: boolean;
	error: string | null;
	refetch: () => void;
}

/**
 * Hook to fetch a single exercise by ID.
 */
export function useExercise(id: string): UseExerciseResult {
	const [exercise, setExercise] = useState<Exercise | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const fetchData = useCallback(() => {
		setLoading(true);
		setError(null);
		getExercise(id)
			.then((data) => {
				setExercise(data);
			})
			.catch((err: unknown) => {
				setError(err instanceof Error ? err.message : "Failed to load exercise");
			})
			.finally(() => {
				setLoading(false);
			});
	}, [id]);

	useEffect(() => {
		fetchData();
	}, [fetchData]);

	return { exercise, loading, error, refetch: fetchData };
}
