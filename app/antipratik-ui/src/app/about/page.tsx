import type { Metadata } from 'next';
import NewsletterBlock from '../../components/NewsletterBlock/NewsletterBlock';
import styles from './page.module.css';

export const metadata: Metadata = {
  title: 'About — antipratik',
  description: 'Developer, musician, writer, and photographer based in Kathmandu, Nepal.',
};

export default function AboutPage() {
  return (
    <main>
      {/* Hero — placeholder for a future photograph */}
      <div className={styles.hero} aria-hidden="true">
        <div className={styles.heroFade} />
        <span className={styles.heroName}>Pratik</span>
      </div>

      <section className={styles.bio}>
        <p className={styles.bioText}>
          I&apos;m a developer, music tinkerer, and occasional writer based in Kathmandu, Nepal.
          I build things with code by day and record ambient music by night.
        </p>
        <p className={styles.bioText}>
          This site is where I collect the work that doesn&apos;t fit anywhere else.
          Essays on craft, music made in small rooms, photographs taken on slow mornings,
          and short notes on whatever is keeping me up at night, without any schedule or format.
          Things that I feel are worth making.
        </p>
        <p className={styles.bioText}>
          I grew up in the Kathmandu, under the shadow of the Himalayas, which taught me two things:
          patience and how to find warmth in small spaces. Both show up in the work,
          whether I mean them to or not.
        </p>
      </section>

      <div className={styles.newsletterWrapper}>
        <NewsletterBlock variant="page" />
      </div>
    </main>
  );
}
