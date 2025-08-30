import type React from "react";
import type { ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { LoadingSpinner } from "./ui/LoadingSpinner";

interface ProtectedRouteProps {
	children: ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
	children,
}) => {
	const { isAuthenticated, isLoading } = useAuth();

	if (isLoading) {
		return (
			<div className="min-h-screen flex items-center justify-center">
				<LoadingSpinner />
			</div>
		);
	}

	if (!isAuthenticated) {
		return <Navigate to="/auth" replace />;
	}

	return <>{children}</>;
};