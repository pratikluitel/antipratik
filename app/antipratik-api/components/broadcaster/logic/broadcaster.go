package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/components/broadcaster"
	"github.com/pratikluitel/antipratik/components/broadcaster/lib/resend"
	"github.com/pratikluitel/antipratik/components/broadcaster/store"
	"github.com/pratikluitel/antipratik/components/posts"
)

// emailRegex validates basic RFC 5322 email format.
var emailRegex = regexp.MustCompile(`(?i)^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// batchSize is the maximum number of emails sent per 24-hour window.
const batchSize = 100

// EmailSender abstracts the Resend client for injection and testing.
type EmailSender interface {
	Send(ctx context.Context, req resend.SendRequest) error
}

// broadcasterConfig mirrors the fields from config that the logic layer needs.
type broadcasterConfig struct {
	AdminEmail string
	SiteDomain string
	SiteName   string
	FromName   string
}

// broadcasterLogic implements broadcaster.BroadcasterLogic.
type broadcasterLogic struct {
	store            broadcaster.BroadcasterStore
	sender           EmailSender
	posts            PostService
	log              logging.Logger
	cfg              broadcasterConfig
	newsletterTmpl   string
	confirmationTmpl string
	contactTmpl      string
}

// NewBroadcasterLogic creates a new broadcasterLogic.
// Email templates are embedded at compile time from the app/emails React Email build output.
func NewBroadcasterLogic(
	s broadcaster.BroadcasterStore,
	sender EmailSender,
	postsSvc posts.PostsService,
	adminEmail, siteDomain, siteName, fromName string,
	log logging.Logger,
) (broadcaster.BroadcasterLogic, error) {
	svc := &broadcasterLogic{
		store:  s,
		sender: sender,
		posts:  newPostAdapter(postsSvc),
		cfg: broadcasterConfig{
			AdminEmail: adminEmail,
			SiteDomain: siteDomain,
			SiteName:   siteName,
			FromName:   fromName,
		},
		log:              log,
		newsletterTmpl:   defaultNewsletterTmpl,
		confirmationTmpl: defaultConfirmationTmpl,
		contactTmpl:      defaultContactTmpl,
	}

	// Rescue any buffered sends from a previous run/crash
	svc.ResumePendingDispatches(context.Background())

	return svc, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func validateEmail(address string) error {
	address = strings.TrimSpace(address)
	if address == "" {
		return commonerrors.New("email address cannot be empty")
	}
	if len(address) > 254 {
		return commonerrors.New("email address is too long")
	}
	if !emailRegex.MatchString(address) {
		return commonerrors.New("please enter a valid email address")
	}
	return nil
}

func postAccentColor(postType string) string {
	switch postType {
	case "essay":
		return "#4A7FBB"
	case "music":
		return "#E03E35"
	case "photo":
		return "#5E9E6A"
	case "short":
		return "#D4A832"
	case "video":
		return "#4A7C6F"
	case "link":
		return "#7A8890"
	default:
		return "#7A9AB4"
	}
}

// broadcastData is the JSON payload stored in broadcasts.data.
type broadcastData struct {
	Caption string   `json:"caption"`
	PostIDs []string `json:"postIDs"`
}

func marshalData(caption string, postIDs []string) (string, error) {
	b, err := json.Marshal(broadcastData{Caption: caption, PostIDs: postIDs})
	if err != nil {
		return "", fmt.Errorf("marshal broadcast data: %w", err)
	}
	return string(b), nil
}

func unmarshalData(raw string) (broadcastData, error) {
	var d broadcastData
	if raw == "" || raw == "{}" {
		return d, nil
	}
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		return d, fmt.Errorf("unmarshal broadcast data: %w", err)
	}
	return d, nil
}

// postLinkURL returns the site URL used for email click-throughs per post type.
func (svc *broadcasterLogic) postLinkURL(p PostSummary) string {
	switch p.Type {
	case "essay":
		if p.Slug != "" {
			return svc.cfg.SiteDomain + "/" + p.Slug
		}
	case "photo":
		return svc.cfg.SiteDomain + "/feed?photo=" + p.ID
	case "music":
		return svc.cfg.SiteDomain + "/feed?track=" + p.ID
	case "video":
		return p.VideoURL
	case "link":
		return p.LinkURL
	}
	return svc.cfg.SiteDomain
}

// absoluteURL returns an absolute URL by prepending the site domain to relative paths.
func (svc *broadcasterLogic) absoluteURL(u string) string {
	if u == "" || strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return svc.cfg.SiteDomain + u
}

// renderThumbnail emits an <a>-wrapped <img> for a post card, linked to linkURL.
// imgURL is made absolute before use so relative /thumbnails/... paths work in emails.
func (svc *broadcasterLogic) renderThumbnail(imgURL, alt, linkURL string) string {
	imgURL = svc.absoluteURL(imgURL)
	if imgURL == "" {
		return ""
	}
	return `<a href="` + linkURL + `" target="_blank" style="display:block;margin:0 0 16px;">` +
		`<img src="` + imgURL + `" alt="` + htmlEscape(alt) + `" width="560" ` +
		`style="width:100%;max-width:560px;height:auto;border-radius:4px;display:block;"></a>`
}

// renderPostsSection generates the posts HTML inserted into __POSTS_HTML__.
// For a single post the full detail is shown; for multiple posts a compact list.
func (svc *broadcasterLogic) renderPostsSection(posts []PostSummary) string {
	if len(posts) == 0 {
		return ""
	}
	var b strings.Builder
	isSingle := len(posts) == 1

	if isSingle {
		p := posts[0]
		color := postAccentColor(p.Type)
		link := svc.postLinkURL(p)
		b.WriteString(`<table align="center" width="100%" border="0" cellpadding="0" cellspacing="0" style="padding:28px 0">`)
		b.WriteString(`<tbody><tr><td>`)
		b.WriteString(`<div style="margin-bottom:14px;"><span style="font-size:10px;font-weight:600;letter-spacing:0.10em;text-transform:uppercase;color:` + color + `;padding:3px 8px;border:1px solid ` + color + `;border-radius:3px;display:inline-block;">` + htmlEscape(p.Type) + `</span></div>`)
		switch p.Type {
		case "photo":
			imgURL := p.ThumbnailMediumURL
			if imgURL == "" {
				imgURL = p.ImageURL
			}
			b.WriteString(svc.renderThumbnail(imgURL, "Photo", link))
			b.WriteString(`<div style="margin-top:8px;"><a href="` + link + `" target="_blank" style="display:inline-block;padding:10px 22px;background-color:#1E2535;color:#E8E4DC;text-decoration:none;border-radius:4px;font-size:13px;font-weight:500;">View photos &#x2192;</a></div>`)
		case "music":
			b.WriteString(`<h1 style="margin:0 0 16px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:26px;font-weight:400;color:#E8E4DC;line-height:1.3;">` + htmlEscape(p.Title) + `</h1>`)
			imgURL := p.AlbumArtMediumURL
			b.WriteString(svc.renderThumbnail(imgURL, htmlEscape(p.Title), link))
			b.WriteString(`<div style="margin-top:8px;"><a href="` + link + `" target="_blank" style="display:inline-block;padding:10px 22px;background-color:#E03E35;color:#E8E4DC;text-decoration:none;border-radius:4px;font-size:13px;font-weight:500;">&#9654; Listen on the site</a></div>`)
		case "video":
			b.WriteString(`<h1 style="margin:0 0 16px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:26px;font-weight:400;color:#E8E4DC;line-height:1.3;">` + htmlEscape(p.Title) + `</h1>`)
			b.WriteString(svc.renderThumbnail(p.ThumbnailMediumURL, htmlEscape(p.Title), link))
			b.WriteString(`<div style="margin-top:8px;"><a href="` + link + `" target="_blank" style="display:inline-block;padding:10px 22px;background-color:#1E2535;color:#E8E4DC;text-decoration:none;border-radius:4px;font-size:13px;font-weight:500;">Watch &#x2192;</a></div>`)
		case "link":
			b.WriteString(`<h1 style="margin:0 0 16px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:26px;font-weight:400;color:#E8E4DC;line-height:1.3;">` + htmlEscape(p.Title) + `</h1>`)
			b.WriteString(svc.renderThumbnail(p.ThumbnailMediumURL, htmlEscape(p.Title), link))
			if p.Domain != "" {
				b.WriteString(`<p style="margin:0 0 10px;font-size:11px;color:#7A9AB4;letter-spacing:0.05em;text-transform:uppercase;">` + htmlEscape(p.Domain) + `</p>`)
			}
			if p.Excerpt != "" {
				b.WriteString(`<p style="margin:0 0 20px;color:#B8B4AC;font-size:15px;line-height:1.7;">` + htmlEscape(p.Excerpt) + `</p>`)
			}
			b.WriteString(`<div><a href="` + link + `" target="_blank" style="display:inline-block;padding:10px 22px;background-color:#1E2535;color:#E8E4DC;text-decoration:none;border-radius:4px;font-size:13px;font-weight:500;">Visit &#x2192;</a></div>`)
		case "essay":
			b.WriteString(`<h1 style="margin:0 0 20px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:26px;font-weight:400;color:#E8E4DC;line-height:1.3;">` + htmlEscape(p.Title) + `</h1>`)
			b.WriteString(`<div style="color:#B8B4AC;font-size:15px;line-height:1.75;">` + p.Body + `</div>`)
			if link != "" {
				b.WriteString(`<div style="margin-top:24px;"><a href="` + link + `" style="display:inline-block;padding:10px 22px;background-color:#1E2535;color:#E8E4DC;text-decoration:none;border-radius:4px;font-size:13px;font-weight:500;">Read on the site &#x2192;</a></div>`)
			}
		default:
			if p.Title != "" {
				b.WriteString(`<h1 style="margin:0 0 20px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:26px;font-weight:400;color:#E8E4DC;line-height:1.3;">` + htmlEscape(p.Title) + `</h1>`)
			}
			if p.Excerpt != "" {
				b.WriteString(`<p style="margin:0 0 20px;color:#B8B4AC;font-size:15px;line-height:1.7;">` + htmlEscape(p.Excerpt) + `</p>`)
			}
		}
		b.WriteString(`</td></tr></tbody></table>`)
	} else {
		for _, p := range posts {
			color := postAccentColor(p.Type)
			link := svc.postLinkURL(p)
			b.WriteString(`<table align="center" width="100%" border="0" cellpadding="0" cellspacing="0" style="padding:24px 0;border-bottom:1px solid #1E2535">`)
			b.WriteString(`<tbody><tr><td>`)
			b.WriteString(`<div style="margin-bottom:10px;"><span style="font-size:10px;font-weight:600;letter-spacing:0.10em;text-transform:uppercase;color:` + color + `;padding:3px 8px;border:1px solid ` + color + `;border-radius:3px;display:inline-block;">` + htmlEscape(p.Type) + `</span></div>`)
			switch p.Type {
			case "photo":
				imgURL := p.ThumbnailMediumURL
				if imgURL == "" {
					imgURL = p.ImageURL
				}
				b.WriteString(svc.renderThumbnail(imgURL, "Photo", link))
				b.WriteString(`<a href="` + link + `" target="_blank" style="color:#7A9AB4;text-decoration:none;font-size:13px;font-weight:500;">View photos &#x2192;</a>`)
			case "music":
				b.WriteString(`<h2 style="margin:0 0 12px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:19px;font-weight:400;color:#E8E4DC;line-height:1.35;">` + htmlEscape(p.Title) + `</h2>`)
				b.WriteString(svc.renderThumbnail(p.AlbumArtMediumURL, htmlEscape(p.Title), link))
				b.WriteString(`<a href="` + link + `" target="_blank" style="color:#E03E35;text-decoration:none;font-size:13px;font-weight:500;">&#9654; Listen &#x2192;</a>`)
			case "video":
				b.WriteString(`<h2 style="margin:0 0 12px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:19px;font-weight:400;color:#E8E4DC;line-height:1.35;">` + htmlEscape(p.Title) + `</h2>`)
				b.WriteString(svc.renderThumbnail(p.ThumbnailMediumURL, htmlEscape(p.Title), link))
				b.WriteString(`<a href="` + link + `" target="_blank" style="color:#7A9AB4;text-decoration:none;font-size:13px;font-weight:500;">Watch &#x2192;</a>`)
			case "link":
				b.WriteString(`<h2 style="margin:0 0 8px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:19px;font-weight:400;color:#E8E4DC;line-height:1.35;">` + htmlEscape(p.Title) + `</h2>`)
				if p.Domain != "" {
					b.WriteString(`<p style="margin:0 0 10px;font-size:11px;color:#7A9AB4;letter-spacing:0.05em;text-transform:uppercase;">` + htmlEscape(p.Domain) + `</p>`)
				}
				b.WriteString(svc.renderThumbnail(p.ThumbnailMediumURL, htmlEscape(p.Title), link))
				if p.Excerpt != "" {
					b.WriteString(`<p style="margin:0 0 12px;color:#B8B4AC;font-size:14px;line-height:1.65;">` + htmlEscape(p.Excerpt) + `</p>`)
				}
				b.WriteString(`<a href="` + link + `" target="_blank" style="color:#7A9AB4;text-decoration:none;font-size:13px;font-weight:500;">Visit &#x2192;</a>`)
			default:
				if p.Title != "" {
					b.WriteString(`<h2 style="margin:0 0 10px;font-family:'DM Serif Display',Georgia,'Times New Roman',serif;font-size:19px;font-weight:400;color:#E8E4DC;line-height:1.35;">` + htmlEscape(p.Title) + `</h2>`)
				}
				if p.Excerpt != "" {
					b.WriteString(`<p style="margin:0 0 14px;color:#B8B4AC;font-size:14px;line-height:1.65;">` + htmlEscape(p.Excerpt) + `</p>`)
				}
				if link != "" && link != svc.cfg.SiteDomain {
					b.WriteString(`<a href="` + link + `" style="color:#7A9AB4;text-decoration:none;font-size:13px;font-weight:500;">Read more &#x2192;</a>`)
				}
			}
			b.WriteString(`</td></tr></tbody></table>`)
		}
	}
	return b.String()
}

