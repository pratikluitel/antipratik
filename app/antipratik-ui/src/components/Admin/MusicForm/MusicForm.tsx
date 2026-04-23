'use client';

import { useState, FormEvent } from 'react';
import type { MusicPost } from '@/lib/types';
import { createMusicPost, updateMusicPost } from '@/lib/api';
import TagInput from '../TagInput';
import f from '../adminForm.module.css';

interface MusicFormProps {
  token: string;
  initial?: MusicPost;
  onSuccess: () => void;
  onCancel: () => void;
}

function readAudioDuration(file: File): Promise<number> {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(file);
    const audio = new Audio();
    audio.addEventListener('loadedmetadata', () => {
      URL.revokeObjectURL(url);
      resolve(Math.round(audio.duration));
    });
    audio.addEventListener('error', () => {
      URL.revokeObjectURL(url);
      reject(new Error('Could not read audio duration.'));
    });
    audio.src = url;
  });
}

export default function MusicForm({ token, initial, onSuccess, onCancel }: MusicFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [duration, setDuration] = useState('');
  const [album, setAlbum] = useState(initial?.album ?? '');
  const [tags, setTags] = useState<string[]>(initial?.tags ?? []);
  const [audioFile, setAudioFile] = useState<File | null>(null);
  const [albumArtFile, setAlbumArtFile] = useState<File | null>(null);
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
        if (album !== (initial.album ?? '')) fd.append('album', album);
        if (albumArtFile) fd.append('albumArtFile', albumArtFile);
        tags.forEach((t) => fd.append('tags[]', t));
        await updateMusicPost(initial.id, fd, token);
      } else {
        if (!audioFile) { setError('Audio file is required.'); setLoading(false); return; }
        if (!duration) { setError('Could not read duration from audio file.'); setLoading(false); return; }
        fd.append('title', title);
        fd.append('duration', duration);
        if (album) fd.append('album', album);
        fd.append('audioFile', audioFile);
        if (albumArtFile) fd.append('albumArtFile', albumArtFile);
        tags.forEach((t) => fd.append('tags[]', t));
        await createMusicPost(fd, token);
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
        <label className={`${f.label} ${f.required}`} htmlFor="music-title">Title</label>
        <input id="music-title" className={f.input} value={title} onChange={(e) => setTitle(e.target.value)} required disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={`${f.label} ${!initial ? f.required : ''}`} htmlFor="music-audio">
          Audio file {initial && <span className={f.immutableNote}>(cannot be changed after creation)</span>}
        </label>
        {initial ? (
          <p className={f.immutableNote}>Current: {initial.audioUrl}</p>
        ) : (
          <input
            id="music-audio"
            className={f.fileInput}
            type="file"
            accept=".mp3,.wav,.ogg,.m4a"
            onChange={async (e) => {
              const file = e.target.files?.[0] ?? null;
              setAudioFile(file);
              if (file) {
                try {
                  const d = await readAudioDuration(file);
                  setDuration(d.toString());
                } catch {
                  setDuration('');
                }
              } else {
                setDuration('');
              }
            }}
            disabled={loading}
          />
        )}
      </div>

      {!initial && (
        <div className={f.field}>
          <label className={`${f.label} ${f.required}`} htmlFor="music-duration">Duration (seconds)</label>
          <input
            id="music-duration"
            className={f.input}
            type="number"
            min="1"
            value={duration}
            readOnly
            disabled={loading || !audioFile}
            placeholder={audioFile ? undefined : 'Select audio file first'}
          />
        </div>
      )}

      <div className={f.field}>
        <label className={f.label} htmlFor="music-art">Album art</label>
        {initial && <p className={f.immutableNote}>Current: {initial.albumArt || 'none'}</p>}
        <input id="music-art" className={f.fileInput} type="file" accept=".jpg,.jpeg,.png,.webp,.heic,.heif" onChange={(e) => setAlbumArtFile(e.target.files?.[0] ?? null)} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label} htmlFor="music-album">Album</label>
        <input id="music-album" className={f.input} value={album} onChange={(e) => setAlbum(e.target.value)} disabled={loading} />
      </div>

      <div className={f.field}>
        <label className={f.label}>Tags</label>
        <TagInput tags={tags} onChange={setTags} disabled={loading} />
      </div>

      <div className={f.actions}>
        <button className={f.submitBtn} type="submit" disabled={loading || success}>
          {loading ? 'Saving…' : initial ? 'Update track' : 'Publish track'}
        </button>
        <button className={f.cancelBtn} type="button" onClick={onCancel} disabled={loading}>Cancel</button>
      </div>
    </form>
  );
}
