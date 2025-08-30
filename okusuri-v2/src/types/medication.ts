import { z } from "zod";

// 服薬ログのスキーマ
export const medicationLogSchema = z.object({
	id: z.string(),
	userId: z.string(),
	hasBleeding: z.boolean(),
	createdAt: z.string(),
	updatedAt: z.string(),
});

export type MedicationLog = z.infer<typeof medicationLogSchema>;

// 服薬ログリクエストのスキーマ
export const medicationLogRequestSchema = z.object({
	hasBleeding: z.boolean(),
});

export type MedicationLogRequest = z.infer<typeof medicationLogRequestSchema>;

// 薬の状態
export const medicationStatusSchema = z.object({
	currentStreak: z.number(), // 現在の連続服用日数
	isRestPeriod: z.boolean(), // 休薬期間中かどうか
	restDaysLeft: z.number(), // 休薬期間の残り日数
	consecutiveBleedingDays: z.number(), // 連続出血日数
});

export type MedicationStatus = z.infer<typeof medicationStatusSchema>;

// API レスポンス
export const medicationLogListSchema = z.array(medicationLogSchema);
export type MedicationLogList = z.infer<typeof medicationLogListSchema>;