// htmlEscape escapes s for safe insertion into HTML text/attribute contexts.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	return s
}

// renderConfirmationEmail substitutes tokens into the confirmation template.
func (svc *broadcasterLogic) renderConfirmationEmail(confirmURL string) string {
	return strings.NewReplacer(
		"__CONFIRM_URL__", confirmURL,
		"__SITE_NAME__", htmlEscape(svc.cfg.SiteName),
		"__SITE_DOMAIN__", svc.cfg.SiteDomain,
	).Replace(svc.confirmationTmpl)
}

// renderNewsletter renders the newsletter template for a broadcast.
// The unsubscribe link contains the literal placeholder __UNSUBSCRIBE_TOKEN__
// which is substituted at send time per-subscriber.
func (svc *broadcasterLogic) renderNewsletter(ctx context.Context, _ string, caption string, postIDs []string) (string, error) {
	posts, err := svc.posts.GetPostsByIDs(ctx, postIDs)
	if err != nil {
		return "", fmt.Errorf("fetch posts: %w", err)
	}

	postsHTML := svc.renderPostsSection(posts)

	r := strings.NewReplacer(
		"__SITE_NAME__", htmlEscape(svc.cfg.SiteName),
		"__SITE_DOMAIN__", svc.cfg.SiteDomain,
		"__CAPTION__", htmlEscape(caption),
		"__POSTS_HTML__", postsHTML,
	)
	return r.Replace(svc.newsletterTmpl), nil
}

