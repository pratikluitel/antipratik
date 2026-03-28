import type { Metadata } from 'next';
import { getLinks } from '../../lib/api';
import ExternalLinksBlock from '../../components/ExternalLinksBlock/ExternalLinksBlock';
import styles from './page.module.css';

export const metadata: Metadata = {
  title: 'Links — antipratik',
  description: 'Things worth reading, watching, and listening to.',
};

export default async function LinksPage() {
  const links = await getLinks();

  return (
    <main className={styles.page}>
      <div className={styles.header}>
        <h1 className={styles.title}>Links</h1>
        <p className={styles.subtitle}>Things worth reading, watching, and listening to.</p>
      </div>
      <ExternalLinksBlock links={links} variant="page" />
    </main>
  );
}
