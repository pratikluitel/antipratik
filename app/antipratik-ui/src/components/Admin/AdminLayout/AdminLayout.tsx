'use client';

import { useRouter } from 'next/navigation';
import styles from './AdminLayout.module.css';

interface AdminLayoutProps {
  children: React.ReactNode;
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const router = useRouter();

  function handleLogout() {
    localStorage.removeItem('admin_token');
    router.replace('/admin/login');
  }

  return (
    <div className={styles.shell}>
      <header className={styles.topBar}>
        <span className={styles.brand}>antipratik / admin</span>
        <button className={styles.logoutBtn} onClick={handleLogout} type="button">
          Logout
        </button>
      </header>
      <main className={styles.content}>{children}</main>
    </div>
  );
}