// ── Subscriber lifecycle ──────────────────────────────────────────────────────

// Subscribe validates the address, registers the subscriber, and sends a confirmation email.
func validateSubscriberType(subType string) error {
	if subType != "email" {
		return commonerrors.New("unsupported subscriber type; only \"email\" is supported")
	}
	return nil
}

func (svc *broadcasterLogic) Subscribe(ctx context.Context, subType, address string) error {
	if err := validateSubscriberType(subType); err != nil {
		return err
	}
	address = strings.ToLower(strings.TrimSpace(address))
	if err := validateEmail(address); err != nil {
		return err
	}

	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("broadcasterLogic.Subscribe: %w", err)
	}

	if err := svc.store.RegisterSubscriber(ctx, subType, address, token); err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			return commonerrors.New("address already subscribed")
		}
		return fmt.Errorf("broadcasterLogic.Subscribe: %w", err)
	}

	confirmURL := svc.cfg.SiteDomain + "/confirm?token=" + token
	html := svc.renderConfirmationEmail(confirmURL)

	if err := svc.sender.Send(ctx, resend.SendRequest{
		To:      []string{address},
		Subject: "Confirm your subscription",
		HTML:    html,
	}); err != nil {
		svc.log.Error("broadcasterLogic.Subscribe send confirmation", "err", err)
	}

	return nil
}

