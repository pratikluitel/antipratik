import type { LayoutMode } from '../../lib/types';
import styles from './ClusterDivider.module.css';

interface Props {
  from: LayoutMode;
  to: LayoutMode;
}

export default function ClusterDivider(_props: Props) {
  return <div className={styles.divider} />;
}
