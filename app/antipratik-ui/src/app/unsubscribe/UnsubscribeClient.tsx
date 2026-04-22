'use client';

import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { unsubscribeNewsletter } from '@/lib/api';
import styles from './unsubscribe.module.css';

type State = 'loading' | 'success' | 'error';

export default function UnsubscribeClient({ token }: { token: string }) {
  const [state, setState] = useState<State>('loading');
  const called = useRef(false);

  useEffect(() => {
    if (called.current) return;
    called.current = true;

    if (!token) {
      Promise.resolve().then(() => setState('error'));
      return;
    }
    unsubscribeNewsletter(token)
      .then(() => setState('success'))
      .catch(() => setState('error'));
  }, [token]);

  return (
    <main className={styles.page}>
      <div className={styles.card}>
        {state === 'loading' && (
          <p className={styles.muted}>Unsubscribing…</p>
        )}
        {state === 'success' && (
          <>
            <h1 className={styles.heading}>You&apos;ve been unsubscribed.</h1>
            <p className={styles.body}>
              It&apos;s sad to see you leave :(. If you change your mind, you can always{' '}
              <Link href="/" className={styles.link}>
                subscribe again
              </Link>
              .
            </p>
          </>
        )}
        {state === 'error' && (
          <>
            <h1 className={styles.heading}>Something went wrong.</h1>
            <p className={styles.body}>
              This unsubscribe link is invalid or has already been used.
            </p>
          </>
        )}
      </div>
    </main>
  );
}