// SendConfirmationEmails generates fresh tokens for unconfirmed subscribers and
// sends them confirmation emails. Returns the count sent.
// This is a temporary admin-only operation for legacy subscribers.
func (svc *broadcasterLogic) SendConfirmationEmails(ctx context.Context, subType string) (int, error) {
	subs, err := svc.store.GetUnconfirmedSubscribers(ctx, subType)
	if err != nil {
		return 0, fmt.Errorf("broadcasterLogic.SendConfirmationEmails: %w", err)
	}

	sent := 0
	for _, sub := range subs {
		token, err := generateToken()
		if err != nil {
			return sent, fmt.Errorf("broadcasterLogic.SendConfirmationEmails: %w", err)
		}
		if err := svc.store.SetSubscriberToken(ctx, sub.Address, token); err != nil {
			svc.log.Error("SendConfirmationEmails set token", "address", sub.Address, "err", err)
			continue
		}

		confirmURL := svc.cfg.SiteDomain + "/confirm?token=" + token
		if err := svc.sender.Send(ctx, resend.SendRequest{
			To:      []string{sub.Address},
			Subject: "Confirm your subscription",
			HTML:    svc.renderConfirmationEmail(confirmURL),
		}); err != nil {
			svc.log.Error("SendConfirmationEmails send", "address", sub.Address, "err", err)
			continue
		}
		sent++
	}
	return sent, nil
}

