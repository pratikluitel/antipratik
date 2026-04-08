'use client';

import { useState, FormEvent } from 'react';
import type { ShortPost } from '@/lib/types';
import { createShortPost, updateShortPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';
import styles from './ShortPostForm.module.css';

const MAX_CHARS = 320;

interface ShortPostFormProps {
  token: string;
  initial?: ShortPost;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function ShortPostForm({ token, initial, onSuccess, onCancel }: ShortPostFormProps) {
  const [body, setBody] = useState(initial?.body ?? '');
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const over = body.length > MAX_CHARS;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (over) return;
    setError(null);
    setLoading(true);
    try {
      if (initial) {
        const patch: Record<string, unknown> = { tags };
        if (body !== initial.body) patch.body = body;
        await updateShortPost(initial.id, patch, token);
      } else {
        await createShortPost({ body, tags }, token);
      }
      setSuccess(true);
      setTimeout(onSuccess, 800);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong.');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form className={f.form} onSubmit={handleSubmit} noValidate>
      {error && <p className={f.error} role="alert">{error}</p>}
      {success && <p className={f.success}>Saved successfully.</p>}

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="short-body">Post</label>
        <textarea
          id="short-body"
          className={f.textarea}
          value={body}
          onChange={(e) => setBody(e.target.value)}
          required
          disabled={loading}
          style={{ minHeight: '120px' }}
        />
        <p className={`${f.charCount} ${over ? f.charCountOver : ''} ${styles.counter}`}>
          {body.length} / {MAX_CHARS}
        </p>
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || over || success}>
          {loading ? 'Saving…' : initial ? 'Update post' : 'Publish post'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>
          Cancel
        </button>
      </div>
    </form>
  );
}
