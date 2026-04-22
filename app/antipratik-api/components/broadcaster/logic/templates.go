package logic

import _ "embed"

// Email templates are pre-built by the app/emails React Email package.
// To regenerate: run `npm run build` in app/emails/ then copy dist/ to
// components/broadcaster/logic/emails/dist/ (done automatically by Dockerfile.api).

//go:embed emails/dist/newsletter.html
var defaultNewsletterTmpl string

//go:embed emails/dist/confirmation.html
var defaultConfirmationTmpl string

//go:embed emails/dist/contact.html
var defaultContactTmpl string