// ConfirmSubscription marks the subscriber with the given token as confirmed.
func (svc *broadcasterLogic) ConfirmSubscription(ctx context.Context, token string) error {
	if err := commonerrors.RequireNonEmpty("token", token); err != nil {
		return err
	}
	if err := svc.store.ConfirmSubscriber(ctx, token); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return store.ErrNotFound
		}
		return fmt.Errorf("broadcasterLogic.ConfirmSubscription: %w", err)
	}
	return nil
}

// Unsubscribe marks the subscriber with the given token as unsubscribed.
func (svc *broadcasterLogic) Unsubscribe(ctx context.Context, token string) error {
	if err := commonerrors.RequireNonEmpty("token", token); err != nil {
		return err
	}
	if err := svc.store.UnsubscribeByToken(ctx, token); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return store.ErrNotFound
		}
		return fmt.Errorf("broadcasterLogic.Unsubscribe: %w", err)
	}
	return nil
}

// GetSubscribers returns all subscribers as safe summary objects.
func (svc *broadcasterLogic) GetSubscribers(ctx context.Context, subType string) ([]broadcaster.SubscriberSummary, error) {
	subs, err := svc.store.GetAllSubscribers(ctx, subType)
	if err != nil {
		return nil, fmt.Errorf("GetAllSubscribers: %w", err)
	}
	var out []broadcaster.SubscriberSummary
	for _, sub := range subs {
		summary := broadcaster.SubscriberSummary{
			Type:      sub.Type,
			Address:   sub.Address,
			Confirmed: sub.Confirmed,
			CreatedAt: sub.CreatedAt.Format(time.RFC3339),
		}
		if sub.ConfirmedAt != nil {
			ts := sub.ConfirmedAt.Format(time.RFC3339)
			summary.ConfirmedAt = &ts
		}
		if sub.UnsubscribedAt != nil {
			ts := sub.UnsubscribedAt.Format(time.RFC3339)
			summary.UnsubscribedAt = &ts
		}
		out = append(out, summary)
	}
	if out == nil {
		out = []broadcaster.SubscriberSummary{} // prevent null json
	}
	return out, nil
}

