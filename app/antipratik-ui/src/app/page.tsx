import Link from 'next/link';
import { getPosts, getFeaturedLinks } from '../lib/api';
import HomeFeedClient from '../components/HomeFeedClient/HomeFeedClient';
import ExternalLinksBlock from '../components/ExternalLinksBlock/ExternalLinksBlock';
import NewsletterBlock from '../components/NewsletterBlock/NewsletterBlock';
import styles from './page.module.css';

export default async function Home() {
  const [posts, links] = await Promise.all([getPosts(), getFeaturedLinks()]);

  return (
    <main>
      {/* Hero — intentionally always dark, theme-resistant.
          The hero represents the "Himalayan Dusk" mood of the site
          and should not shift to the light parchment theme. */}
      <section className={styles.hero} style={{ background: '#0F1118' }}>
        <div className={styles.heroInner}>
          <p className={styles.eyebrow}>Kathmandu · Nepal</p>
          <h1 className={styles.displayName}>antipratik</h1>
          <p className={styles.description}>
            Writing, music, and photographs from the edge of the Himalayas.
            Occasionally, thoughts on code.
          </p>
          <div className={styles.tags} aria-hidden="true">
            <span className={`${styles.tag} ${styles.tagEssays}`}>Essays</span>
            <span className={`${styles.tag} ${styles.tagShort}`}>Short</span>
            <span className={`${styles.tag} ${styles.tagMusic}`}>Music</span>
            <span className={`${styles.tag} ${styles.tagPhotos}`}>Photos</span>
            <span className={`${styles.tag} ${styles.tagVideos}`}>Videos</span>
            <span className={`${styles.tag} ${styles.tagLinks}`}>Links</span>
          </div>
        </div>
      </section>

      {/* Feed snippet */}
      <section className={styles.section}>
        <div className={styles.sectionInner}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>Recent</h2>
            <Link href="/feed" className={styles.sectionLink}>View feed →</Link>
          </div>
        </div>
        <div className={styles.feedWrapper}>
          <HomeFeedClient posts={posts.slice(0, 6)} />
        </div>
      </section>

      {/* Links snippet */}
      <section className={styles.section}>
        <div className={styles.sectionInner}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>Links</h2>
            <Link href="/links" className={styles.sectionLink}>View all →</Link>
          </div>
          <ExternalLinksBlock links={links} variant="homepage" />
        </div>
      </section>

      {/* Newsletter */}
      <section className={styles.section}>
        <div className={styles.sectionInner}>
          <NewsletterBlock variant="page" />
        </div>
      </section>
    </main>
  );
}
