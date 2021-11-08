package main

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	"github.com/pomerium/verify"
)

func main() {
	bindAddress := os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = ":8080"
	}

	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ps := []string{
			path.Join("dist", r.URL.Path, "index.html"),
			path.Join("dist", r.URL.Path),
		}

		for _, p := range ps {
			f, err := verify.FS.Open(p)
			if err != nil {
				continue
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// ignore directories
			if fi.IsDir() {
				continue
			}

			http.ServeContent(w, r, path.Base(r.URL.Path), fi.ModTime(), f.(io.ReadSeeker))
			return
		}

		http.NotFound(w, r)
	})

	log.Info().Str("bind-address", bindAddress).Msg("starting http listener")
	http.ListenAndServe(bindAddress, r)
}