// DeleteSubscriber hard-deletes the subscriber with the given address.
func (svc *broadcasterLogic) DeleteSubscriber(ctx context.Context, address string) error {
	if err := commonerrors.RequireNonEmpty(address, "address"); err != nil {
		return err
	}
	if err := svc.store.DeleteSubscriber(ctx, address); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return store.ErrNotFound
		}
		return fmt.Errorf("broadcasterLogic.DeleteSubscriber: %w", err)
	}
	return nil
}

// ── Broadcast management ──────────────────────────────────────────────────────

func (svc *broadcasterLogic) validateBroadcastInput(bType, title string, postIDs []string) error {
	if bType != "email" {
		return commonerrors.New("unsupported broadcast type; only \"email\" is supported")
	}
	if err := commonerrors.RequireNonEmpty("title", title); err != nil {
		return err
	}
	if len(title) > 200 {
		return commonerrors.New("title cannot be longer than 200 characters")
	}
	if len(postIDs) == 0 {
		return commonerrors.New("postIDs must contain at least one post ID")
	}
	for _, id := range postIDs {
		if strings.TrimSpace(id) == "" {
			return commonerrors.New("postIDs must not contain empty values")
		}
	}
	return nil
}

// CreateBroadcast validates input, renders the HTML, and saves the broadcast.
func (svc *broadcasterLogic) CreateBroadcast(ctx context.Context, input broadcaster.BroadcastInput) (broadcaster.BroadcastPreview, error) {
	if err := svc.validateBroadcastInput(input.Type, input.Title, input.PostIDs); err != nil {
		return broadcaster.BroadcastPreview{}, err
	}
	if len(input.Caption) > 500 {
		return broadcaster.BroadcastPreview{}, commonerrors.New("caption cannot be longer than 500 characters")
	}

	html, err := svc.renderNewsletter(ctx, input.Title, input.Caption, input.PostIDs)
	if err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.CreateBroadcast: %w", err)
	}

	data, err := marshalData(input.Caption, input.PostIDs)
	if err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.CreateBroadcast: %w", err)
	}

	id, err := svc.store.CreateBroadcast(ctx, broadcaster.StoreBroadcastInput{
		Type:  input.Type,
		Title: input.Title,
		Data:  data,
	}, html)
	if err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.CreateBroadcast: %w", err)
	}
	return broadcaster.BroadcastPreview{ID: id, HTML: html}, nil
}

// UpdateBroadcast re-validates, re-renders, and updates the stored broadcast.
func (svc *broadcasterLogic) UpdateBroadcast(ctx context.Context, id int64, input broadcaster.BroadcastUpdateInput) (broadcaster.BroadcastPreview, error) {
	if err := svc.validateBroadcastInput("email", input.Title, input.PostIDs); err != nil {
		return broadcaster.BroadcastPreview{}, err
	}
	if len(input.Caption) > 500 {
		return broadcaster.BroadcastPreview{}, commonerrors.New("caption cannot be longer than 500 characters")
	}

	// Verify the broadcast exists before re-rendering.
	existing, err := svc.store.GetBroadcast(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return broadcaster.BroadcastPreview{}, store.ErrNotFound
		}
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.UpdateBroadcast: %w", err)
	}
	_ = existing // type and id are immutable; nothing to copy

	html, err := svc.renderNewsletter(ctx, input.Title, input.Caption, input.PostIDs)
	if err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.UpdateBroadcast: %w", err)
	}

	data, err := marshalData(input.Caption, input.PostIDs)
	if err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.UpdateBroadcast: %w", err)
	}

	if err := svc.store.UpdateBroadcast(ctx, id, broadcaster.StoreBroadcastUpdateInput{
		Title: input.Title,
		Data:  data,
	}, html); err != nil {
		return broadcaster.BroadcastPreview{}, fmt.Errorf("broadcasterLogic.UpdateBroadcast: %w", err)
	}
	return broadcaster.BroadcastPreview{ID: id, HTML: html}, nil
}

