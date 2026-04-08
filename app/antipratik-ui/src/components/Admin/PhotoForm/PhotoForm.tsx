'use client';

import { useState, FormEvent, useRef, useEffect } from 'react';
import type { PhotoPost } from '@/lib/types';
import { createPhotoPost, updatePhotoPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';
import styles from './PhotoForm.module.css';
import Image from 'next/image';

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
  previewUrl: string;
}

export default function PhotoForm({ token, initial, onSuccess, onCancel }: PhotoFormProps) {
  const [photos, setPhotos] = useState<PhotoEntry[]>([]);
  const [location, setLocation] = useState(initial?.location ?? '');
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Cleanup object URLs to prevent memory leaks
  useEffect(() => {
    return () => {
      photos.forEach((p) => URL.revokeObjectURL(p.previewUrl));
    };
  }, [photos]);

  function handleFilesChange(files: FileList | null) {
    if (!files) return;
    const entries: PhotoEntry[] = Array.from(files).map((file) => ({
      file,
      alt: '',
      caption: '',
      previewUrl: URL.createObjectURL(file),
    }));
    setPhotos((prev) => [...prev, ...entries]);
    
    // Reset file input so same file can be selected again if needed
    if (fileInputRef.current) fileInputRef.current.value = '';
  }

  function updateEntry(index: number, field: 'alt' | 'caption', value: string) {
    setPhotos((prev) => prev.map((p, i) => i === index ? { ...p, [field]: value } : p));
  }

  function removeEntry(index: number) {
    setPhotos((prev) => {
      const removed = prev[index];
      if (removed) URL.revokeObjectURL(removed.previewUrl);
      return prev.filter((_, i) => i !== index);
    });
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

      <div className={f.field}>
        <label className={f.label}>Images</label>
        {initial?.images && initial.images.length > 0 && (
          <div className={styles.existingImages}>
            <p className={f.immutableNote}>Existing photos ({initial.images.length}):</p>
            <div className={styles.thumbnailGrid}>
              {initial.images.map((img, i) => (
                <div key={i} className={styles.thumbnailItem}>
                  {img.thumbnailSmallUrl && (
                    <Image 
                      src={img.thumbnailSmallUrl} 
                      alt={img.alt} 
                      width={80} 
                      height={80} 
                      className={styles.smallThumb} 
                    />
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
        
        {!initial && (
          <>
            <input 
              ref={fileInputRef}
              id="photo-files" 
              className={styles.hiddenInput} 
              type="file" 
              accept=".jpg,.jpeg,.png,.webp,.heic,.heif" 
              multiple 
              onChange={(e) => handleFilesChange(e.target.files)} 
              disabled={loading} 
            />
            <button 
              type="button" 
              className={f.secondaryBtn} 
              onClick={() => fileInputRef.current?.click()}
              disabled={loading}
            >
              {photos.length > 0 ? 'Add more photos' : 'Select photos'}
            </button>
          </>
        )}
      </div>

      {photos.length > 0 && (
        <div className={styles.photoList}>
          {photos.map((p, i) => (
            <div key={i} className={styles.photoEntry}>
              <div className={styles.previewContainer}>
                <div className={styles.previewWrapper}>
                  <Image 
                    src={p.previewUrl} 
                    alt="Preview" 
                    fill 
                    sizes="120px"
                    className={styles.previewImage}
                  />
                </div>
                <p className={styles.fileName}>{p.file.name}</p>
              </div>
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
