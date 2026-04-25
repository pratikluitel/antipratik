'use client';

import { useState, FormEvent } from 'react';
import type { VideoPost } from '@/lib/types';
import { createVideoPost, updateVideoPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';

interface VideoFormProps {
  token: string;
  initial?: VideoPost;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function VideoForm({ token, initial, onSuccess, onCancel }: VideoFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [description, setDescription] = useState(initial?.description ?? '');
  const [videoFile, setVideoFile] = useState<File | null>(null);
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
        if (description !== (initial.description ?? '')) fd.append('description', description);
        if (thumbnailFile) fd.append('thumbnailFile', thumbnailFile);
        tags.forEach((t) => fd.append('tags[]', t));
        await updateVideoPost(initial.id, fd, token);
      } else {
        if (!videoFile) {
          setError('A video file is required.');
          setLoading(false);
          return;
        }
        fd.append('title', title);
        if (description) fd.append('description', description);
        fd.append('videoFile', videoFile);
        if (thumbnailFile) fd.append('thumbnailFile', thumbnailFile);
        tags.forEach((t) => fd.append('tags[]', t));
        await createVideoPost(fd, token);
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
        <label className={`${f.label} ${f.required}`} htmlFor="video-title">Title</label>
        <input id="video-title" className={f.input} value={title} onChange={(e) => setTitle(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label} htmlFor="video-description">Description</label>
        <textarea id="video-description" className={f.textarea} value={description} onChange={(e) => setDescription(e.target.value)} disabled={loading} rows={3} />
      </div>

      {!initial && (
        <div className={f.field}>
          <label className={`${f.label} ${f.required}`} htmlFor="video-file">Video file</label>
          <input
            id="video-file"
            className={f.fileInput}
            type="file"
            accept="video/mp4,video/webm,video/quicktime"
            onChange={(e) => setVideoFile(e.target.files?.[0] ?? null)}
            required
            disabled={loading}
          />
        </div>
      )}

      <div className={f.field}>
        <label className={f.label} htmlFor="video-thumb">Thumbnail</label>
        {initial && <p className={f.immutableNote}>Current: {initial.thumbnailUrl ?? 'none'}</p>}
        <input id="video-thumb" className={f.fileInput} type="file" accept=".jpg,.jpeg,.png,.webp,.heic,.heif" onChange={(e) => setThumbnailFile(e.target.files?.[0] ?? null)} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update video' : 'Publish video'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>Cancel</button>
      </div>
    </form>
  );
}
