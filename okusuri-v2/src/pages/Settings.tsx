import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { NotificationSetting } from '@/components/NotificationSetting';

export const Settings: React.FC = () => {
  const { signOut } = useAuth();

  return (
    <div className='container max-w-md mx-auto pt-8 pb-24 px-4 space-y-6'>
      <NotificationSetting />

      <Card>
        <CardHeader>
          <CardTitle>アカウント</CardTitle>
        </CardHeader>
        <CardContent>
          <Button variant='destructive' onClick={signOut}>
            ログアウト
          </Button>
        </CardContent>
      </Card>
    </div>
  );
};
