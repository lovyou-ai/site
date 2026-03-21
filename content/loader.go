// Package content loads embedded blog posts from markdown files.
package content

import (
	"bytes"
	"embed"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lovyou-ai/site/views"
	"github.com/yuin/goldmark"
)

//go:embed posts/*.md
var postsFS embed.FS

var (
	monthYear = regexp.MustCompile(`·\s+(\w+)\s+(\d{4})`)
	postNum   = regexp.MustCompile(`^post(\d+)`)
)

// LoadPosts reads all embedded markdown posts and returns them newest-first.
func LoadPosts() ([]views.Post, error) {
	entries, err := postsFS.ReadDir("posts")
	if err != nil {
		return nil, fmt.Errorf("read posts dir: %w", err)
	}

	md := goldmark.New()
	var posts []views.Post

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		raw, err := postsFS.ReadFile("posts/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}

		post, err := parsePost(md, e.Name(), raw)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func parsePost(md goldmark.Markdown, filename string, raw []byte) (views.Post, error) {
	lines := strings.SplitN(string(raw), "\n", -1)

	// Title: first line starting with "# "
	var title string
	for _, l := range lines {
		if strings.HasPrefix(l, "# ") {
			title = strings.TrimPrefix(l, "# ")
			break
		}
	}

	// Summary: first line wrapped in *...*  (italic)
	var summary string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "*") && strings.HasSuffix(l, "*") && len(l) > 2 {
			summary = strings.Trim(l, "*")
			break
		}
	}
	// Fallback: use the first non-empty, non-title, non-author line
	if summary == "" {
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "---") || strings.Contains(l, "·") {
				continue
			}
			summary = l
			if len(summary) > 200 {
				summary = summary[:200] + "..."
			}
			break
		}
	}

	// Date: parse "· Month Year" from byline
	date := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	if m := monthYear.FindStringSubmatch(string(raw)); m != nil {
		year, _ := strconv.Atoi(m[2])
		month := parseMonth(m[1])
		day := 1
		// Use post number as day offset for ordering
		if nm := postNum.FindStringSubmatch(filename); nm != nil {
			n, _ := strconv.Atoi(nm[1])
			day = n
			if day > 28 {
				day = 28
			}
		}
		date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}

	// Slug: strip "postNN-" prefix and ".md" suffix
	slug := strings.TrimSuffix(filename, ".md")
	if idx := strings.Index(slug, "-"); idx > 0 {
		slug = slug[idx+1:]
	}

	// Body: convert full markdown to HTML (skip the header metadata lines)
	bodyStart := findBodyStart(lines)
	bodyMD := strings.Join(lines[bodyStart:], "\n")

	var buf bytes.Buffer
	if err := md.Convert([]byte(bodyMD), &buf); err != nil {
		return views.Post{}, fmt.Errorf("convert markdown: %w", err)
	}

	return views.Post{
		Slug:    slug,
		Title:   title,
		Summary: summary,
		Date:    date,
		Body:    buf.String(),
	}, nil
}

// findBodyStart skips past the title, subtitle, author, and first --- separator.
func findBodyStart(lines []string) int {
	pastFirstHR := false
	for i, l := range lines {
		l = strings.TrimSpace(l)
		if l == "---" {
			if pastFirstHR {
				// Second --- means start of body after it
				return i + 1
			}
			pastFirstHR = true
		}
	}
	// No separator found — start after title
	for i, l := range lines {
		if strings.HasPrefix(l, "# ") {
			return i + 1
		}
	}
	return 0
}

func parseMonth(s string) time.Month {
	months := map[string]time.Month{
		"January": time.January, "February": time.February, "March": time.March,
		"April": time.April, "May": time.May, "June": time.June,
		"July": time.July, "August": time.August, "September": time.September,
		"October": time.October, "November": time.November, "December": time.December,
	}
	if m, ok := months[s]; ok {
		return m
	}
	return time.January
}
