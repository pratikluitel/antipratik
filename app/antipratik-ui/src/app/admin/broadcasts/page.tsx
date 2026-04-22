'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { useRouter } from 'next/navigation';
import {
  getPosts,
  getBroadcasts,
  getSubscribers,
  createBroadcast,
  updateBroadcast,
  dispatchBroadcast,
  getBroadcastSendDetails,
} from '@/lib/api';
import type {
  Post,
  BroadcastSummary,
  BroadcastSendDetail,
  SubscriberSummary,
  CreateBroadcastInput,
} from '@/lib/types';
import AdminLayout from '@/components/Admin/AdminLayout';
import styles from './broadcasts.module.css';
import formStyles from '@/components/Admin/adminForm.module.css';

type SubTab = 'broadcasts' | 'subscribers';

type FormMode =
  | { kind: 'none' }
  | { kind: 'create' }
  | { kind: 'edit'; broadcast: BroadcastSummary };

const TYPE_OPTIONS = ['email'] as const;

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

function postLabel(post: Post): string {
  switch (post.type) {
    case 'essay': return post.title;
    case 'short': return post.body.slice(0, 60) + (post.body.length > 60 ? '…' : '');
    case 'music': return post.title;
    case 'photo': return post.images[0]?.alt ?? `Photo — ${formatDate(post.createdAt)}`;
    case 'video': return post.title;
    case 'link': return post.title;
  }
}

function postTypeLabel(type: Post['type']): string {
  return type.charAt(0).toUpperCase() + type.slice(1);
}

// ─── Post selector dropdown ───────────────────────────────────────────────────

interface PostSelectorProps {
  posts: Post[];
  selectedIDs: string[];
  onChange: (ids: string[]) => void;
}

function PostSelector({ posts, selectedIDs, onChange }: PostSelectorProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function onOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener('mousedown', onOutside);
    return () => document.removeEventListener('mousedown', onOutside);
  }, []);

  const filtered = posts.filter((p) => {
    const label = postLabel(p).toLowerCase();
    return query === '' || label.includes(query.toLowerCase()) || p.type.includes(query.toLowerCase());
  });

  function toggle(id: string) {
    if (selectedIDs.includes(id)) {
      onChange(selectedIDs.filter((x) => x !== id));
    } else {
      onChange([...selectedIDs, id]);
    }
  }

  function removeSelected(id: string) {
    onChange(selectedIDs.filter((x) => x !== id));
  }

  function moveUp(index: number) {
    if (index === 0) return;
    const next = [...selectedIDs];
    [next[index - 1], next[index]] = [next[index], next[index - 1]];
    onChange(next);
  }

  function moveDown(index: number) {
    if (index === selectedIDs.length - 1) return;
    const next = [...selectedIDs];
    [next[index], next[index + 1]] = [next[index + 1], next[index]];
    onChange(next);
  }

  const selectedPosts = selectedIDs
    .map((id) => posts.find((p) => p.id === id))
    .filter((p): p is Post => p !== undefined);

  return (
    <div className={styles.postSelector} ref={ref}>
      <button
        type="button"
        className={styles.selectorTrigger}
        onClick={() => setOpen((o) => !o)}
      >
        {selectedIDs.length === 0
          ? 'Select posts…'
          : `${selectedIDs.length} post${selectedIDs.length !== 1 ? 's' : ''} selected`}
        <span className={styles.selectorChevron}>{open ? '▴' : '▾'}</span>
      </button>

      {open && (
        <div className={styles.selectorDropdown}>
          <input
            className={styles.selectorSearch}
            type="text"
            placeholder="Search posts…"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            autoFocus
          />
          <div className={styles.selectorList}>
            {filtered.length === 0 && (
              <div className={styles.selectorEmpty}>No posts match</div>
            )}
            {filtered.map((post) => {
              const selected = selectedIDs.includes(post.id);
              return (
                <button
                  key={post.id}
                  type="button"
                  className={`${styles.selectorOption} ${selected ? styles.selectorOptionSelected : ''}`}
                  onClick={() => toggle(post.id)}
                >
                  <span className={`${styles.postTypeBadge} ${styles[`postType_${post.type}`]}`}>
                    {postTypeLabel(post.type)}
                  </span>
                  <span className={styles.postOptionLabel}>{postLabel(post)}</span>
                  <span className={styles.postOptionDate}>{formatDate(post.createdAt)}</span>
                  {selected && <span className={styles.checkmark}>✓</span>}
                </button>
              );
            })}
          </div>
        </div>
      )}

      {selectedPosts.length > 0 && (
        <ol className={styles.selectedList}>
          {selectedPosts.map((post, i) => (
            <li key={post.id} className={styles.selectedItem}>
              <div className={styles.selectedItemInfo}>
                <span className={`${styles.postTypeBadge} ${styles[`postType_${post.type}`]}`}>
                  {postTypeLabel(post.type)}
                </span>
                <span className={styles.selectedItemLabel}>{postLabel(post)}</span>
                <span className={styles.selectedItemDate}>{formatDate(post.createdAt)}</span>
              </div>
              <div className={styles.selectedItemActions}>
                <button type="button" className={styles.orderBtn} onClick={() => moveUp(i)} disabled={i === 0} aria-label="Move up">↑</button>
                <button type="button" className={styles.orderBtn} onClick={() => moveDown(i)} disabled={i === selectedIDs.length - 1} aria-label="Move down">↓</button>
                <button type="button" className={styles.removeBtn} onClick={() => removeSelected(post.id)} aria-label="Remove">✕</button>
              </div>
            </li>
          ))}
        </ol>
      )}
    </div>
  );
}

