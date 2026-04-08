/**
 * Admin panel segment layout.
 * Wraps all /admin routes. The public root layout (Navbar, MusicPlayer) still
 * renders above this, but the admin section provides its own full-screen chrome.
 */
import styles from './layout.module.css';

export default function AdminSegmentLayout({ children }: { children: React.ReactNode }) {
  return <div className={styles.adminRoot}>{children}</div>;
}
