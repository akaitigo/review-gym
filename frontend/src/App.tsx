import { ExerciseList } from "@/components/ExerciseList";
import { ExerciseReview } from "@/components/ExerciseReview";
import { BrowserRouter, Navigate, Route, Routes, useNavigate, useParams } from "react-router-dom";

function ExerciseListPage() {
	const navigate = useNavigate();
	return (
		<main>
			<h1>Review Gym</h1>
			<p>コードレビュースキルを鍛えるトレーニングプラットフォーム</p>
			<ExerciseList onSelect={(id) => navigate(`/exercises/${id}`)} />
		</main>
	);
}

function ExerciseReviewPage() {
	const { id } = useParams<{ id: string }>();
	const navigate = useNavigate();

	if (id === undefined) {
		return <Navigate to="/" replace />;
	}

	return (
		<main>
			<ExerciseReview exerciseId={id} onBack={() => navigate("/")} onComplete={() => navigate("/")} />
		</main>
	);
}

export function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<ExerciseListPage />} />
				<Route path="/exercises/:id" element={<ExerciseReviewPage />} />
			</Routes>
		</BrowserRouter>
	);
}
