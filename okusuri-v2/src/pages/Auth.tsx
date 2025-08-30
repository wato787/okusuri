import type React from "react";
import { useEffect } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { LoadingSpinner } from "../components/ui/LoadingSpinner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";

export const Auth: React.FC = () => {
	const { isAuthenticated, isLoading, signInWithGoogle } = useAuth();

	useEffect(() => {
		if (isAuthenticated) {
			// Redirect will be handled by Navigate component below
		}
	}, [isAuthenticated]);

	if (isLoading) {
		return (
			<div className="w-full h-screen flex items-center justify-center">
				<LoadingSpinner size="lg" />
			</div>
		);
	}

	if (isAuthenticated) {
		return <Navigate to="/" replace />;
	}

	return (
		<div className="w-full h-screen flex items-center justify-center">
			<div className="flex flex-col gap-6 w-1/2">
				<Card>
					<CardHeader>
						<CardTitle className="text-2xl">ログイン</CardTitle>
						<CardDescription>
							Google アカウントでログインしてください。
						</CardDescription>
					</CardHeader>
					<CardContent>
						<div className="flex flex-col gap-6">
							<Button onClick={signInWithGoogle} className="w-full">
								Googleでログイン
							</Button>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
};