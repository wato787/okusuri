import type React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  useMedicationStatus,
  useCreateMedicationLog,
} from '@/hooks/useMedicationData';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Plus } from 'lucide-react';

export const Home: React.FC = () => {
  const { data: medicationStatus, isLoading: statusLoading } =
    useMedicationStatus();
  const createLogMutation = useCreateMedicationLog();

  const handleAddMedicationLog = (hasBleeding: boolean) => {
    createLogMutation.mutate({ hasBleeding });
  };

  if (statusLoading) {
    return (
      <div className='p-4 flex justify-center items-center min-h-96'>
        <LoadingSpinner size='lg' />
      </div>
    );
  }

  return (
    <div className='container max-w-md mx-auto pt-6 pb-24 px-4'>
      {/* 今日のアクション */}
      <Card className='mb-6'>
        <CardHeader>
          <CardTitle className='flex items-center gap-2'>
            <Plus className='h-5 w-5' />
            今日の服薬記録
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex gap-3'>
            <Button
              onClick={() => handleAddMedicationLog(false)}
              disabled={createLogMutation.isPending}
              className='flex-1'
            >
              {createLogMutation.isPending ? (
                <LoadingSpinner size='sm' />
              ) : (
                '服薬済み'
              )}
            </Button>
            <Button
              onClick={() => handleAddMedicationLog(true)}
              disabled={createLogMutation.isPending}
              variant='outline'
              className='flex-1'
            >
              {createLogMutation.isPending ? (
                <LoadingSpinner size='sm' />
              ) : (
                '服薬済み + 出血'
              )}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* 服薬状況 */}
      {medicationStatus && (
        <Card className='mb-6'>
          <CardHeader>
            <CardTitle>現在の服薬状況</CardTitle>
          </CardHeader>
          <CardContent>
            <div className='grid grid-cols-2 gap-4'>
              <div className='text-center'>
                <div className='text-2xl font-bold text-blue-600'>
                  {medicationStatus.currentStreak}
                </div>
                <div className='text-sm text-gray-600'>連続服用日数</div>
              </div>
              <div className='text-center'>
                <div className='text-2xl font-bold text-purple-600'>
                  {medicationStatus.isRestPeriod
                    ? medicationStatus.restDaysLeft
                    : '0'}
                </div>
                <div className='text-sm text-gray-600'>
                  {medicationStatus.isRestPeriod
                    ? '休薬残り日数'
                    : '休薬期間外'}
                </div>
              </div>
            </div>
            {medicationStatus.consecutiveBleedingDays > 0 && (
              <div className='mt-4 p-3 bg-red-50 rounded-lg'>
                <div className='text-center'>
                  <div className='text-lg font-semibold text-red-700'>
                    連続出血 {medicationStatus.consecutiveBleedingDays} 日
                  </div>
                  <div className='text-sm text-red-600'>
                    体調に注意してください
                  </div>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
};