// ─── Broadcast form ───────────────────────────────────────────────────────────

interface BroadcastFormProps {
  posts: Post[];
  mode: FormMode;
  token: string;
  onSaved: (preview: { id: number; html: string }) => void;
  onCancel: () => void;
}

function BroadcastForm({ posts, mode, token, onSaved, onCancel }: BroadcastFormProps) {
  const existing = mode.kind === 'edit' ? mode.broadcast : null;
  const [bType, setBType] = useState(existing?.type ?? 'email');
  const [title, setTitle] = useState(existing?.title ?? '');
  const [caption, setCaption] = useState(existing?.caption ?? '');
  const [selectedPostIDs, setSelectedPostIDs] = useState<string[]>(existing?.postIDs ?? ([] as string[]));
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [preview, setPreview] = useState<{ id: number; html: string } | null>(null);
  const [showPreview, setShowPreview] = useState(false);

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setError('');
    setSaving(true);
    try {
      let result: { id: number; html: string };
      if (mode.kind === 'edit') {
        result = await updateBroadcast(mode.broadcast.id, { title, caption, postIDs: selectedPostIDs }, token);
      } else {
        const input: CreateBroadcastInput = {
          type: bType,
          title,
          data: { caption, postIDs: selectedPostIDs },
        };
        result = await createBroadcast(input, token);
      }
      setPreview(result);
      onSaved(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save broadcast');
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className={styles.formSection}>
      <h3 className={styles.formTitle}>
        {mode.kind === 'edit' ? 'Edit Broadcast' : 'New Broadcast'}
      </h3>
      <form className={formStyles.form} onSubmit={handleSubmit}>
        {mode.kind === 'create' && (
          <div className={formStyles.field}>
            <label className={formStyles.label}>Type</label>
            <select
              className={formStyles.input}
              value={bType}
              onChange={(e) => setBType(e.target.value)}
              required
            >
              {TYPE_OPTIONS.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>
        )}
        <div className={formStyles.field}>
          <label className={formStyles.label}><span className={formStyles.required}>Title</span></label>
          <input
            className={formStyles.input}
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Email subject / broadcast title"
            required
          />
        </div>
        <div className={formStyles.field}>
          <label className={formStyles.label}>Caption</label>
          <input
            className={formStyles.input}
            type="text"
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            placeholder="Optional caption / preview text"
          />
        </div>
        <div className={formStyles.field}>
          <label className={formStyles.label}>Posts</label>
          <PostSelector posts={posts} selectedIDs={selectedPostIDs} onChange={setSelectedPostIDs} />
        </div>

        {error && <p className={formStyles.error}>{error}</p>}

        <div className={formStyles.actions}>
          <button type="submit" className={formStyles.submitBtn} disabled={saving}>
            {saving ? 'Saving…' : mode.kind === 'edit' ? 'Update' : 'Create'}
          </button>
          <button type="button" className={formStyles.cancelBtn} onClick={onCancel}>
            Cancel
          </button>
        </div>
      </form>

      {preview && (
        <div className={styles.previewSection}>
          <button
            type="button"
            className={styles.previewToggle}
            onClick={() => setShowPreview((s) => !s)}
          >
            {showPreview ? 'Hide preview' : 'Show email preview'} {showPreview ? '▴' : '▾'}
          </button>
          {showPreview && (
            <iframe
              className={styles.previewFrame}
              srcDoc={preview.html}
              title="Email preview"
              sandbox="allow-same-origin"
            />
          )}
        </div>
      )}
    </div>
  );
}

// ─── Broadcast row ────────────────────────────────────────────────────────────

interface BroadcastRowProps {
  broadcast: BroadcastSummary;
  token: string;
  onEdit: () => void;
  onDispatched: (id: number, buffered: number) => void;
}

function BroadcastRow({ broadcast, token, onEdit, onDispatched }: BroadcastRowProps) {
  const [dispatching, setDispatching] = useState(false);
  const [error, setError] = useState('');
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [sendsOpen, setSendsOpen] = useState(false);
  const [sends, setSends] = useState<BroadcastSendDetail[] | null>(null);
  const [sendsLoading, setSendsLoading] = useState(false);
  const [sendsError, setSendsError] = useState('');

  const isDispatched = broadcast.buffered > 0 || broadcast.success > 0 || broadcast.failed > 0;
  const postCount = (broadcast.postIDs ?? []).length;

  async function handleDispatch() {
    setError('');
    setDispatching(true);
    setConfirmOpen(false);
    try {
      const result = await dispatchBroadcast(broadcast.id, token);
      onDispatched(broadcast.id, result.buffered_count);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Dispatch failed');
    } finally {
      setDispatching(false);
    }
  }

  async function toggleSends() {
    if (sendsOpen) { setSendsOpen(false); return; }
    setSendsOpen(true);
    if (sends !== null) return;
    setSendsLoading(true);
    setSendsError('');
    try {
      const data = await getBroadcastSendDetails(broadcast.id, token);
      setSends(data);
    } catch (err) {
      setSendsError(err instanceof Error ? err.message : 'Failed to load sends');
    } finally {
      setSendsLoading(false);
    }
  }

  return (
    <div className={styles.broadcastRow}>
      <div className={styles.broadcastRowMain}>
        <div className={styles.broadcastRowMeta}>
          <span className={styles.broadcastId}>#{broadcast.id}</span>
          <span className={`${styles.typeBadge} ${styles[`typeBadge_${broadcast.type}`]}`}>
            {broadcast.type}
          </span>
        </div>
        <div className={styles.broadcastRowTitle}>{broadcast.title}</div>
        {broadcast.caption && (
          <div className={styles.broadcastRowCaption}>{broadcast.caption}</div>
        )}
        <div className={styles.broadcastRowStats}>
          {isDispatched ? (
            <>
              <span className={styles.statBuffered}>{broadcast.buffered} buffered</span>
              <span className={styles.statSuccess}>{broadcast.success} sent</span>
              <span className={styles.statFailed}>{broadcast.failed} failed</span>
            </>
          ) : (
            <span className={styles.statPending}>Not dispatched</span>
          )}
          <span className={styles.statPosts}>{postCount} post{postCount !== 1 ? 's' : ''}</span>
        </div>
      </div>
      <div className={styles.broadcastRowActions}>
        <button type="button" className={styles.actionBtn} onClick={onEdit}>
          Edit
        </button>
        {isDispatched && (
          <button
            type="button"
            className={`${styles.actionBtn} ${sendsOpen ? styles.actionBtnActive : ''}`}
            onClick={toggleSends}
          >
            {sendsOpen ? 'Hide sends' : 'View sends'}
          </button>
        )}
        {!confirmOpen ? (
          <button
            type="button"
            className={`${styles.actionBtn} ${styles.actionBtnDispatch}`}
            onClick={() => setConfirmOpen(true)}
            disabled={dispatching}
          >
            {dispatching ? 'Sending…' : 'Dispatch'}
          </button>
        ) : (
          <div className={styles.confirmInline}>
            <span className={styles.confirmText}>Send to all subscribers?</span>
            <button type="button" className={`${styles.actionBtn} ${styles.actionBtnDanger}`} onClick={handleDispatch}>
              Confirm
            </button>
            <button type="button" className={styles.actionBtn} onClick={() => setConfirmOpen(false)}>
              Cancel
            </button>
          </div>
        )}
      </div>
      {error && <p className={styles.rowError}>{error}</p>}

      {sendsOpen && (
        <div className={styles.sendsPanel}>
          {sendsLoading && <p className={styles.sendsPanelMsg}>Loading…</p>}
          {sendsError && <p className={styles.sendsPanelMsg}>{sendsError}</p>}
          {sends !== null && sends.length === 0 && (
            <p className={styles.sendsPanelMsg}>No send records yet.</p>
          )}
          {sends !== null && sends.length > 0 && (
            <div className={styles.sendsTable}>
              <div className={styles.sendsHeader}>
                <span>Address</span>
                <span>Status</span>
                <span>Scheduled</span>
                <span>Sent</span>
                <span>Message</span>
              </div>
              {sends.map((s, i) => (
                <div key={i} className={styles.sendsRow}>
                  <span className={styles.sendsAddress}>{s.address}</span>
                  <span>
                    {s.status === 'SUCCESS' && <span className={styles.badgeSuccess}>Sent</span>}
                    {s.status === 'FAILED' && <span className={styles.badgeFailed}>Failed</span>}
                    {s.status === 'BUFFERED' && <span className={styles.badgeBuffered}>Queued</span>}
                  </span>
                  <span className={styles.sendsDate}>{formatDate(s.scheduledAt)}</span>
                  <span className={styles.sendsDate}>{s.sentAt ? formatDate(s.sentAt) : '—'}</span>
                  <span className={styles.sendsMessage}>{s.message ?? '—'}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

// ─── Main page ────────────────────────────────────────────────────────────────

export default function BroadcastsPage() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [subTab, setSubTab] = useState<SubTab>('broadcasts');
  const [posts, setPosts] = useState<Post[]>([]);
  const [broadcasts, setBroadcasts] = useState<BroadcastSummary[]>([]);
  const [subscribers, setSubscribers] = useState<SubscriberSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [formMode, setFormMode] = useState<FormMode>({ kind: 'none' });
  const [typeFilter, setTypeFilter] = useState('email');
  const [confirmedFilter, setConfirmedFilter] = useState<'all' | 'confirmed' | 'unconfirmed'>('all');

  useEffect(() => {
    const t = localStorage.getItem('admin_token');
    if (!t) { router.replace('/admin/login'); return; }
    setToken(t);
  }, [router]);

  const loadData = useCallback(async (t: string) => {
    setLoading(true);
    try {
      const [allPosts, allBroadcasts, allSubscribers] = await Promise.all([
        getPosts(),
        getBroadcasts(typeFilter, t),
        getSubscribers(typeFilter, t),
      ]);
      setPosts(allPosts);
      setBroadcasts(allBroadcasts);
      setSubscribers(allSubscribers);
    } catch {
      // silently ignore — may not be available in dev
    } finally {
      setLoading(false);
    }
  }, [typeFilter]);

  useEffect(() => {
    if (token) loadData(token);
  }, [token, loadData]);

  function handleSaved(_preview: { id: number; html: string }) {
    if (token) loadData(token);
    setFormMode({ kind: 'none' });
  }

  function handleDispatched(id: number, buffered: number) {
    setBroadcasts((prev) =>
      prev.map((b) => b.id === id ? { ...b, buffered } : b)
    );
  }

  const filteredSubscribers = subscribers.filter((s) => {
    if (confirmedFilter === 'confirmed') return s.confirmed;
    if (confirmedFilter === 'unconfirmed') return !s.confirmed;
    return true;
  });

  if (!token) return null;

  return (
    <AdminLayout>
      <div className={styles.page}>
        <div className={styles.subTabs}>
          <button
            type="button"
            className={`${styles.subTab} ${subTab === 'broadcasts' ? styles.subTabActive : ''}`}
            onClick={() => setSubTab('broadcasts')}
          >
            Broadcasts
          </button>
          <button
            type="button"
            className={`${styles.subTab} ${subTab === 'subscribers' ? styles.subTabActive : ''}`}
            onClick={() => setSubTab('subscribers')}
          >
            Subscribers
            {subscribers.length > 0 && (
              <span className={styles.subTabCount}>{subscribers.length}</span>
            )}
          </button>
        </div>

        {/* ─── Broadcasts tab ─────────────────────────────────── */}
        {subTab === 'broadcasts' && (
          <div className={styles.tabContent}>
            <div className={styles.tabToolbar}>
              <div className={styles.filterRow}>
                <label className={styles.filterLabel}>Type</label>
                {TYPE_OPTIONS.map((t) => (
                  <button
                    key={t}
                    type="button"
                    className={`${styles.filterPill} ${typeFilter === t ? styles.filterPillActive : ''}`}
                    onClick={() => setTypeFilter(t)}
                  >
                    {t}
                  </button>
                ))}
              </div>
              {formMode.kind === 'none' && (
                <button
                  type="button"
                  className={styles.newBtn}
                  onClick={() => setFormMode({ kind: 'create' })}
                >
                  + New Broadcast
                </button>
              )}
            </div>

            {formMode.kind !== 'none' && (
              <BroadcastForm
                posts={posts}
                mode={formMode}
                token={token}
                onSaved={handleSaved}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}

            {loading && <p className={styles.loading}>Loading…</p>}

            {!loading && broadcasts.length === 0 && (
              <p className={styles.empty}>No broadcasts yet.</p>
            )}

            <div className={styles.broadcastList}>
              {broadcasts.map((b) => (
                <BroadcastRow
                  key={b.id}
                  broadcast={b}
                  token={token}
                  onEdit={() => setFormMode({ kind: 'edit', broadcast: b })}
                  onDispatched={handleDispatched}
                />
              ))}
            </div>
          </div>
        )}

        {/* ─── Subscribers tab ────────────────────────────────── */}
        {subTab === 'subscribers' && (
          <div className={styles.tabContent}>
            <div className={styles.tabToolbar}>
              <div className={styles.filterRow}>
                <label className={styles.filterLabel}>Type</label>
                {TYPE_OPTIONS.map((t) => (
                  <button
                    key={t}
                    type="button"
                    className={`${styles.filterPill} ${typeFilter === t ? styles.filterPillActive : ''}`}
                    onClick={() => setTypeFilter(t)}
                  >
                    {t}
                  </button>
                ))}
              </div>
              <div className={styles.filterRow}>
                <label className={styles.filterLabel}>Status</label>
                {(['all', 'confirmed', 'unconfirmed'] as const).map((f) => (
                  <button
                    key={f}
                    type="button"
                    className={`${styles.filterPill} ${confirmedFilter === f ? styles.filterPillActive : ''}`}
                    onClick={() => setConfirmedFilter(f)}
                  >
                    {f.charAt(0).toUpperCase() + f.slice(1)}
                  </button>
                ))}
              </div>
            </div>

            {loading && <p className={styles.loading}>Loading…</p>}

            {!loading && filteredSubscribers.length === 0 && (
              <p className={styles.empty}>No subscribers match.</p>
            )}

            {filteredSubscribers.length > 0 && (
              <div className={styles.subscriberTable}>
                <div className={styles.subscriberHeader}>
                  <span>Address</span>
                  <span>Type</span>
                  <span>Status</span>
                  <span>Signed up</span>
                  <span>Confirmed</span>
                </div>
                {filteredSubscribers.map((s) => (
                  <div key={s.address} className={styles.subscriberRow}>
                    <span className={styles.subscriberAddress}>{s.address}</span>
                    <span className={styles.subscriberType}>{s.type}</span>
                    <span>
                      {s.unsubscribedAt ? (
                        <span className={styles.badgeUnsubscribed}>Unsubscribed</span>
                      ) : s.confirmed ? (
                        <span className={styles.badgeConfirmed}>Confirmed</span>
                      ) : (
                        <span className={styles.badgePending}>Pending</span>
                      )}
                    </span>
                    <span className={styles.subscriberDate}>{formatDate(s.createdAt)}</span>
                    <span className={styles.subscriberDate}>
                      {s.confirmedAt ? formatDate(s.confirmedAt) : '—'}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </AdminLayout>
  );
}
