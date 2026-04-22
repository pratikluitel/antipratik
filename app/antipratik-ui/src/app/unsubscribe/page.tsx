import type { Metadata } from 'next';
import { Suspense } from 'react';
import UnsubscribeClient from './UnsubscribeClient';

export const metadata: Metadata = {
  title: 'Unsubscribe — antipratik',
};

export default async function UnsubscribePage({
  searchParams,
}: {
  searchParams: Promise<{ token?: string }>;
}) {
  const { token } = await searchParams;
  return (
    <Suspense>
      <UnsubscribeClient token={token ?? ''} />
    </Suspense>
  );
}
