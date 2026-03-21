// Command site serves lovyou.ai — the home of the hive's products.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lovyou-ai/site/content"
	"github.com/lovyou-ai/site/views"
)

func main() {
	port := flag.String("port", "", "HTTP port (default: $PORT or 8080)")
	flag.Parse()

	p := *port
	if p == "" {
		p = os.Getenv("PORT")
	}
	if p == "" {
		p = "8080"
	}

	posts, err := content.LoadPosts()
	if err != nil {
		log.Fatalf("load posts: %v", err)
	}
	log.Printf("loaded %d blog posts", len(posts))

	handleHome, handleBlogIndex, handleBlogPost := makeHandlers(posts)

	mux := http.NewServeMux()

	// Static files.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Pages.
	mux.HandleFunc("GET /{$}", handleHome)
	mux.HandleFunc("GET /blog", handleBlogIndex)
	mux.HandleFunc("GET /blog/{slug}", handleBlogPost)

	// Health check for Fly.io.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	addr := ":" + p
	log.Printf("lovyou.ai listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// ────────────────────────────────────────────────────────────────────
// Handlers
// ────────────────────────────────────────────────────────────────────

func makeHandlers(posts []views.Post) (home, blogIndex, blogPost http.HandlerFunc) {
	home = func(w http.ResponseWriter, r *http.Request) {
		views.Home().Render(r.Context(), w)
	}
	blogIndex = func(w http.ResponseWriter, r *http.Request) {
		views.BlogIndex(posts).Render(r.Context(), w)
	}
	blogPost = func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		for _, post := range posts {
			if post.Slug == slug {
				views.BlogPost(post).Render(r.Context(), w)
				return
			}
		}
		http.NotFound(w, r)
	}
	return
}
