import { render } from '@react-email/render';
import * as fs from 'fs';
import * as path from 'path';
import * as React from 'react';

import ConfirmationEmail from '../emails/confirmation';
import ContactEmail from '../emails/contact';
import NewsletterEmail from '../emails/newsletter';

const outDir = path.join(__dirname, '..', 'dist');
fs.mkdirSync(outDir, { recursive: true });

async function renderAll() {
  const templates: Array<{ name: string; component: React.ComponentType }> = [
    { name: 'confirmation', component: ConfirmationEmail },
    { name: 'contact', component: ContactEmail },
    { name: 'newsletter', component: NewsletterEmail },
  ];

  for (const { name, component } of templates) {
    const html = await render(React.createElement(component), { pretty: false });
    const outPath = path.join(outDir, `${name}.html`);
    fs.writeFileSync(outPath, html, 'utf8');
    console.log(`✓ ${name}.html`);
  }
}

renderAll().catch((err) => {
  console.error(err);
  process.exit(1);
});
