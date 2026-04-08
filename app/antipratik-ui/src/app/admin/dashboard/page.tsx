'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { getPosts, getLinks, deletePost, deleteExternalLink } from '@/lib/api';
import type { Post, EssayPost, ShortPost, MusicPost, PhotoPost, VideoPost, LinkPost, ExternalLink } from '@/lib/types';
import AdminLayout from '@/components/Admin/AdminLayout';
import EssayForm from '@/components/Admin/EssayForm';
import ShortPostForm from '@/components/Admin/ShortPostForm';
import MusicForm from '@/components/Admin/MusicForm';
import PhotoForm from '@/components/Admin/PhotoForm';
import VideoForm from '@/components/Admin/VideoForm';
import LinkForm from '@/components/Admin/LinkForm';
import ExternalLinkForm from '@/components/Admin/ExternalLinkForm';
import styles from './page.module.css';

type TabId = 'essay' | 'short' | 'music' | 'photo' | 'video' | 'link' | 'externallink';

const TABS: { id: TabId; label: string }[] = [
  { id: 'essay', label: 'Essays' },
  { id: 'short', label: 'Short Posts' },
  { id: 'music', label: 'Music' },
  { id: 'photo', label: 'Photos' },
  { id: 'video', label: 'Videos' },
  { id: 'link', label: 'Links' },
  { id: 'externallink', label: 'External Links' },
];

type FormMode =
  | { kind: 'none' }
  | { kind: 'create' }
  | { kind: 'edit-post'; post: Post }
  | { kind: 'edit-link'; link: ExternalLink };

export default function AdminDashboardPage() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabId>('essay');
  const [posts, setPosts] = useState<Post[]>([]);
  const [links, setLinks] = useState<ExternalLink[]>([]);
  const [formMode, setFormMode] = useState<FormMode>({ kind: 'none' });
  const [loadingData, setLoadingData] = useState(false);

  useEffect(() => {
    const t = localStorage.getItem('admin_token');
    if (!t) { router.replace('/admin/login'); return; }
    setToken(t);
  }, [router]);

  const refreshData = useCallback(async () => {
    setLoadingData(true);
    try {
      const [allPosts, allLinks] = await Promise.all([getPosts(), getLinks()]);
      setPosts(allPosts);
      setLinks(allLinks);
    } catch {
      // silently ignore — data may not be available in dev
    } finally {
      setLoadingData(false);
    }
  }, []);

  useEffect(() => {
    if (token) refreshData();
  }, [token, refreshData]);

  function handleSuccess() {
    setFormMode({ kind: 'none' });
    refreshData();
  }

  async function handleDeletePost(id: string) {
    if (!token || !confirm('Delete this post? This cannot be undone.')) return;
    await deletePost(id, token);
    refreshData();
  }

  async function handleDeleteLink(id: string) {
    if (!token || !confirm('Delete this link? This cannot be undone.')) return;
    await deleteExternalLink(id, token);
    refreshData();
  }

  if (!token) return null;

  const tabPosts = activeTab !== 'externallink'
    ? posts.filter((p) => p.type === activeTab)
    : [];

  return (
    <AdminLayout>
      <div className={styles.tabBar} role="tablist">
        {TABS.map((tab) => (
          <button
            key={tab.id}
            role="tab"
            type="button"
            className={`${styles.tab} ${activeTab === tab.id ? styles.tabActive : ''}`}
            onClick={() => { setActiveTab(tab.id); setFormMode({ kind: 'none' }); }}
            aria-selected={activeTab === tab.id}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className={styles.panel}>
        {formMode.kind === 'none' && (
          <button
            className={styles.newBtn}
            type="button"
            onClick={() => setFormMode({ kind: 'create' })}
          >
            + New {TABS.find((t) => t.id === activeTab)?.label.replace(/s$/, '')}
          </button>
        )}

        {formMode.kind !== 'none' && (
          <div className={styles.formSection}>
            <h2 className={styles.formTitle}>
              {formMode.kind === 'create' ? 'New' : 'Edit'}{' '}
              {TABS.find((t) => t.id === activeTab)?.label.replace(/s$/, '')}
            </h2>
            {activeTab === 'essay' && (
              <EssayForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as EssayPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'short' && (
              <ShortPostForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as ShortPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'music' && (
              <MusicForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as MusicPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'photo' && (
              <PhotoForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as PhotoPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'video' && (
              <VideoForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as VideoPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'link' && (
              <LinkForm
                token={token}
                initial={formMode.kind === 'edit-post' ? formMode.post as LinkPost : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
            {activeTab === 'externallink' && (
              <ExternalLinkForm
                token={token}
                initial={formMode.kind === 'edit-link' ? formMode.link : undefined}
                onSuccess={handleSuccess}
                onCancel={() => setFormMode({ kind: 'none' })}
              />
            )}
          </div>
        )}

        {formMode.kind === 'none' && (
          <div className={styles.list}>
            {loadingData && <p className={styles.loading}>Loading…</p>}

            {activeTab !== 'externallink' && tabPosts.length === 0 && !loadingData && (
              <p className={styles.empty}>No {TABS.find((t) => t.id === activeTab)?.label.toLowerCase()} yet.</p>
            )}

            {activeTab !== 'externallink' && tabPosts.map((post) => (
              <div key={post.id} className={styles.row}>
                <div className={styles.rowInfo}>
                  <span className={styles.rowTitle}>
                    {'title' in post ? post.title : post.type === 'short' ? (post as ShortPost).body.slice(0, 60) + '…' : post.id}
                  </span>
                  <span className={styles.rowMeta}>{new Date(post.createdAt).toLocaleDateString()}</span>
                </div>
                <div className={styles.rowActions}>
                  <button
                    type="button"
                    className={styles.editBtn}
                    onClick={() => setFormMode({ kind: 'edit-post', post })}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className={styles.deleteBtn}
                    onClick={() => handleDeletePost(post.id)}
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}

            {activeTab === 'externallink' && links.length === 0 && !loadingData && (
              <p className={styles.empty}>No external links yet.</p>
            )}

            {activeTab === 'externallink' && links.map((link) => (
              <div key={link.id} className={styles.row}>
                <div className={styles.rowInfo}>
                  <span className={styles.rowTitle}>{link.title}</span>
                  <span className={styles.rowMeta}>{link.domain} · {link.category}{link.featured ? ' · featured' : ''}</span>
                </div>
                <div className={styles.rowActions}>
                  <button
                    type="button"
                    className={styles.editBtn}
                    onClick={() => setFormMode({ kind: 'edit-link', link })}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className={styles.deleteBtn}
                    onClick={() => handleDeleteLink(link.id)}
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </AdminLayout>
  );
}
