import {
  Body,
  Container,
  Head,
  Html,
  Link,
  Preview,
  Section,
  Text,
} from '@react-email/components';
import * as React from 'react';

export default function NewsletterEmail() {
  return (
    <Html lang="en">
      <Head />
      <Preview>__SITE_NAME__</Preview>
      <Body style={body}>
        <Container style={container}>

          {/* Header */}
          <Section style={header}>
            <Link href="__SITE_DOMAIN__" style={siteLink}>__SITE_NAME__</Link>
          </Section>

          {/* Caption — Go strips this whole block if empty */}
          <Section style={captionSection}>
            <Text style={captionText}>__CAPTION__</Text>
          </Section>

          {/* Posts — Go substitutes __POSTS_HTML__ with generated post rows */}
          <div dangerouslySetInnerHTML={{ __html: '__POSTS_HTML__' }} />

          {/* Footer */}
          <Section style={footer}>
            <Text style={footerText}>
              You&apos;re receiving this because you subscribed at{' '}
              <Link href="__SITE_DOMAIN__" style={footerLink}>__SITE_DOMAIN__</Link>.
            </Text>
            <Link
              href="__SITE_DOMAIN__/unsubscribe?token=__UNSUBSCRIBE_TOKEN__"
              style={unsubLink}
            >
              Unsubscribe
            </Link>
          </Section>

        </Container>
      </Body>
    </Html>
  );
}

const body: React.CSSProperties = {
  margin: 0,
  padding: 0,
  backgroundColor: '#0F1118',
  fontFamily: "'DM Sans', Arial, Helvetica, sans-serif",
};

const container: React.CSSProperties = {
  maxWidth: '600px',
  margin: '0 auto',
  padding: '40px 20px',
};

const header: React.CSSProperties = {
  paddingBottom: '28px',
  borderBottom: '1px solid #1E2535',
};

const siteLink: React.CSSProperties = {
  textDecoration: 'none',
  color: '#E8E4DC',
  fontFamily: "'DM Serif Display', Georgia, 'Times New Roman', serif",
  fontSize: '22px',
  fontWeight: 400,
  letterSpacing: '0.01em',
};

const captionSection: React.CSSProperties = {
  padding: '24px 0 0',
};

const captionText: React.CSSProperties = {
  margin: 0,
  color: '#7A9AB4',
  fontSize: '14px',
  lineHeight: '1.65',
  fontStyle: 'italic',
};

const footer: React.CSSProperties = {
  paddingTop: '32px',
  paddingBottom: '8px',
  borderTop: '1px solid #1E2535',
};

const footerText: React.CSSProperties = {
  margin: '0 0 8px',
  color: '#4A6A84',
  fontSize: '12px',
  lineHeight: '1.5',
};

const footerLink: React.CSSProperties = {
  color: '#4A6A84',
};

const unsubLink: React.CSSProperties = {
  color: '#4A6A84',
  fontSize: '12px',
  textDecoration: 'underline',
};
