'use client';

import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { confirmSubscription } from '@/lib/api';
import styles from './confirm.module.css';

type State = 'loading' | 'success' | 'error';

export default function ConfirmClient({ token }: { token: string }) {
  const [state, setState] = useState<State>('loading');
  const called = useRef(false);

  useEffect(() => {
    if (called.current) return;
    called.current = true;

    if (!token) {
      // Use a microtask to avoid synchronous setState inside effect
      Promise.resolve().then(() => setState('error'));
      return;
    }
    confirmSubscription(token)
      .then(() => setState('success'))
      .catch(() => setState('error'));
  }, [token]);

  return (
    <main className={styles.page}>
      <div className={styles.card}>
        {state === 'loading' && (
          <p className={styles.muted}>Confirming…</p>
        )}
        {state === 'success' && (
          <>
            <span className={styles.icon}>✓</span>
            <h1 className={styles.heading}>You&apos;re in! Make sure you check your email to confirm your subscription!</h1>
            <p className={styles.body}>
              Visit{' '}
              <Link href="/" className={styles.link}>
                antipratik
              </Link>{' '}
              and explore what&apos;s here.
            </p>
          </>
        )}
        {state === 'error' && (
          <>
            <h1 className={styles.heading}>Something went wrong.</h1>
            <p className={styles.body}>
              This confirmation link is invalid or has already been used.
            </p>
          </>
        )}
      </div>
    </main>
  );
}
