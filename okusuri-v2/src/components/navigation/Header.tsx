export function Header() {
  return (
    <header className='sticky top-0 z-50 w-full border-b bg-white/95 backdrop-blur border-gray-200'>
      <div className='container mx-auto flex h-14 items-center justify-between relative px-4'>
        <div className='absolute left-0'>
          <div className='w-10' />
        </div>

        <div className='font-medium absolute left-1/2 transform -translate-x-1/2'>
          おくすり管理
        </div>

        <div className='absolute right-0'>{/* 何も表示しない */}</div>
      </div>
    </header>
  );
}
