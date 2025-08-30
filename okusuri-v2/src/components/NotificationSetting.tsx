import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Bell, BellOff } from 'lucide-react';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import toast from 'react-hot-toast';

import { initializeWebPush } from '@/lib/webpush';
import {
  registerNotificationSetting,
  getNotificationSetting,
} from '@/services/notificationApi';

interface Props {
  onSettingChange?: () => void;
}

export function NotificationSetting({ onSettingChange }: Props) {
  const queryClient = useQueryClient();

  // 通知設定を取得
  const { data: notificationSetting, isLoading } = useQuery({
    queryKey: ['notification-setting'],
    queryFn: getNotificationSetting,
  });

  // 通知設定登録のミューテーション
  const registerMutation = useMutation({
    mutationFn: registerNotificationSetting,
    onSuccess: () => {
      toast.success('通知の設定が完了しました');
      queryClient.invalidateQueries({ queryKey: ['notification-setting'] });
      onSettingChange?.();
    },
    onError: (error) => {
      console.error(error);
      toast.error('通知の設定に失敗しました');
    },
  });

  const handleNotificationSetting = async () => {
    try {
      // Web Pushの初期化とサブスクリプション取得
      const subscription = await initializeWebPush();

      if (!subscription) {
        toast.error('通知の設定に失敗しました。通知の許可を確認してください。');
        return;
      }

      // 設定を登録
      registerMutation.mutate({
        subscription,
        isEnabled: true,
        platform: 'web',
      });
    } catch (error) {
      console.error(error);
      toast.error('通知の設定に失敗しました');
    }
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className='flex items-center gap-2'>
            <Bell className='h-5 w-5 text-amber-500' />
            <span>通知設定</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex justify-center py-8'>
            <LoadingSpinner size='md' />
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className='flex items-center gap-2'>
          <Bell className='h-5 w-5 text-amber-500' />
          <span>通知設定</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className='text-sm text-muted-foreground mb-6'>
          {notificationSetting
            ? '通知設定が完了しています。お薬の服用時間になると通知が届きます。'
            : 'お薬の服用時間になると通知でお知らせします。下のボタンから通知を有効にしてください。'}
        </div>

        <Button
          variant={notificationSetting ? 'outline' : 'default'}
          className='w-full h-12 text-lg font-medium'
          onClick={handleNotificationSetting}
          disabled={!!notificationSetting || registerMutation.isPending}
        >
          {registerMutation.isPending ? (
            <>
              <LoadingSpinner size='sm' className='mr-2' />
              設定中...
            </>
          ) : notificationSetting ? (
            <>
              <BellOff className='mr-2 h-5 w-5' />
              通知設定済み
            </>
          ) : (
            <>
              <Bell className='mr-2 h-5 w-5' />
              PUSH通知を有効にする
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  );
}
