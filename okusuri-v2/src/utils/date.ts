import { format, parseISO, startOfDay, endOfDay, subDays, addDays } from "date-fns";
import { ja } from "date-fns/locale";

// 日付フォーマット関数
export function formatDate(date: Date | string, formatStr = "yyyy-MM-dd"): string {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	return format(dateObj, formatStr, { locale: ja });
}

// 日本語での日付表示
export function formatDateJapanese(date: Date | string): string {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	return format(dateObj, "M月d日(E)", { locale: ja });
}

// 今日の日付を取得
export function getToday(): string {
	return formatDate(new Date());
}

// 昨日の日付を取得
export function getYesterday(): string {
	return formatDate(subDays(new Date(), 1));
}

// 明日の日付を取得
export function getTomorrow(): string {
	return formatDate(addDays(new Date(), 1));
}

// 日付範囲を取得
export function getDateRange(days: number): { start: string; end: string } {
	const end = new Date();
	const start = subDays(end, days - 1);
	
	return {
		start: formatDate(start),
		end: formatDate(end),
	};
}

// 月の開始日と終了日を取得
export function getMonthRange(date: Date | string = new Date()): { start: string; end: string } {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	const year = dateObj.getFullYear();
	const month = dateObj.getMonth();
	
	const start = new Date(year, month, 1);
	const end = new Date(year, month + 1, 0);
	
	return {
		start: formatDate(start),
		end: formatDate(end),
	};
}

// 日付が今日かどうかチェック
export function isToday(date: Date | string): boolean {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	const today = new Date();
	
	return formatDate(dateObj) === formatDate(today);
}

// 日付の差を計算（日数）
export function getDaysBetween(startDate: Date | string, endDate: Date | string): number {
	const start = typeof startDate === "string" ? parseISO(startDate) : startDate;
	const end = typeof endDate === "string" ? parseISO(endDate) : endDate;
	
	const timeDiff = end.getTime() - start.getTime();
	return Math.ceil(timeDiff / (1000 * 60 * 60 * 24));
}

// ISO文字列からDateオブジェクトに変換
export function parseDate(dateString: string): Date {
	return parseISO(dateString);
}

// 日付の開始時刻を取得
export function getStartOfDay(date: Date | string): Date {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	return startOfDay(dateObj);
}

// 日付の終了時刻を取得
export function getEndOfDay(date: Date | string): Date {
	const dateObj = typeof date === "string" ? parseISO(date) : date;
	return endOfDay(dateObj);
}