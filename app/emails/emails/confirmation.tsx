import {
  Body,
  Container,
  Head,
  Heading,
  Html,
  Link,
  Preview,
  Section,
  Text,
} from '@react-email/components';
import * as React from 'react';

export default function ConfirmationEmail() {
  return (
    <Html lang="en">
      <Head />
      <Preview>Confirm your subscription</Preview>
      <Body style={body}>
        <Container style={container}>
          <Section style={header}>
            <Link href="__SITE_DOMAIN__" style={siteLink}>__SITE_NAME__</Link>
          </Section>

          <Section style={content}>
            <Heading style={heading}>Confirm your subscription</Heading>
            <Text style={paragraph}>
              Welcome to the AntiPratik newsletter! Click below to confirm your subscription to __SITE_NAME__.
            </Text>
            <Link href="__CONFIRM_URL__" style={button}>
              Confirm subscription
            </Link>
            <Text style={muted}>
              If you didn&apos;t subscribe, you can safely ignore this email.
            </Text>
          </Section>

          <Section style={footer}>
            <Text style={footerText}>
              You&apos;re receiving this because someone signed up at{' '}
              <Link href="__SITE_DOMAIN__" style={footerLink}>__SITE_DOMAIN__</Link>.
            </Text>
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

const content: React.CSSProperties = {
  padding: '32px 0',
};

const heading: React.CSSProperties = {
  margin: '0 0 16px',
  fontFamily: "'DM Serif Display', Georgia, 'Times New Roman', serif",
  fontSize: '24px',
  fontWeight: 400,
  color: '#E8E4DC',
  lineHeight: '1.3',
};

const paragraph: React.CSSProperties = {
  margin: '0 0 24px',
  color: '#B8B4AC',
  fontSize: '15px',
  lineHeight: '1.7',
};

const button: React.CSSProperties = {
  display: 'inline-block',
  padding: '10px 22px',
  backgroundColor: '#1E2535',
  color: '#E8E4DC',
  textDecoration: 'none',
  borderRadius: '4px',
  fontSize: '13px',
  fontWeight: 500,
};

const muted: React.CSSProperties = {
  margin: '24px 0 0',
  color: '#4A6A84',
  fontSize: '13px',
  lineHeight: '1.5',
};

const footer: React.CSSProperties = {
  paddingTop: '24px',
  borderTop: '1px solid #1E2535',
};

const footerText: React.CSSProperties = {
  margin: 0,
  color: '#4A6A84',
  fontSize: '12px',
  lineHeight: '1.5',
};

const footerLink: React.CSSProperties = {
  color: '#4A6A84',
  textDecoration: 'none',
};
