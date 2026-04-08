'use client';

import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import styles from './MarkdownEditor.module.css';

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
}

export default function MarkdownEditor({ value, onChange, disabled }: MarkdownEditorProps) {
  const [activeTab, setActiveTab] = useState<'write' | 'preview'>('write');

  return (
    <div className={styles.editor}>
      <div className={styles.tabs} role="tablist">
        <button
          role="tab"
          type="button"
          className={`${styles.tab} ${activeTab === 'write' ? styles.tabActive : ''}`}
          onClick={() => setActiveTab('write')}
          aria-selected={activeTab === 'write'}
        >
          Write
        </button>
        <button
          role="tab"
          type="button"
          className={`${styles.tab} ${activeTab === 'preview' ? styles.tabActive : ''}`}
          onClick={() => setActiveTab('preview')}
          aria-selected={activeTab === 'preview'}
        >
          Preview
        </button>
      </div>

      {activeTab === 'write' ? (
        <textarea
          className={styles.textarea}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="Write markdown here…"
          disabled={disabled}
          spellCheck
        />
      ) : (
        <div className={styles.preview}>
          {value.trim() ? (
            <ReactMarkdown>{value}</ReactMarkdown>
          ) : (
            <p className={styles.empty}>Nothing to preview yet.</p>
          )}
        </div>
      )}
    </div>
  );
}
