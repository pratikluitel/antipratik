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

export default function ContactEmail() {
  return (
    <Html lang="en">
      <Head />
      <Preview>New contact message from __NAME__</Preview>
      <Body style={body}>
        <Container style={container}>
          <Section style={header}>
            <Text style={label}>Contact Form</Text>
            <Text style={heading}>New message</Text>
          </Section>

          <Section style={senderSection}>
            <Text style={fieldLabel}>From</Text>
            <Text style={fieldValue}>__NAME__</Text>
            <Text style={fieldLabel}>Email</Text>
            <Link href="mailto:__EMAIL__" style={emailLink}>__EMAIL__</Link>
          </Section>

          <Section style={messageBox}>
            <Text style={messageText}>__MESSAGE__</Text>
          </Section>

          <Section style={footer}>
            <Text style={footerText}>
              This message was sent via the contact form on your site.
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
  paddingBottom: '24px',
  borderBottom: '1px solid #1E2535',
};

const label: React.CSSProperties = {
  margin: '0 0 8px',
  color: '#7A9AB4',
  fontSize: '12px',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  fontWeight: 600,
};

const heading: React.CSSProperties = {
  margin: 0,
  fontFamily: "'DM Serif Display', Georgia, 'Times New Roman', serif",
  fontSize: '22px',
  fontWeight: 400,
  color: '#E8E4DC',
  lineHeight: '1.3',
};

const senderSection: React.CSSProperties = {
  padding: '24px 0 16px',
};

const fieldLabel: React.CSSProperties = {
  margin: '0 0 4px',
  color: '#7A9AB4',
  fontSize: '12px',
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
  display: 'block',
};

const fieldValue: React.CSSProperties = {
  margin: '0 0 14px',
  color: '#E8E4DC',
  fontSize: '15px',
};

const emailLink: React.CSSProperties = {
  color: '#4A7FBB',
  fontSize: '15px',
  textDecoration: 'none',
};

const messageBox: React.CSSProperties = {
  padding: '20px',
  backgroundColor: '#181D28',
  borderRadius: '4px',
  border: '1px solid #1E2535',
};

const messageText: React.CSSProperties = {
  margin: 0,
  color: '#B8B4AC',
  fontSize: '15px',
  lineHeight: '1.75',
  whiteSpace: 'pre-wrap',
};

const footer: React.CSSProperties = {
  paddingTop: '24px',
};

const footerText: React.CSSProperties = {
  margin: 0,
  color: '#4A6A84',
  fontSize: '12px',
};
