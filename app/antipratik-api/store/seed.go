package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// SeedIfEmpty inserts all seed data into the database if the posts table is empty.
// Safe to call on every startup — does nothing if data already exists.
func SeedIfEmpty(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count); err != nil {
		return fmt.Errorf("seed check: %w", err)
	}
	if count > 0 {
		log.Printf("database already has %d posts, skipping seed", count)
		return nil
	}

	log.Println("seeding database...")
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("seed begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	if err := seedPosts(tx); err != nil {
		return fmt.Errorf("seed posts: %w", err)
	}
	if err := seedLinks(tx); err != nil {
		return fmt.Errorf("seed links: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed commit: %w", err)
	}
	log.Println("database seeded successfully")
	return nil
}

// insertPost inserts into posts and post_tags tables.
func insertPost(tx *sql.Tx, id, postType, createdAt string, tags []string) error {
	_, err := tx.Exec("INSERT INTO posts (id, type, created_at) VALUES (?, ?, ?)", id, postType, createdAt)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if _, err := tx.Exec("INSERT INTO post_tags (post_id, tag) VALUES (?, ?)", id, tag); err != nil {
			return err
		}
	}
	return nil
}

func seedPosts(tx *sql.Tx) error {
	// ── Essays ────────────────────────────────────────────────────────────────

	if err := insertPost(tx, "essay-001", "essay", "2026-03-15T08:30:00Z",
		[]string{"philosophy", "music", "impermanence", "buddhism"}); err != nil {
		return err
	}
	_, err := tx.Exec(`INSERT INTO essay_posts (post_id, title, slug, excerpt, body, reading_time_minutes) VALUES (?, ?, ?, ?, ?, ?)`,
		"essay-001",
		"On Impermanence and Code",
		"on-impermanence-and-code",
		"A sound exists only in the moment it is heard. A function exists only as long as the runtime allows. These are not different things.",
		`
A sound exists only in the moment it is heard. Unlike a painting, which persists on its canvas across centuries, a piece of music is always already disappearing — each note a small act of vanishing. This is not a flaw in the medium. It is the medium.

I've been thinking about this since spending a week at a gompa above Langtang, watching monks build a sand mandala. Three days of careful work, every grain placed with intention. On the fourth morning, they swept it into a bowl and poured it into the river.

The mandala was not destroyed. It was completed.

## What Altitude Does to Thinking

There's something about being at 4,000 metres that clarifies the mind. Maybe it's the thinner air, the slower pace of everything, the absence of notifications. Maybe it's proximity to something larger than yourself.

I was refactoring a payment service when I left for the trek. When I returned, I saw it differently. The code I'd been precious about — the clever abstractions, the elegant edge-case handling — none of it would survive the next major version of the library it depended on. I'd been building a sand mandala.

This isn't nihilism. The monks weren't nihilists. The point is that impermanence is not a problem to be solved. It's a condition to be worked *with*.

## The Repository Pattern as Buddhist Practice

There's a pattern in software called the repository pattern. You write a thin layer that stands between your business logic and your data source. When the data source changes — when you move from SQLite to Postgres, or from REST to GraphQL — you change the repository, not the logic.

The business logic doesn't grieve the lost SQLite connection. It never knew about SQLite. It only knew about the interface.

I think about this when I write functions that will outlast the infrastructure they currently run on. The function is the mandala. The infrastructure is the sand.

Write good interfaces. Hold the implementation loosely.

## On Music Production in a City That Doesn't Sleep

Kathmandu is loud in a specific way. It's not the uniform drone of a highway. It's layered — a dog starting a chain reaction across three neighbourhoods, a temple bell, a truck navigating a lane too small for it, rain on corrugated iron.

I record this. Not as samples, exactly — more as texture. The city is a collaborator.

When I sit down to mix a track at 2am, the sounds outside become part of the sonic context, even when I'm wearing headphones. The mix that sounds right in that context is the mix I keep.

The present moment of music is all there is. The past note is memory. The future note is anticipation. Only the note you're hearing now exists.

I try to make music that knows this about itself.`,
		7,
	)
	if err != nil {
		return err
	}

	if err := insertPost(tx, "essay-002", "essay", "2026-02-20T10:00:00Z",
		[]string{"code", "craft", "tools", "unix"}); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO essay_posts (post_id, title, slug, excerpt, body, reading_time_minutes) VALUES (?, ?, ?, ?, ?, ?)`,
		"essay-002",
		"The Terminal as a Musical Instrument",
		"the-terminal-as-a-musical-instrument",
		"A musician who truly knows their instrument stops thinking about the instrument. The same is true for a developer who has made the terminal their own.",
		`
There's a guitarist I know who has played the same guitar for twelve years. He doesn't think about the instrument anymore. He thinks about the music, and the instrument disappears.

I've been using the same terminal configuration for four years. I barely notice it. This is the goal.

## What "Making It Your Own" Actually Means

There's a phase in learning an instrument — and a development environment — where you're constantly aware of the tool. You're thinking about key bindings, about where things live, about the friction between your intention and the outcome.

When this phase ends, something important happens. The tool becomes transparent. You see through it to the problem.

This isn't about memorising shortcuts. It's about building muscle memory through genuine use. The shortcuts you use every day become unconscious. The ones you don't need don't matter.

## The Prompt as Interface Design

Your shell prompt is the most frequently rendered piece of UI in your workflow. I spent more time designing mine than I've spent on any UI component in production code.

It shows: current directory (abbreviated), git branch and status, and the exit code of the last command if it was non-zero. That's it. Everything else is noise.

Less information, more quickly understood, is always better than more information.

## On Portability

I keep my dotfiles in a repository. When I set up a new machine, everything I need is one command away. This is not just convenience — it's a discipline.

Every configuration decision I've made is documented in code. I can read the git log and understand why I made each choice. The setup is reproducible and, more importantly, understandable.

Write tools you can explain.`,
		5,
	)
	if err != nil {
		return err
	}

	if err := insertPost(tx, "essay-003", "essay", "2026-01-10T07:00:00Z",
		[]string{"nepal", "monsoon", "writing", "place"}); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO essay_posts (post_id, title, slug, excerpt, body, reading_time_minutes) VALUES (?, ?, ?, ?, ?, ?)`,
		"essay-003",
		"What the Monsoon Teaches About Latency",
		"what-the-monsoon-teaches-about-latency",
		"For four months each year, Kathmandu operates at a different clock speed. There are things you can only learn by waiting.",
		`# What the Monsoon Teaches About Latency

For four months each year — roughly June through September — the valley fills with cloud and rain. The light changes. The pace changes. You stop planning more than a day ahead because the weather will change your plans anyway.

I've lived through five monsoons here. Each one has taught me something about working with conditions you can't control.

## Latency Is Not Lag

In distributed systems, latency is the time between a request and its response. We optimise to reduce it. We build caches. We precompute. We get frustrated when it exceeds expectations.

But there's a kind of latency that can't be optimised away, and shouldn't be. The latency between planting a thought and seeing what grows from it. The latency between writing a piece of music and understanding whether it's good.

The monsoon enforces this latency. You can't rush the river.

## On Working Slowly

I write better code in the monsoon. I'm not sure why. Maybe the grey light is easier on the eyes. Maybe the rain is good ambient sound. Maybe the general slowing of the city's metabolism affects mine.

I've stopped fighting it. I've started scheduling the exploratory work — the architectural thinking, the research — for the monsoon months. The heads-down execution work for the dry season.

There's wisdom in matching your pace to the environment's pace.

## The Specific Smell of Rain on Hot Stone

There's a word for it: petrichor. The smell of rain on dry earth. In Kathmandu it smells different from anywhere else I've been, because the stone and soil are different.

I've tried to put this smell into music twice. Both times I failed, but the failures were interesting. They led me somewhere else.

Sometimes the unreachable thing is the most useful thing.`,
		6,
	)
	if err != nil {
		return err
	}

	// ── Shorts ────────────────────────────────────────────────────────────────

	if err := insertPost(tx, "short-001", "short", "2026-03-20T14:22:00Z",
		[]string{"kathmandu", "fog", "morning"}); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO short_posts (post_id, body) VALUES (?, ?)",
		"short-001",
		"just spent an hour watching fog roll over Shivapuri. some mornings the city disappears completely. those are the good mornings.",
	); err != nil {
		return err
	}

	if err := insertPost(tx, "short-002", "short", "2026-03-12T02:15:00Z",
		[]string{"code", "night", "rain"}); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO short_posts (post_id, body) VALUES (?, ?)",
		"short-002",
		"refactoring at 2am while rain hits the tin roofs. there is no better context for this kind of work. the city is quiet and the problem is the only thing.",
	); err != nil {
		return err
	}

	if err := insertPost(tx, "short-003", "short", "2026-02-28T09:45:00Z",
		[]string{"boudhanath", "observation", "ritual"}); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO short_posts (post_id, body) VALUES (?, ?)",
		"short-003",
		"watched an old man do his morning kora at Boudha. same direction, same pace, same beads. 108 times. wondered how many of the world's problems would dissolve if we each had one ritual we trusted that completely.",
	); err != nil {
		return err
	}

	// ── Music ─────────────────────────────────────────────────────────────────

	if err := insertPost(tx, "music-001", "music", "2026-03-01T12:00:00Z",
		[]string{"ambient", "kathmandu", "field-recording"}); err != nil {
		return err
	}
	album1 := "Altitude Studies"
	if _, err := tx.Exec("INSERT INTO music_posts (post_id, title, album_art, audio_url, duration, album) VALUES (?, ?, ?, ?, ?, ?)",
		"music-001", "Threshold (Nagarkot Dawn)", "/images/music/threshold.jpg", "", 144, &album1,
	); err != nil {
		return err
	}

	if err := insertPost(tx, "music-002", "music", "2026-01-18T16:00:00Z",
		[]string{"electronic", "folk-fusion", "sarangi"}); err != nil {
		return err
	}
	album2 := "River Sounds"
	if _, err := tx.Exec("INSERT INTO music_posts (post_id, title, album_art, audio_url, duration, album) VALUES (?, ?, ?, ?, ?, ?)",
		"music-002", "Trisuli", "/images/music/trisuli.jpg", "/audio/trisuli.mp3", 144, &album2,
	); err != nil {
		return err
	}

	// ── Photos ────────────────────────────────────────────────────────────────

	if err := insertPost(tx, "photo-001", "photo", "2026-02-14T06:30:00Z",
		[]string{"langtang", "mountains", "winter", "trek"}); err != nil {
		return err
	}
	location1 := "Langtang Valley, 3,500m"
	if _, err := tx.Exec("INSERT INTO photo_posts (post_id, location) VALUES (?, ?)", "photo-001", &location1); err != nil {
		return err
	}
	caption1a := "The valley floor was still dark when the peaks caught the first light."
	caption1b := "Langtang village. A week after the snowfall."
	photoImages1 := []struct {
		url, alt string
		caption  *string
		order    int
	}{
		{"/images/photos/langtang-ridge-dawn.jpg", "Langtang ridge at dawn, thin clouds below the peaks", &caption1a, 0},
		{"/images/photos/langtang-village-snow.jpg", "Stone houses in Langtang village with fresh snow", &caption1b, 1},
		{"/images/photos/langtang-prayer-flags.jpg", "Prayer flags strung between two stone pillars against a clear blue sky", nil, 2},
	}
	for _, img := range photoImages1 {
		if _, err := tx.Exec("INSERT INTO photo_images (post_id, url, alt, caption, sort_order) VALUES (?, ?, ?, ?, ?)",
			"photo-001", img.url, img.alt, img.caption, img.order,
		); err != nil {
			return err
		}
	}

	if err := insertPost(tx, "photo-002", "photo", "2026-01-28T17:00:00Z",
		[]string{"kathmandu", "patan", "durbar-square", "evening"}); err != nil {
		return err
	}
	location2 := "Patan, Lalitpur"
	if _, err := tx.Exec("INSERT INTO photo_posts (post_id, location) VALUES (?, ?)", "photo-002", &location2); err != nil {
		return err
	}
	caption2 := "Patan Durbar, January. The square empties after sunset."
	if _, err := tx.Exec("INSERT INTO photo_images (post_id, url, alt, caption, sort_order) VALUES (?, ?, ?, ?, ?)",
		"photo-002", "/images/photos/patan-durbar-evening.jpg",
		"Patan Durbar Square at dusk, warm light on stone temples", &caption2, 0,
	); err != nil {
		return err
	}

	// ── Videos ────────────────────────────────────────────────────────────────

	if err := insertPost(tx, "video-001", "video", "2026-02-05T11:00:00Z",
		[]string{"music", "studio", "process", "ambient"}); err != nil {
		return err
	}
	playlist1 := "Studio Sessions"
	if _, err := tx.Exec("INSERT INTO video_posts (post_id, title, thumbnail_url, video_url, duration, playlist) VALUES (?, ?, ?, ?, ?, ?)",
		"video-001", `Making "Threshold" — a studio session`,
		"/images/videos/threshold-session-thumb.jpg",
		"https://vimeo.com/placeholder/threshold-session",
		1140, &playlist1,
	); err != nil {
		return err
	}

	// ── Link posts ────────────────────────────────────────────────────────────

	if err := insertPost(tx, "link-001", "link", "2026-03-10T09:00:00Z",
		[]string{"music", "technology", "generative"}); err != nil {
		return err
	}
	desc := "A long-form piece on why most generative music sounds like waiting room ambience, and what the exceptions do differently."
	cat := "music"
	if _, err := tx.Exec("INSERT INTO link_posts (post_id, title, url, domain, description, thumbnail_url, category) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"link-001",
		"Why generative music keeps failing to be interesting",
		"https://example.com/generative-music-failure",
		"example.com",
		&desc, nil, &cat,
	); err != nil {
		return err
	}

	return nil
}

