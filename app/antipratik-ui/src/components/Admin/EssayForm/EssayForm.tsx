'use client';

import { useState, FormEvent } from 'react';
import type { EssayPost } from '@/lib/types';
import { createEssay, updateEssay } from '@/lib/api';
import TagInput from '../TagInput';
import MarkdownEditor from '../MarkdownEditor';
import f from '../adminForm.module.css';

interface EssayFormProps {
  token: string;
  initial?: EssayPost;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function EssayForm({ token, initial, onSuccess, onCancel }: EssayFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [slug, setSlug] = useState(initial?.slug ?? '');
  const [excerpt, setExcerpt] = useState(initial?.excerpt ?? '');
  const [body, setBody] = useState(initial?.body ?? '');
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  function autoSlug(value: string) {
    return value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
  }

  function handleTitleChange(value: string) {
    setTitle(value);
    if (!initial) setSlug(autoSlug(value));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      if (initial) {
        const patch: Record<string, unknown> = {};
        if (title !== initial.title) patch.title = title;
        if (slug !== initial.slug) patch.slug = slug;
        if (excerpt !== initial.excerpt) patch.excerpt = excerpt;
        if (body !== initial.body) patch.body = body;
        patch.tags = tags;
        await updateEssay(initial.id, patch, token);
      } else {
        await createEssay({ title, slug, excerpt, body, tags }, token);
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
        <label className={`${f.label} ${f.required}`} htmlFor="essay-title">Title</label>
        <input id="essay-title" className={f.input} value={title} onChange={(e) => handleTitleChange(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="essay-slug">Slug</label>
        <input id="essay-slug" className={f.input} value={slug} onChange={(e) => setSlug(e.target.value)} required disabled={loading} />
        <p className={f.hint}>URL path: /{slug || '…'}</p>
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="essay-excerpt">Excerpt</label>
        <textarea id="essay-excerpt" className={f.textarea} value={excerpt} onChange={(e) => setExcerpt(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`}>Body</label>
        <MarkdownEditor value={body} onChange={setBody} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update essay' : 'Publish essay'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>
          Cancel
        </button>
      </div>
    </form>
  );
}
