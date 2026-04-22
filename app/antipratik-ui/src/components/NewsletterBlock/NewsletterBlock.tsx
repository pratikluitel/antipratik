'use client';

import { useState } from 'react';
import styles from './NewsletterBlock.module.css';
import { subscribeNewsletter } from '@/lib/api';

interface Props {
  variant: 'page' | 'footer';
}

type NLState = 'idle' | 'submitting' | 'success' | 'error';

export default function NewsletterBlock({ variant }: Props) {
  const [nlState, setNLState] = useState<NLState>('idle');
  const [email, setEmail] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
      setErrorMsg('Please enter a valid email address.');
      setNLState('error');
      return;
    }
    setNLState('submitting');
    setErrorMsg('');
    try {
      await subscribeNewsletter(email.trim());
      setNLState('success');
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Something went wrong.');
      setNLState('error');
    }
  }

  return (
    <div className={`${styles.block} ${variant === 'footer' ? styles.footer : ''}`}>
      <h2 className={styles.heading}>Occasionally, something worth sharing.</h2>
      <p className={styles.subtext}>No schedule. No noise. Just the good stuff.</p>

      {nlState === 'success' ? (
        <div className={styles.success}>
          <span className={styles.successIcon}>✓</span>
          <span className={styles.successText}>You&apos;re in! Make sure you check your email to confirm your subscription!</span>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className={styles.form} noValidate>
          <div
            className={`${styles.formWrapper} ${nlState === 'error' ? styles.formWrapperError : ''}`}
          >
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your@email.com"
              required
              disabled={nlState === 'submitting'}
              className={`${styles.input} ${nlState === 'error' ? styles.inputError : ''}`}
              aria-label="Email address"
            />
            <button
              type="submit"
              disabled={nlState === 'submitting'}
              className={styles.button}
            >
              {nlState === 'submitting' ? '...' : 'Subscribe'}
            </button>
          </div>
          {nlState === 'error' && errorMsg && (
            <p className={styles.errorMsg}>{errorMsg}</p>
          )}
        </form>
      )}

      <p className={styles.legal}>No spam, ever. Unsubscribe any time.</p>
    </div>
  );
}