func seedLinks(tx *sql.Tx) error {
	links := []struct {
		id, title, url, domain, description string
		featured                             int
		category                             string
	}{
		{"link-ext-001", "antipratik on SoundCloud", "https://soundcloud.com/antipratik", "soundcloud.com",
			"Ambient and electronic tracks. Field recordings from Kathmandu and the Himalayas.", 1, "music"},
		{"link-ext-002", "antipratik on Bandcamp", "https://antipratik.bandcamp.com", "bandcamp.com",
			"Full albums and EPs. Pay what you want or nothing — the music is meant to be heard.", 1, "music"},
		{"link-ext-003", "Essays on Substack", "https://antipratik.substack.com", "substack.com",
			"Longer essays on music, code, and living at altitude. Published when ready, not on a schedule.", 1, "writing"},
		{"link-ext-004", "Writing on Medium", "https://medium.com/@antipratik", "medium.com",
			"Technical writing on distributed systems, developer tooling, and the intersection of craft in music and code.", 0, "writing"},
		{"link-ext-005", "YouTube — Studio Sessions", "https://youtube.com/@antipratik", "youtube.com",
			"Behind-the-scenes studio sessions, gear walkthroughs, and long-form process videos.", 1, "video"},
		{"link-ext-006", "Vimeo — Short Films", "https://vimeo.com/antipratik", "vimeo.com",
			"Short films and visual essays shot in Nepal. Higher quality than YouTube for the cinematic work.", 0, "video"},
		{"link-ext-007", "@antipratik on X", "https://x.com/antipratik", "x.com",
			"Sporadic thoughts on music, code, and Kathmandu. The short-form version of everything else.", 0, "social"},
		{"link-ext-008", "GitHub", "https://github.com/antipratik", "github.com",
			"Open source code. Tools, utilities, and the occasional library. Most of it is small and useful.", 0, "social"},
	}

	for _, l := range links {
		if _, err := tx.Exec(
			"INSERT INTO links (id, title, url, domain, description, featured, category) VALUES (?, ?, ?, ?, ?, ?, ?)",
			l.id, l.title, l.url, l.domain, l.description, l.featured, l.category,
		); err != nil {
			return err
		}
	}
	return nil
}
