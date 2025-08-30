import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import {
	addMonths,
	eachDayOfInterval,
	endOfMonth,
	format,
	isSameDay,
	parseISO,
	startOfMonth,
	subMonths,
} from "date-fns";
import { ja } from "date-fns/locale";
import { Check, ChevronLeft, ChevronRight, Droplet, X } from "lucide-react";
import { useState } from "react";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import {
	useCreateMedicationLog,
	useUpdateMedicationLog,
	useMedicationLogs,
} from "@/hooks/useMedicationData";

export function MedicationCalendar() {
	const [currentMonth, setCurrentMonth] = useState(new Date());
	const [selectedDay, setSelectedDay] = useState<Date | undefined>(undefined);
	const [activeButton, setActiveButton] = useState<
		"bleeding" | "normal" | null
	>(null);

	const { data: logs = [], isLoading } = useMedicationLogs();
	const createLogMutation = useCreateMedicationLog();
	const updateLogMutation = useUpdateMedicationLog();

	// 選択された日付のログを取得
	const selectedDayLog = selectedDay
		? logs.find((log) => isSameDay(parseISO(log.createdAt), selectedDay))
		: undefined;

	// 月の最初と最後の日を取得
	const monthStart = startOfMonth(currentMonth);
	const monthEnd = endOfMonth(currentMonth);

	// 週の始まりの日を取得（日曜日から始まるように調整）
	const startDate = new Date(monthStart);
	startDate.setDate(startDate.getDate() - startDate.getDay());

	// 週の終わりの日を取得（土曜日で終わるように調整）
	const endDate = new Date(monthEnd);
	const daysToAdd = 6 - endDate.getDay();
	endDate.setDate(endDate.getDate() + daysToAdd);

	// カレンダーに表示する全ての日付を取得
	const calendarDays = eachDayOfInterval({ start: startDate, end: endDate });

	// 前の月へ
	const prevMonth = () => {
		setCurrentMonth(subMonths(currentMonth, 1));
	};

	// 次の月へ
	const nextMonth = () => {
		setCurrentMonth(addMonths(currentMonth, 1));
	};

	// 日付をクリックしたときの処理
	const handleDayClick = (day: Date) => {
		setSelectedDay(day);
		setActiveButton(null);
	};

	// 日付に対応するログを取得
	const getLogForDay = (day: Date) => {
		return logs.find((log) => isSameDay(parseISO(log.createdAt), day));
	};

	// 記録を追加または更新する処理
	const handleAddRecord = async (hasBleeding: boolean) => {
		if (!selectedDay) return;

		// クリックされたボタンをアクティブにする
		setActiveButton(hasBleeding ? "bleeding" : "normal");

		try {
			// 既存のログがあれば更新、なければ新規登録
			if (selectedDayLog) {
				await updateLogMutation.mutateAsync({
					id: selectedDayLog.id,
					data: { hasBleeding },
				});
			} else {
				await createLogMutation.mutateAsync({ hasBleeding });
			}
		} catch (error) {
			console.error("記録の処理に失敗しました", error);
		} finally {
			setActiveButton(null);
		}
	};

	// 曜日の配列
	const weekdays = ["日", "月", "火", "水", "木", "金", "土"];

	// 未来の日付かどうかを判定
	const isFutureDate = (date: Date) => {
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		return date > today;
	};

	// 選択を解除する
	const clearSelection = () => {
		setSelectedDay(undefined);
	};

	const isPending =
		createLogMutation.isPending || updateLogMutation.isPending;

	if (isLoading) {
		return (
			<div className="flex justify-center items-center min-h-96">
				<LoadingSpinner size="lg" />
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<Card className="overflow-hidden shadow-lg border-0 rounded-xl">
				<CardHeader className="pb-2">
					<div className="flex items-center justify-between">
						<Button
							variant="outline"
							size="icon"
							onClick={prevMonth}
							aria-label="前月へ"
							className="rounded-full h-9 w-9"
						>
							<ChevronLeft className="h-4 w-4" />
						</Button>
						<CardTitle className="text-center text-lg">
							{format(currentMonth, "yyyy年MM月", { locale: ja })}
						</CardTitle>
						<Button
							variant="outline"
							size="icon"
							onClick={nextMonth}
							aria-label="次月へ"
							className="rounded-full h-9 w-9"
						>
							<ChevronRight className="h-4 w-4" />
						</Button>
					</div>
				</CardHeader>
				<CardContent className="p-4">
					<div className="grid grid-cols-7 gap-1">
						{/* 曜日のヘッダー */}
						{weekdays.map((day, index) => (
							<div
								key={day}
								className={cn(
									"text-center text-sm font-medium py-2",
									index === 0
										? "text-red-500"
										: index === 6
											? "text-blue-500"
											: "",
								)}
							>
								{day}
							</div>
						))}

						{/* カレンダーの日付 */}
						{calendarDays.map((day) => {
							const dayLog = getLogForDay(day);
							const isCurrentMonth = day.getMonth() === currentMonth.getMonth();
							const isToday = isSameDay(day, new Date());
							const isSelected = selectedDay && isSameDay(day, selectedDay);
							const isFuture = isFutureDate(day);

							return (
								<button
									key={day.toString()}
									type="button"
									className={cn(
										"relative h-10 flex items-center justify-center text-sm rounded-md transition-colors hover:scale-105 active:scale-95",
										!isCurrentMonth && "text-gray-400 opacity-40",
										isToday && "border border-blue-600",
										isSelected && "bg-blue-100 font-medium",
										dayLog &&
											!dayLog.hasBleeding &&
											"bg-green-100 text-green-800",
										dayLog?.hasBleeding && "bg-red-100 text-red-800",
										isFuture && "opacity-50 cursor-not-allowed",
									)}
									onClick={() => !isFuture && handleDayClick(day)}
									disabled={isFuture}
									aria-label={`${format(day, "yyyy年MM月dd日")}${
										dayLog
											? dayLog.hasBleeding
												? "、出血あり"
												: "、正常に服用"
											: ""
									}${isFuture ? "、未来の日付" : ""}`}
									aria-pressed={isSelected}
								>
									<span>{format(day, "d")}</span>
									{dayLog && (
										<div className="absolute bottom-0.5 left-1/2 transform -translate-x-1/2">
											{dayLog.hasBleeding ? (
												<Droplet className="h-3 w-3 text-red-500" />
											) : (
												<Check className="h-3 w-3 text-green-500" />
											)}
										</div>
									)}
								</button>
							);
						})}
					</div>
				</CardContent>
			</Card>

			{/* 選択された日付の詳細 */}
			{selectedDay && !isFutureDate(selectedDay) && (
				<Card className="overflow-hidden shadow-lg border-0 rounded-xl">
					<CardHeader className="pb-2 pt-4 px-4 flex flex-row items-center justify-between">
						<CardTitle className="text-center text-base">
							{format(selectedDay, "yyyy年MM月dd日 (eee)", { locale: ja })}
						</CardTitle>
						<Button
							variant="ghost"
							size="icon"
							className="h-8 w-8 rounded-full"
							onClick={clearSelection}
						>
							<X className="h-4 w-4" />
						</Button>
					</CardHeader>
					<CardContent className="px-4 pb-4">
						{selectedDayLog ? (
							<div className="flex flex-col items-center justify-center p-2">
								<div className="flex items-center justify-center mb-2">
									{selectedDayLog.hasBleeding ? (
										<div className="flex items-center text-red-500">
											<Droplet className="mr-2 h-5 w-5" />
											<span>出血あり</span>
										</div>
									) : (
										<div className="flex items-center text-green-500">
											<Check className="mr-2 h-5 w-5" />
											<span>正常に服用</span>
										</div>
									)}
								</div>
								<div className="text-sm text-gray-600">
									記録時間:{" "}
									{format(parseISO(selectedDayLog.createdAt), "HH:mm")}
								</div>

								<div className="grid grid-cols-2 gap-3 w-full mt-4">
									<Button
										variant="outline"
										className={cn(
											"flex-1 border-2",
											!selectedDayLog.hasBleeding
												? "bg-green-100 border-green-300"
												: "border-green-200 hover:border-green-300",
										)}
										onClick={() => handleAddRecord(false)}
										disabled={
											isPending || (!selectedDayLog.hasBleeding && !isPending)
										}
									>
										{isPending && activeButton === "normal" ? (
											<LoadingSpinner size="sm" />
										) : (
											<Check className="mr-2 h-4 w-4 text-green-500" />
										)}
										正常に服用
									</Button>
									<Button
										variant="outline"
										className={cn(
											"flex-1 border-2",
											selectedDayLog.hasBleeding
												? "bg-red-100 border-red-300"
												: "border-red-200 hover:border-red-300",
										)}
										onClick={() => handleAddRecord(true)}
										disabled={
											isPending || (selectedDayLog.hasBleeding && !isPending)
										}
									>
										{isPending && activeButton === "bleeding" ? (
											<LoadingSpinner size="sm" />
										) : (
											<Droplet className="mr-2 h-4 w-4 text-red-500" />
										)}
										出血あり
									</Button>
								</div>
							</div>
						) : (
							<div className="grid grid-cols-2 gap-3 w-full">
								<Button
									variant="outline"
									className={cn(
										"flex-1 h-14 w-full border-2",
										activeButton === "normal"
											? "bg-green-100 border-green-300"
											: "border-green-200 hover:border-green-300",
									)}
									onClick={() => handleAddRecord(false)}
									disabled={isPending}
								>
									{isPending && activeButton === "normal" ? (
										<LoadingSpinner size="sm" />
									) : (
										<Check className="mr-2 h-4 w-4 text-green-500" />
									)}
									正常に服用
								</Button>
								<Button
									variant="outline"
									className={cn(
										"flex-1 h-14 w-full border-2",
										activeButton === "bleeding"
											? "bg-red-100 border-red-300"
											: "border-red-200 hover:border-red-300",
									)}
									onClick={() => handleAddRecord(true)}
									disabled={isPending}
								>
									{isPending && activeButton === "bleeding" ? (
										<LoadingSpinner size="sm" />
									) : (
										<Droplet className="mr-2 h-4 w-4 text-red-500" />
									)}
									出血あり
								</Button>
							</div>
						)}
					</CardContent>
				</Card>
			)}

			{/* 凡例 */}
			<div className="flex justify-center space-x-4">
				<div className="flex items-center">
					<div className="w-3 h-3 rounded-full bg-green-500 mr-2" />
					<span className="text-sm">正常に服用</span>
				</div>
				<div className="flex items-center">
					<div className="w-3 h-3 rounded-full bg-red-500 mr-2" />
					<span className="text-sm">出血あり</span>
				</div>
			</div>
		</div>
	);
}