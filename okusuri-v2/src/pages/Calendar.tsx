import type React from "react";
import { MedicationCalendar } from "@/components/calendar/MedicationCalendar";

export const Calendar: React.FC = () => {
	return (
		<div className="container max-w-md mx-auto pt-8 pb-24 px-4">
			<MedicationCalendar />
		</div>
	);
};