'use client';

import { useState, FormEvent } from 'react';
import type { ExternalLink } from '@/lib/types';
import { createExternalLink, updateExternalLink } from '@/lib/api';
import f from '../adminForm.module.css';

const CATEGORIES = ['music', 'writing', 'video', 'social'] as const;
type Category = typeof CATEGORIES[number];

interface ExternalLinkFormProps {
  token: string;
  initial?: ExternalLink;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function ExternalLinkForm({ token, initial, onSuccess, onCancel }: ExternalLinkFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [url, setUrl] = useState(initial?.url ?? '');
  const [description, setDescription] = useState(initial?.description ?? '');
  const [category, setCategory] = useState<Category>(initial?.category ?? 'writing');
  const [featured, setFeatured] = useState(initial?.featured ?? false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      if (initial) {
        const patch: Record<string, unknown> = {};
        if (title !== initial.title) patch.title = title;
        if (url !== initial.url) patch.url = url;
        if (description !== initial.description) patch.description = description;
        if (category !== initial.category) patch.category = category;
        if (featured !== initial.featured) patch.featured = featured;
        await updateExternalLink(initial.id, patch, token);
      } else {
        await createExternalLink({ title, url, description, category, featured }, token);
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
        <label className={`${f.label} ${f.required}`} htmlFor="el-title">Title</label>
        <input id="el-title" className={f.input} value={title} onChange={(e) => setTitle(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="el-url">URL</label>
        <input id="el-url" className={f.input} type="url" value={url} onChange={(e) => setUrl(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="el-description">Description</label>
        <textarea id="el-description" className={f.textarea} value={description} onChange={(e) => setDescription(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="el-category">Category</label>
        <select id="el-category" className={f.select} value={category} onChange={(e) => setCategory(e.target.value as Category)} required disabled={loading}>
          {CATEGORIES.map((c) => <option key={c} value={c}>{c}</option>)}
        </select>
      </div>

      <div className={f.checkboxRow}>
        <input id="el-featured" className={f.checkbox} type="checkbox" checked={featured} onChange={(e) => setFeatured(e.target.checked)} disabled={loading} />
        <label className={f.label} htmlFor="el-featured">Featured on homepage</label>
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update link' : 'Add link'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>Cancel</button>
      </div>
    </form>
  );
}
