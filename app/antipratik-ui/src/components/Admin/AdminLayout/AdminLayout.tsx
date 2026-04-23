'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import styles from './AdminLayout.module.css';

interface AdminLayoutProps {
  children: React.ReactNode;
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const router = useRouter();
  const pathname = usePathname();

  function handleLogout() {
    localStorage.removeItem('admin_token');
    router.replace('/admin/login');
  }

  return (
    <div className={styles.shell}>
      <header className={styles.topBar}>
        <div className={styles.topBarLeft}>
          <span className={styles.brand}>antipratik / admin</span>
          <nav className={styles.nav}>
            <Link
              href="/admin/dashboard"
              className={`${styles.navLink} ${pathname.startsWith('/admin/dashboard') ? styles.navLinkActive : ''}`}
            >
              Posts
            </Link>
            <Link
              href="/admin/broadcasts"
              className={`${styles.navLink} ${pathname.startsWith('/admin/broadcasts') ? styles.navLinkActive : ''}`}
            >
              Broadcasts
            </Link>
          </nav>
        </div>
        <button className={styles.logoutBtn} onClick={handleLogout} type="button">
          Logout
        </button>
      </header>
      <main className={styles.content}>{children}</main>
    </div>
  );
}
