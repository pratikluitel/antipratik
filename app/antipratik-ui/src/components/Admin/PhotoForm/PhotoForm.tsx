'use client';

import { useState, FormEvent } from 'react';
import type { PhotoPost } from '@/lib/types';
import { createPhotoPost, updatePhotoPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';
import styles from './PhotoForm.module.css';

interface PhotoFormProps {
  token: string;
  initial?: PhotoPost;
  onSuccess: () => void;
  onCancel: () => void;
}

interface PhotoEntry {
  file: File;
  alt: string;
  caption: string;
}

export default function PhotoForm({ token, initial, onSuccess, onCancel }: PhotoFormProps) {
  const [photos, setPhotos] = useState<PhotoEntry[]>([]);
  const [location, setLocation] = useState(initial?.location ?? '');
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  function handleFilesChange(files: FileList | null) {
    if (!files) return;
    const entries: PhotoEntry[] = Array.from(files).map((file) => ({ file, alt: '', caption: '' }));
    setPhotos((prev) => [...prev, ...entries]);
  }

  function updateEntry(index: number, field: 'alt' | 'caption', value: string) {
    setPhotos((prev) => prev.map((p, i) => i === index ? { ...p, [field]: value } : p));
  }

  function removeEntry(index: number) {
    setPhotos((prev) => prev.filter((_, i) => i !== index));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const fd = new FormData();
      if (initial) {
        if (location !== (initial.location ?? '')) fd.append('location', location);
        tags.forEach((t) => fd.append('tags[]', t));
        await updatePhotoPost(initial.id, fd, token);
      } else {
        if (photos.length === 0) { setError('At least one image is required.'); setLoading(false); return; }
        const missingAlt = photos.findIndex((p) => !p.alt.trim());
        if (missingAlt !== -1) { setError(`Alt text is required for image ${missingAlt + 1}.`); setLoading(false); return; }
        photos.forEach((p) => {
          fd.append('images[]', p.file);
          fd.append('alt[]', p.alt);
          fd.append('caption[]', p.caption);
        });
        if (location) fd.append('location', location);
        tags.forEach((t) => fd.append('tags[]', t));
        await createPhotoPost(fd, token);
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

      {initial ? (
        <div className={f.field}>
          <span className={f.label}>Images <span className={f.immutableNote}>(cannot be changed after creation)</span></span>
          <p className={f.immutableNote}>{initial.images.length} image(s) in this post.</p>
        </div>
      ) : (
        <div className={f.field}>
          <label className={`${f.label} ${f.required}`} htmlFor="photo-files">Images</label>
          <input id="photo-files" className={f.fileInput} type="file" accept=".jpg,.jpeg,.png,.webp" multiple onChange={(e) => handleFilesChange(e.target.files)} disabled={loading} />
        </div>
      )}

      {!initial && photos.length > 0 && (
        <div className={styles.photoList}>
          {photos.map((p, i) => (
            <div key={i} className={styles.photoEntry}>
              <p className={styles.fileName}>{p.file.name}</p>
              <div className={f.field}>
                <label className={`${f.label} ${f.required}`}>Alt text</label>
                <input className={f.input} value={p.alt} onChange={(e) => updateEntry(i, 'alt', e.target.value)} disabled={loading} />
              </div>
              <div className={f.field}>
                <label className={f.label}>Caption</label>
                <input className={f.input} value={p.caption} onChange={(e) => updateEntry(i, 'caption', e.target.value)} disabled={loading} />
              </div>
              <button type="button" className={f.cancelBtn} onClick={() => removeEntry(i)} disabled={loading}>Remove</button>
            </div>
          ))}
        </div>
      )}

      <div className={f.field}>
        <label className={f.label} htmlFor="photo-location">Location</label>
        <input id="photo-location" className={f.input} value={location} onChange={(e) => setLocation(e.target.value)} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update post' : 'Publish photos'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>Cancel</button>
      </div>
    </form>
  );
}
