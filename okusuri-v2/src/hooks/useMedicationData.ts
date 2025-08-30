import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/apiClient';
import type {
  MedicationLog,
  MedicationLogRequest,
  MedicationStatus,
} from '@/types/medication';

// 服薬ログを取得
export const useMedicationLogs = () => {
  return useQuery({
    queryKey: ['medication-logs'],
    queryFn: () => apiClient.get<MedicationLog[]>('/api/medication-log'),
  });
};

// 服薬ログを作成
export const useCreateMedicationLog = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: MedicationLogRequest) =>
      apiClient.post<MedicationLog>('/api/medication-log', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medication-logs'] });
      queryClient.invalidateQueries({ queryKey: ['medication-status'] });
    },
  });
};

// 服薬ログを更新
export const useUpdateMedicationLog = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: Partial<MedicationLogRequest>;
    }) => apiClient.patch<MedicationLog>(`/api/medication-log/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medication-logs'] });
      queryClient.invalidateQueries({ queryKey: ['medication-status'] });
    },
  });
};

// 服薬ログを削除
export const useDeleteMedicationLog = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => apiClient.delete(`/api/medication-log/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medication-logs'] });
      queryClient.invalidateQueries({ queryKey: ['medication-status'] });
    },
  });
};

// 服薬ステータスを取得
export const useMedicationStatus = () => {
  return useQuery({
    queryKey: ['medication-status'],
    queryFn: () => apiClient.get<MedicationStatus>('/api/medication-status'),
  });
};

// 特定期間の服薬ログを取得
export const useMedicationLogsByDateRange = (
  startDate: string,
  endDate: string
) => {
  return useQuery({
    queryKey: ['medication-logs', startDate, endDate],
    queryFn: () =>
      apiClient.get<MedicationLog[]>(
        `/api/medication-log?start=${startDate}&end=${endDate}`
      ),
    enabled: !!startDate && !!endDate,
  });
};
