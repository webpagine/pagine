package environment

import (
	"github.com/webpagine/pagine/site"
	"net/http"
)

func Serve(site *site.Site, publicDir string) error {

	// TODO

	err := http.ListenAndServe("/", http.FileServer(http.Dir(publicDir)))
	if err != nil {
		return err
	}

	return nil
}