// GetBroadcasts returns all broadcasts of the given type with their send summaries.
func (svc *broadcasterLogic) GetBroadcasts(ctx context.Context, bType string) ([]broadcaster.BroadcastSummary, error) {
	rows, err := svc.store.GetBroadcasts(ctx, bType)
	if err != nil {
		return nil, fmt.Errorf("broadcasterLogic.GetBroadcasts: %w", err)
	}

	out := make([]broadcaster.BroadcastSummary, 0, len(rows))
	for _, r := range rows {
		summary, err := svc.store.GetBroadcastSendSummary(ctx, r.ID)
		if err != nil {
			return nil, fmt.Errorf("broadcasterLogic.GetBroadcasts summary: %w", err)
		}
		d, _ := unmarshalData(r.Data)
		out = append(out, broadcaster.BroadcastSummary{
			ID:        r.ID,
			Type:      r.Type,
			Title:     r.Title,
			Caption:   d.Caption,
			PostIDs:   d.PostIDs,
			EmailBody: r.EmailBody,
			Buffered:  summary.Buffered,
			Success:   summary.Success,
			Failed:    summary.Failed,
		})
	}
	return out, nil
}

// ── Dispatch ──────────────────────────────────────────────────────────────────

// DispatchBroadcast creates BUFFERED send records for all confirmed subscribers
// and spawns a background goroutine to process due batches asynchronously.
// Returns the number of BUFFERED records created.
func (svc *broadcasterLogic) DispatchBroadcast(ctx context.Context, id int64) (int, error) {
	broadcast, err := svc.store.GetBroadcast(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, store.ErrNotFound
		}
		return 0, fmt.Errorf("broadcasterLogic.DispatchBroadcast: %w", err)
	}
	if broadcast.EmailBody == "" {
		return 0, commonerrors.New("broadcast has no rendered email body")
	}

	subs, err := svc.store.GetConfirmedSubscribers(ctx, "email")
	if err != nil {
		return 0, fmt.Errorf("broadcasterLogic.DispatchBroadcast get subscribers: %w", err)
	}
	if len(subs) == 0 {
		return 0, nil
	}

	// Build per-subscriber send inputs with batch-staggered schedule.
	sends := make([]broadcaster.BroadcastSendInput, len(subs))
	now := time.Now().UTC()
	for i, sub := range subs {
		batchIndex := i / batchSize
		sends[i] = broadcaster.BroadcastSendInput{
			SubscriberID: sub.ID,
			ScheduledAt:  now.Add(time.Duration(batchIndex) * 24 * time.Hour),
		}
	}

	if err := svc.store.CreateBroadcastSends(ctx, id, sends); err != nil {
		return 0, fmt.Errorf("broadcasterLogic.DispatchBroadcast create sends: %w", err)
	}

	// Spawn a background goroutine that processes due batches until all are sent.
	go svc.runDispatch(id, broadcast.EmailBody, broadcast.Title)

	return len(sends), nil
}

// runDispatch is the background worker that processes due broadcast_sends rows.
func (svc *broadcasterLogic) runDispatch(broadcastID int64, emailBody, subject string) {
	ctx := context.Background()
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// Process immediately for batch 0.
	svc.processDueSends(ctx, broadcastID, emailBody, subject)

	for range ticker.C {
		remaining, err := svc.store.GetRemainingBuffered(ctx, broadcastID)
		if err != nil {
			svc.log.Error("runDispatch check remaining", "broadcast_id", broadcastID, "err", err)
		}
		if remaining == 0 {
			return
		}
		svc.processDueSends(ctx, broadcastID, emailBody, subject)
	}
}

// processDueSends fetches and sends all due BUFFERED rows for the broadcast.
func (svc *broadcasterLogic) processDueSends(ctx context.Context, broadcastID int64, emailBody, subject string) {
	sends, err := svc.store.GetDueSends(ctx, broadcastID)
	if err != nil {
		svc.log.Error("processDueSends get", "broadcast_id", broadcastID, "err", err)
		return
	}

	for _, send := range sends {
		html := strings.ReplaceAll(emailBody, "__UNSUBSCRIBE_TOKEN__", send.SubscriberToken)
		if err := svc.sendWithRetry(ctx, send.SubscriberAddress, subject, html); err != nil {
			msg := err.Error()
			_ = svc.store.UpdateSendStatus(ctx, send.ID, "FAILED", msg, nil)
			svc.log.Error("processDueSends send failed", "broadcast_id", broadcastID, "subscriber", send.SubscriberAddress, "err", err)
			continue
		}
		now := time.Now().UTC()
		_ = svc.store.UpdateSendStatus(ctx, send.ID, "SUCCESS", "", &now)
	}
}

