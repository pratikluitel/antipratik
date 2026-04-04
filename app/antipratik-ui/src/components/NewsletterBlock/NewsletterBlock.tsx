'use client';

import { useState } from 'react';
import styles from './NewsletterBlock.module.css';

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
    setNLState('submitting');
    setErrorMsg('');
    // Fake 1000ms delay simulating an API call.
    // When wiring a real backend, replace with:
    // await fetch(`${process.env.NEXT_PUBLIC_API_URL}/subscribe`, { method: 'POST', body: JSON.stringify({ email }) })
    await new Promise((r) => setTimeout(r, 1000));
    setNLState('success');
  }

  return (
    <div className={`${styles.block} ${variant === 'footer' ? styles.footer : ''}`}>
      <h2 className={styles.heading}>Occasionally, something worth sharing.</h2>
      <p className={styles.subtext}>No schedule. No noise. Just the good stuff.</p>

      {nlState === 'success' ? (
        <div className={styles.success}>
          <span className={styles.successIcon}>✓</span>
          <span className={styles.successText}>You&apos;re in.</span>
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
