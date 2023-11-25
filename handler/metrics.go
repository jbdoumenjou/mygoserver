package handler

import (
	"fmt"
	"net/http"
)

type Metrics struct {
	fileserverHits int
}

func (m *Metrics) MetricsTextHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + fmt.Sprintf("%d", m.fileserverHits)))
}

func (m *Metrics) MetricsHTMLHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	const html = `<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
`
	w.Write([]byte(fmt.Sprintf(html, m.fileserverHits)))
}

func (m *Metrics) ResetHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	m.fileserverHits = 0
}

func (m *Metrics) MiddlewareInc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.fileserverHits++

		next.ServeHTTP(w, r)
	}
}
