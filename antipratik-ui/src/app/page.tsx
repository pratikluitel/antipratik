import { getPosts } from '../lib/api';

export default async function Home() {
  const posts = await getPosts();

  return (
    <div>
      <main style={{ paddingTop: 'var(--space-10)' }}>
        <h1
          style={{
            fontFamily: 'var(--font-serif)',
            fontSize: 'var(--text-h1)',
            lineHeight: 'var(--lh-heading)',
            marginBottom: 'var(--space-4)',
          }}
        >
          antipratik — scaffold check
        </h1>
        <p style={{ fontSize: 'var(--text-body)', fontFamily: 'var(--font-sans)' }}>
          Current posts: {posts.length}
        </p>
      </main>
    </div>
  );
}
