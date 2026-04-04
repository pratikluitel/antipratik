import type { ExternalLink } from '../../lib/types';
import styles from './ExternalLinksBlock.module.css';

interface Props {
  links: ExternalLink[];
  variant: 'page' | 'homepage';
}

const ICON: Record<ExternalLink['category'], string> = {
  music: '♪',
  writing: '✦',
  video: '▶',
  social: '◐',
};

const CATEGORY_LABEL: Record<ExternalLink['category'], string> = {
  music: 'Music',
  writing: 'Writing',
  video: 'Video',
  social: 'Social',
};

const CATEGORY_ORDER: ExternalLink['category'][] = ['music', 'writing', 'video', 'social'];

function LinkRow({ link, variant }: { link: ExternalLink; variant: 'page' | 'homepage' }) {
  return (
    <a
      href={link.url}
      target="_blank"
      rel="noopener noreferrer"
      className={`${styles.row} ${variant === 'homepage' ? styles.rowHomepage : ''}`}
      data-category={link.category}
    >
      <span className={styles.accent} data-category={link.category} aria-hidden="true" />
      <span
        className={`${styles.icon} ${variant === 'homepage' ? styles.iconHomepage : ''}`}
        data-category={link.category}
        aria-hidden="true"
      >
        {ICON[link.category]}
      </span>
      <span className={styles.text}>
        <span className={`${styles.title} ${variant === 'homepage' ? styles.titleHomepage : ''}`}>
          {link.title}
        </span>
        <span className={styles.domain}>{link.domain}</span>
        {link.description && variant === 'page' && (
          <span className={styles.desc}>{link.description}</span>
        )}
      </span>
      <span className={styles.arrow} aria-hidden="true">↗</span>
    </a>
  );
}

export default function ExternalLinksBlock({ links, variant }: Props) {
  if (variant === 'homepage') {
    return (
      <div className={styles.block}>
        <div className={styles.list}>
          {links.map((link) => (
            <LinkRow key={link.id} link={link} variant="homepage" />
          ))}
        </div>
      </div>
    );
  }

  const grouped = CATEGORY_ORDER.map((cat) => ({
    cat,
    items: links.filter((l) => l.category === cat),
  })).filter((g) => g.items.length > 0);

  return (
    <div className={styles.block}>
      {grouped.map(({ cat, items }) => (
        <div key={cat} className={styles.group}>
          <div className={styles.categoryLabel}>
            <span className={styles.labelText}>{CATEGORY_LABEL[cat]}</span>
            <span className={styles.hairline} aria-hidden="true" />
          </div>
          <div className={styles.list}>
            {items.map((link) => (
              <LinkRow key={link.id} link={link} variant="page" />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
