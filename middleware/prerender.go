package middleware

import (
	"io"
	"log"
	"net/http"
	"strings"
)

var crawlerUserAgents = []string{
	"googlebot", "bingbot", "yandex", "baiduspider", "facebookexternalhit",
	"twitterbot", "linkedinbot", "embedly", "slackbot", "discordbot",
}

func shouldPrerender(r *http.Request) bool {
	ua := strings.ToLower(r.Header.Get("User-Agent"))
	for _, bot := range crawlerUserAgents {
		if strings.Contains(ua, bot) {
			return true
		}
	}
	return r.URL.Query().Has("_escaped_fragment_")
}

func PrerenderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldPrerender(r) && r.Method == http.MethodGet {
			prerenderUrl := "https://service.prerender.io" + r.URL.RequestURI()

			req, _ := http.NewRequest("GET", prerenderUrl, nil)
			req.Header.Set("User-Agent", r.Header.Get("User-Agent"))
			req.Header.Set("X-Prerender-Token", "JTv89Qhqb1AdhaDBRDj9") // replace this

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error calling Prerender.io:", err)
				next.ServeHTTP(w, r)
				return
			}
			defer resp.Body.Close()

			for k, v := range resp.Header {
				w.Header().Set(k, v[0])
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}
		next.ServeHTTP(w, r)
	})
}
