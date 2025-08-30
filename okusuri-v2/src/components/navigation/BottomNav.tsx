import { Calendar, Home, Settings } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import { cn } from '@/lib/utils';

export function BottomNav() {
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  const navItems = [
    {
      name: 'ホーム',
      href: '/',
      icon: Home,
    },
    {
      name: 'カレンダー',
      href: '/calendar',
      icon: Calendar,
    },
    {
      name: '設定',
      href: '/settings',
      icon: Settings,
    },
  ];

  return (
    <div className='fixed bottom-0 left-0 right-0 border-t bg-white border-gray-200'>
      <nav className='flex justify-around items-center h-16'>
        {navItems.map((item) => (
          <Link
            key={item.href}
            to={item.href}
            className={cn(
              'flex flex-col items-center justify-center w-full h-full transition-colors',
              isActive(item.href)
                ? 'text-blue-600'
                : 'text-gray-400 hover:text-gray-600'
            )}
          >
            <item.icon className='h-5 w-5 mb-1' />
            <span className='text-xs'>{item.name}</span>
          </Link>
        ))}
      </nav>
    </div>
  );
}
