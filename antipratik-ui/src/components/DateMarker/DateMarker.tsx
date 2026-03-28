import styles from './DateMarker.module.css';

interface Props {
  date: string; // ISO 8601 string
}

export default function DateMarker({ date }: Props) {
  const formatted = new Intl.DateTimeFormat('en', {
    month: 'short',
    year: 'numeric',
  }).format(new Date(date));

  return (
    <div className={styles.marker}>
      <time dateTime={date}>{formatted}</time>
    </div>
  );
}
