import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Download, X } from "lucide-react";

interface BeforeInstallPromptEvent extends Event {
	prompt(): Promise<void>;
	userChoice: Promise<{ outcome: "accepted" | "dismissed"; platform: string }>;
}

export function PWAInstallPrompt() {
	const [deferredPrompt, setDeferredPrompt] =
		useState<BeforeInstallPromptEvent | null>(null);
	const [showInstallPrompt, setShowInstallPrompt] = useState(false);

	useEffect(() => {
		const handleBeforeInstallPrompt = (e: Event) => {
			// Prevent the mini-infobar from appearing on mobile
			e.preventDefault();
			// Stash the event so it can be triggered later.
			setDeferredPrompt(e as BeforeInstallPromptEvent);
			// Show the install prompt
			setShowInstallPrompt(true);
		};

		window.addEventListener("beforeinstallprompt", handleBeforeInstallPrompt);

		return () => {
			window.removeEventListener(
				"beforeinstallprompt",
				handleBeforeInstallPrompt,
			);
		};
	}, []);

	const handleInstallClick = async () => {
		if (deferredPrompt) {
			// Show the install prompt
			await deferredPrompt.prompt();
			// Wait for the user to respond to the prompt
			const { outcome } = await deferredPrompt.userChoice;
			
			if (outcome === "accepted") {
				console.log("User accepted the install prompt");
			} else {
				console.log("User dismissed the install prompt");
			}
			
			// Clear the deferredPrompt
			setDeferredPrompt(null);
			setShowInstallPrompt(false);
		}
	};

	const handleDismiss = () => {
		setShowInstallPrompt(false);
		// Hide for this session
		setDeferredPrompt(null);
	};

	if (!showInstallPrompt || !deferredPrompt) {
		return null;
	}

	return (
		<div className="fixed bottom-20 left-4 right-4 z-50 md:left-auto md:right-4 md:max-w-sm">
			<Card className="shadow-lg border border-blue-200">
				<CardHeader className="pb-2">
					<div className="flex items-center justify-between">
						<CardTitle className="text-base flex items-center gap-2">
							<Download className="h-4 w-4" />
							アプリをインストール
						</CardTitle>
						<Button
							variant="ghost"
							size="icon"
							className="h-6 w-6"
							onClick={handleDismiss}
						>
							<X className="h-4 w-4" />
						</Button>
					</div>
				</CardHeader>
				<CardContent className="pt-0">
					<p className="text-sm text-gray-600 mb-3">
						Okusuriをホーム画面に追加して、より便利に使用できます。
					</p>
					<div className="flex gap-2">
						<Button onClick={handleInstallClick} size="sm" className="flex-1">
							インストール
						</Button>
						<Button
							variant="outline"
							size="sm"
							onClick={handleDismiss}
							className="flex-1"
						>
							後で
						</Button>
					</div>
				</CardContent>
			</Card>
		</div>
	);
}