// sendWithRetry attempts to send an email up to 3 times with exponential backoff.
func (svc *broadcasterLogic) sendWithRetry(ctx context.Context, to, subject, html string) error {
	req := resend.SendRequest{To: []string{to}, Subject: subject, HTML: html}
	delays := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for i, delay := range delays {
		if err := svc.sender.Send(ctx, req); err == nil {
			return nil
		} else {
			lastErr = err
			if !resend.IsTransient(err) {
				return err // permanent failure; no retry
			}
		}
		if i < len(delays)-1 {
			time.Sleep(delay)
		}
	}
	return lastErr
}

// GetBroadcastSends returns per-subscriber send details for a broadcast.
func (svc *broadcasterLogic) GetBroadcastSends(ctx context.Context, id int64) ([]broadcaster.BroadcastSendDetail, error) {
	_, err := svc.store.GetBroadcast(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("broadcasterLogic.GetBroadcastSends: %w", err)
	}

	sends, err := svc.store.GetAllBroadcastSends(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("broadcasterLogic.GetBroadcastSends: %w", err)
	}

	out := make([]broadcaster.BroadcastSendDetail, 0, len(sends))
	for _, s := range sends {
		detail := broadcaster.BroadcastSendDetail{
			Address:     s.SubscriberAddress,
			Status:      s.Status,
			Message:     s.Message,
			ScheduledAt: s.ScheduledAt.Format(time.RFC3339),
		}
		if s.SentAt != nil {
			ts := s.SentAt.Format(time.RFC3339)
			detail.SentAt = &ts
		}
		out = append(out, detail)
	}
	return out, nil
}

// ResumePendingDispatches queries all email broadcasts and restarts background loops for any with BUFFERED sends.
func (svc *broadcasterLogic) ResumePendingDispatches(ctx context.Context) {
	broadcasts, err := svc.GetBroadcasts(ctx, "email")
	if err != nil {
		svc.log.Error("ResumePendingDispatches get broadcasts", "err", err)
		return
	}
	resumedCount := 0
	for _, b := range broadcasts {
		if b.Buffered > 0 {
			resumedCount++
			svc.log.Info("ResumePendingDispatches starting background loop", "broadcast_id", b.ID, "buffered", b.Buffered)
			go svc.runDispatch(b.ID, b.EmailBody, b.Title)
		}
	}
	if resumedCount > 0 {
		svc.log.Info("ResumePendingDispatches complete", "resumed_broadcasts", resumedCount)
	}
}

// ── Contact form ──────────────────────────────────────────────────────────────

// SendContactMessage validates input, saves it, and emails the admin.
func (svc *broadcasterLogic) SendContactMessage(ctx context.Context, input broadcaster.ContactInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Message = strings.TrimSpace(input.Message)

	if err := commonerrors.RequireNonEmpty("name", input.Name); err != nil {
		return err
	}
	if len(input.Name) > 100 {
		return commonerrors.New("name cannot be longer than 100 characters")
	}
	if err := validateEmail(input.Email); err != nil {
		return err
	}
	if err := commonerrors.RequireNonEmpty("message", input.Message); err != nil {
		return err
	}
	if len(input.Message) > 2000 {
		return commonerrors.New("message cannot be longer than 2000 characters")
	}

	if err := svc.store.SaveContactMessage(ctx, input.Name, input.Email, input.Message); err != nil {
		return fmt.Errorf("broadcasterLogic.SendContactMessage: %w", err)
	}

	contactHTML := strings.NewReplacer(
		"__NAME__", htmlEscape(input.Name),
		"__EMAIL__", htmlEscape(input.Email),
		"__MESSAGE__", htmlEscape(input.Message),
	).Replace(svc.contactTmpl)

	if svc.cfg.AdminEmail != "" {
		if err := svc.sender.Send(ctx, resend.SendRequest{
			To:      []string{svc.cfg.AdminEmail},
			Subject: "New contact message from " + input.Name,
			HTML:    contactHTML,
		}); err != nil {
			svc.log.Error("broadcasterLogic.SendContactMessage send", "err", err)
		}
	}

	return nil
}
