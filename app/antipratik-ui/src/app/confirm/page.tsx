import type { Metadata } from 'next';
import { Suspense } from 'react';
import ConfirmClient from './ConfirmClient';

export const metadata: Metadata = {
  title: 'Confirm subscription — antipratik',
};

export default async function ConfirmPage({
  searchParams,
}: {
  searchParams: Promise<{ token?: string }>;
}) {
  const { token } = await searchParams;
  return (
    <Suspense>
      <ConfirmClient token={token ?? ''} />
    </Suspense>
  );
}
