'use client';

import { useState, FormEvent, useRef, useEffect } from 'react';
import type { PhotoPost, PhotoImage } from '@/lib/types';
import { createPhotoPost, updatePhotoPost, addPhotoImage, updatePhotoImage, deletePhotoImage } from '@/lib/api';
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

  // Edit-mode: existing images state
  const [existingImages, setExistingImages] = useState<PhotoImage[]>(initial?.images ?? []);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editCaption, setEditCaption] = useState('');
  const [editAlt, setEditAlt] = useState('');
  const [imageError, setImageError] = useState<string | null>(null);

  // Edit-mode: add a single new image
  const [newFile, setNewFile] = useState<File | null>(null);
  const [newFilePreview, setNewFilePreview] = useState<string | null>(null);
  const [newAlt, setNewAlt] = useState('');
  const [newCaption, setNewCaption] = useState('');
  const [addingImage, setAddingImage] = useState(false);
  const addFileInputRef = useRef<HTMLInputElement>(null);

  // Cleanup object URLs to prevent memory leaks
  useEffect(() => {
    return () => {
      photos.forEach((p) => URL.revokeObjectURL(p.previewUrl));
    };
  }, [photos]);

  useEffect(() => {
    return () => {
      if (newFilePreview) URL.revokeObjectURL(newFilePreview);
    };
  }, [newFilePreview]);

  function handleFilesChange(files: FileList | null) {
    if (!files) return;
    const entries: PhotoEntry[] = Array.from(files).map((file) => ({
      file,
      alt: '',
      caption: '',
      previewUrl: URL.createObjectURL(file),
    }));
    setPhotos((prev) => [...prev, ...entries]);
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

  // Edit-mode: start editing an existing image's alt + caption
  function startEdit(img: PhotoImage) {
    setEditingId(img.id!);
    setEditAlt(img.alt);
    setEditCaption(img.caption ?? '');
    setImageError(null);
  }

  // Edit-mode: save alt + caption for an existing image
  async function saveImageMeta(img: PhotoImage) {
    if (!initial) return;
    if (!editAlt.trim()) { setImageError('Alt text is required.'); return; }
    setImageError(null);
    try {
      const updated = await updatePhotoImage(initial.id, img.id!, { alt: editAlt, caption: editCaption }, token);
      setExistingImages((prev) => prev.map((i) => i.id === img.id ? { ...i, alt: updated.alt, caption: updated.caption } : i));
      setEditingId(null);
    } catch (err) {
      setImageError(err instanceof Error ? err.message : 'Failed to save.');
    }
  }

  // Edit-mode: delete an existing image
  async function handleDeleteImage(img: PhotoImage) {
    if (!initial) return;
    setImageError(null);
    try {
      await deletePhotoImage(initial.id, img.id!, token);
      setExistingImages((prev) => prev.filter((i) => i.id !== img.id));
    } catch (err) {
      setImageError(err instanceof Error ? err.message : 'Failed to delete image.');
    }
  }

  // Edit-mode: select new file to add
  function handleNewFileChange(files: FileList | null) {
    if (!files || files.length === 0) return;
    if (newFilePreview) URL.revokeObjectURL(newFilePreview);
    setNewFile(files[0]);
    setNewFilePreview(URL.createObjectURL(files[0]));
    if (addFileInputRef.current) addFileInputRef.current.value = '';
  }

  // Edit-mode: upload new single image
  async function handleAddImage() {
    if (!initial || !newFile) return;
    if (!newAlt.trim()) { setImageError('Alt text is required.'); return; }
    setImageError(null);
    setAddingImage(true);
    try {
      const fd = new FormData();
      fd.append('image', newFile);
      fd.append('alt', newAlt);
      if (newCaption.trim()) fd.append('caption', newCaption);
      const img = await addPhotoImage(initial.id, fd, token);
      setExistingImages((prev) => [...prev, img]);
      setNewFile(null);
      if (newFilePreview) URL.revokeObjectURL(newFilePreview);
      setNewFilePreview(null);
      setNewAlt('');
      setNewCaption('');
    } catch (err) {
      setImageError(err instanceof Error ? err.message : 'Failed to add image.');
    } finally {
      setAddingImage(false);
    }
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

        {/* Edit mode: interactive existing image list */}
        {initial && (
          <div className={styles.existingImages}>
            {imageError && <p className={f.error} role="alert">{imageError}</p>}
            {existingImages.map((img) => (
              <div key={img.id} className={styles.existingImageRow}>
                {img.thumbnailSmallUrl && (
                  <div className={styles.thumbnailItem}>
                    <Image
                      src={img.thumbnailSmallUrl}
                      alt={img.alt}
                      width={80}
                      height={80}
                      className={styles.smallThumb}
                    />
                  </div>
                )}
                <div className={styles.existingImageMeta}>
                  {editingId === img.id ? (
                    <>
                      <div className={f.field}>
                        <label className={`${f.label} ${f.required}`}>Alt text</label>
                        <input
                          className={f.input}
                          value={editAlt}
                          onChange={(e) => setEditAlt(e.target.value)}
                          autoFocus
                        />
                      </div>
                      <div className={f.field}>
                        <label className={f.label}>Caption</label>
                        <input
                          className={f.input}
                          value={editCaption}
                          onChange={(e) => setEditCaption(e.target.value)}
                          placeholder="Optional"
                        />
                      </div>
                      <div className={styles.imageRowActions}>
                        <button type="button" className={f.submitBtn} onClick={() => saveImageMeta(img)}>Save</button>
                        <button type="button" className={f.cancelBtn} onClick={() => setEditingId(null)}>Cancel</button>
                      </div>
                    </>
                  ) : (
                    <>
                      <p className={styles.captionDisplay}><strong>Alt:</strong> {img.alt}</p>
                      {img.caption
                        ? <p className={styles.captionEmpty}><strong>Caption:</strong> {img.caption}</p>
                        : <p className={styles.captionEmpty}>No caption</p>
                      }
                      <div className={styles.imageRowActions}>
                        <button type="button" className={f.cancelBtn} onClick={() => startEdit(img)}>Edit</button>
                        <button
                          type="button"
                          className={f.cancelBtn}
                          onClick={() => handleDeleteImage(img)}
                          disabled={existingImages.length <= 1}
                          title={existingImages.length <= 1 ? 'Cannot delete the only image' : 'Delete image'}
                        >
                          Delete
                        </button>
                      </div>
                    </>
                  )}
                </div>
              </div>
            ))}

            {/* Add a new image in edit mode */}
            <div className={styles.addImageSection}>
              <p className={f.label}>Add a photo</p>
              <input
                ref={addFileInputRef}
                type="file"
                accept=".jpg,.jpeg,.png,.webp,.heic,.heif"
                className={styles.hiddenInput}
                onChange={(e) => handleNewFileChange(e.target.files)}
                disabled={addingImage}
              />
              {newFilePreview && newFile ? (
                <>
                  <div className={styles.previewContainer}>
                    <div className={styles.previewWrapper}>
                      <Image src={newFilePreview} alt="Preview" fill sizes="120px" className={styles.previewImage} />
                    </div>
                    <p className={styles.fileName}>{newFile.name}</p>
                  </div>
                  <div className={f.field}>
                    <label className={`${f.label} ${f.required}`}>Alt text</label>
                    <input className={f.input} value={newAlt} onChange={(e) => setNewAlt(e.target.value)} disabled={addingImage} />
                  </div>
                  <div className={f.field}>
                    <label className={f.label}>Caption</label>
                    <input className={f.input} value={newCaption} onChange={(e) => setNewCaption(e.target.value)} disabled={addingImage} />
                  </div>
                  <div className={styles.imageRowActions}>
                    <button type="button" className={f.submitBtn} onClick={handleAddImage} disabled={addingImage}>
                      {addingImage ? 'Uploading…' : 'Add photo'}
                    </button>
                    <button type="button" className={f.cancelBtn} onClick={() => { setNewFile(null); if (newFilePreview) URL.revokeObjectURL(newFilePreview); setNewFilePreview(null); setNewAlt(''); setNewCaption(''); }} disabled={addingImage}>
                      Cancel
                    </button>
                  </div>
                </>
              ) : (
                <button
                  type="button"
                  className={f.cancelBtn}
                  onClick={() => addFileInputRef.current?.click()}
                  disabled={addingImage}
                >
                  Select photo to add
                </button>
              )}
            </div>
          </div>
        )}

        {/* Create mode: multi-file upload */}
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
              className={f.cancelBtn}
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
