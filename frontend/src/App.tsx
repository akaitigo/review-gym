import { Dashboard } from "@/components/Dashboard";
import { ExerciseList } from "@/components/ExerciseList";
import { ExerciseReview } from "@/components/ExerciseReview";
import { BrowserRouter, Navigate, Route, Routes, useNavigate, useParams } from "react-router-dom";

/** Default user ID for development (no auth yet). */
const DEFAULT_USER_ID = "anonymous";

function ExerciseListPage() {
	const navigate = useNavigate();
	return (
		<main>
			<h1>Review Gym</h1>
			<p>コードレビュースキルを鍛えるトレーニングプラットフォーム</p>
			<div className="nav-actions">
				<button type="button" className="btn-secondary" onClick={() => navigate("/dashboard")}>
					Dashboard
				</button>
			</div>
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

function DashboardPage() {
	const navigate = useNavigate();

	return (
		<main>
			<h1>Review Gym</h1>
			<Dashboard
				userId={DEFAULT_USER_ID}
				onSelectExercise={(id) => navigate(`/exercises/${id}`)}
				onBack={() => navigate("/")}
			/>
		</main>
	);
}

export function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<ExerciseListPage />} />
				<Route path="/exercises/:id" element={<ExerciseReviewPage />} />
				<Route path="/dashboard" element={<DashboardPage />} />
			</Routes>
		</BrowserRouter>
	);
}
