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
          I&apos;m a developer, musician, and occasional writer based in Kathmandu, Nepal.
          I build things with code by day and record ambient music by night — usually with
          the window open so the city sounds bleed in.
        </p>
        <p className={styles.bioText}>
          This site is where I collect the work that doesn&apos;t fit anywhere else.
          Essays on craft, music made in small rooms, photographs taken on slow mornings,
          and short notes on whatever is keeping me up at night. No schedule, no format.
          Just the things that feel worth making.
        </p>
        <p className={styles.bioText}>
          I grew up in the shadow of the Himalayas, which taught me two things:
          patience and how to find warmth in small spaces. Both show up in the work,
          whether I mean them to or not.
        </p>

        <div className={styles.currently}>
          <h2 className={styles.currentlyHeading}>Currently</h2>
          <ul className={styles.currentlyList}>
            <li>Working on a new EP — field recordings from the Boudhanath area</li>
            <li>Reading <em>The Snow Leopard</em> by Peter Matthiessen, again</li>
            <li>Listening to a lot of Nils Frahm and the sound of rain on tin roofs</li>
            <li>Somewhere between Kathmandu and wherever the next project takes me</li>
          </ul>
        </div>
      </section>

      <div className={styles.newsletterWrapper}>
        <NewsletterBlock variant="page" />
      </div>
    </main>
  );
}
