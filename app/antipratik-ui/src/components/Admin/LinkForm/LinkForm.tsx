'use client';

import { useState, FormEvent } from 'react';
import type { LinkPost } from '@/lib/types';
import { createLinkPost, updateLinkPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';

const CATEGORIES = ['music', 'writing', 'video', 'social'] as const;
type Category = typeof CATEGORIES[number];

interface LinkFormProps {
  token: string;
  initial?: LinkPost;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function LinkForm({ token, initial, onSuccess, onCancel }: LinkFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [url, setUrl] = useState(initial?.url ?? '');
  const [description, setDescription] = useState(initial?.description ?? '');
  const [category, setCategory] = useState<Category | ''>(initial?.category ?? '');
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null);
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const fd = new FormData();
      if (initial) {
        if (title !== initial.title) fd.append('title', title);
        if (url !== initial.url) fd.append('url', url);
        if (description !== (initial.description ?? '')) fd.append('description', description);
        if (category && category !== initial.category) fd.append('category', category);
        tags.forEach((t) => fd.append('tags[]', t));
        await updateLinkPost(initial.id, fd, token);
      } else {
        fd.append('title', title);
        fd.append('url', url);
        if (description) fd.append('description', description);
        if (category) fd.append('category', category);
        if (thumbnailFile) fd.append('thumbnailFile', thumbnailFile);
        tags.forEach((t) => fd.append('tags[]', t));
        await createLinkPost(fd, token);
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
        <label className={`${f.label} ${f.required}`} htmlFor="link-title">Title</label>
        <input id="link-title" className={f.input} value={title} onChange={(e) => setTitle(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${f.required}`} htmlFor="link-url">URL</label>
        <input id="link-url" className={f.input} type="url" value={url} onChange={(e) => setUrl(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label} htmlFor="link-description">Description</label>
        <textarea id="link-description" className={f.textarea} value={description} onChange={(e) => setDescription(e.target.value)} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label} htmlFor="link-category">Category</label>
        <select id="link-category" className={f.select} value={category} onChange={(e) => setCategory(e.target.value as Category | '')} disabled={loading}>
          <option value="">— none —</option>
          {CATEGORIES.map((c) => <option key={c} value={c}>{c}</option>)}
        </select>
      </div>

      <div className={f.field}>
        <label className={f.label} htmlFor="link-thumb">
          Thumbnail {initial && <span className={f.immutableNote}>(cannot be changed after creation)</span>}
        </label>
        {initial ? (
          <p className={f.immutableNote}>Current: {initial.thumbnailUrl || 'none'}</p>
        ) : (
          <input id="link-thumb" className={f.fileInput} type="file" accept=".jpg,.jpeg,.png,.webp" onChange={(e) => setThumbnailFile(e.target.files?.[0] ?? null)} disabled={loading} />
        )}
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update link' : 'Publish link'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>Cancel</button>
      </div>
    </form>
  );
}
