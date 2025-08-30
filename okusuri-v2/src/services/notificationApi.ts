import { apiClient } from '@/lib/apiClient';

export interface NotificationSetting {
  id: string;
  userId: string;
  isEnabled: boolean;
  subscription: string;
  platform: 'ios' | 'android' | 'web';
  createdAt: string;
  updatedAt: string;
}

export interface RegisterNotificationSetting {
  subscription: string;
  isEnabled: boolean;
  platform: 'ios' | 'android' | 'web';
}

export interface BaseResponse {
  success: boolean;
  message?: string;
}

// 通知設定を取得
export const getNotificationSetting = async (): Promise<
  NotificationSetting | undefined
> => {
  try {
    return await apiClient.get<NotificationSetting>(
      '/api/notification/setting'
    );
  } catch (error) {
    console.error('通知設定取得エラー:', error);
    return undefined;
  }
};

// 通知設定を登録
export const registerNotificationSetting = async (
  data: RegisterNotificationSetting
): Promise<BaseResponse> => {
  try {
    await apiClient.post('/api/notification/setting', data);
    return { success: true };
  } catch (error) {
    console.error('通知設定登録エラー:', error);
    return { success: false, message: '通知設定の登録に失敗しました' };
  }
};
