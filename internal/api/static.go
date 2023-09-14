package api

import (
	"github.com/AlexxIT/go2rtc/www"
	"net/http"
)

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set your desired headers here
		w.Header().Set("Access-Control-Allow-Origin", "*")  // allows any domain, for security you might specify allowed domains
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		// If it's a preflight OPTIONS request, respond no content
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		// Continue the chain of handlers
		next.ServeHTTP(w, r)
	})
}


func initStatic(staticDir string) {
	var root http.FileSystem
	if staticDir != "" {
		log.Info().Str("dir", staticDir).Msg("[api] serve static")
		root = http.Dir(staticDir)
	} else {
		root = http.FS(www.Static)
	}

	base := len(basePath)
	fileServer := http.FileServer(root)

	// Using the CORS middleware with your static file server
	corsFileServer := withCORS(fileServer)

	HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		if base > 0 {
			r.URL.Path = r.URL.Path[base:]
		}
		corsFileServer.ServeHTTP(w, r)  // Use the cors wrapped handler here
	})
}