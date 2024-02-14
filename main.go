package main

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var (
	port = "8080"

	//go:embed assets
	assets embed.FS
)

func main() {
	http.Handle("GET /assets/", http.FileServerFS(assets))

	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("new request")
		t, err := template.New("landing").Parse(landing)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Whoeps!", http.StatusInternalServerError)
		}
		t.Execute(w, nil)
	})

	slog.Info("server started...", "port", port)

	if err := http.ListenAndServe(":"+port, http.DefaultServeMux); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

const landing = `
<!DOCTYPE html>
  <html lang="en" style="height: 100%;">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="Creating digital planners, digital stickers, and digital notebooks for platforms like GoodNotes" />

    <meta property="og:title" content="Four Rabbits Studio" />
    <meta property="og:url" content="https://www.fourrabbitsstudio.com/" />
    <meta property="og:description" content="Creating digital planners, digital stickers, and digital notebooks for platforms like GoodNotes" />
    <meta property="og:type" content="website" />
    <meta property="og:locale" content="en_US" />

    <meta property="og:image" content="http://www.fourrabbitsstudio.com/assets/logo.png" />
    <meta property="og:image:secure_url" content="http://www.fourrabbitsstudio.com/assets/logo.png" />
    <meta property="og:image:type" content="image/png" />
    <meta property="og:image:width" content="400" />
    <meta property="og:image:height" content="300" />
    <meta property="og:image:alt" content="A shiny red apple with a bite taken out" />

    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Four Rabbits Studio</title>
    <link rel="icon" href="/assets/favicon.ico" type="image/x-icon">
  </head>
  <body style="margin: 0 2rem; height: 100%;">
  <main style="height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center;">
        <img style="max-width: 100%;"src="/assets/logo.svg" alt="logo">
        <div style="margin: 6rem 0; display: flex; flex-direction: row; justify-content: center; align-items: center; gap: 4rem;">
            <a href="https://fourrabbitsstudio.etsy.com/"><img style="height: 3rem;" src="/assets/etsy.svg" alt="etsy"></img></a>
            <a href="https://www.instagram.com/"><img style="height: 3rem;" src="/assets/instagram.svg" alt="instagram"></img></a>
            <a href="https://www.tiktok.com/"><img style="height: 3rem;" src="/assets/tiktok.svg" alt="tiktok"></img></a>
        </div>
    </main>
  </body>
</html>
`
