// API共通型定義
export interface ApiResponse<T> {
	data: T;
	success: boolean;
	message?: string;
}

export interface ApiError {
	message: string;
	code?: string | number;
	details?: unknown;
}

export interface PaginatedResponse<T> {
	data: T[];
	pagination: {
		page: number;
		limit: number;
		total: number;
		totalPages: number;
	};
}