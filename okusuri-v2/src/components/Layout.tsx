import type { ReactNode } from "react";
import { Header } from "./navigation/Header";
import { BottomNav } from "./navigation/BottomNav";
import { useAuth } from "@/contexts/AuthContext";

interface LayoutProps {
	children: ReactNode;
	showNavigation?: boolean;
}

export function Layout({ children, showNavigation = true }: LayoutProps) {
	const { isAuthenticated } = useAuth();

	return (
		<div className="min-h-screen bg-gray-50 flex flex-col">
			{showNavigation && <Header />}
			<main className={`flex-1 ${showNavigation && isAuthenticated ? "pb-16" : ""}`}>
				{children}
			</main>
			{showNavigation && isAuthenticated && <BottomNav />}
		</div>
	);
}