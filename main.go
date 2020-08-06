package main

import (
	"context"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type Configuration struct {
	QuireAuthToken string
	ClientId       string
	ClientSecret   string
}

func main() {

	httpclient.Defaults(httpclient.Map{
		"Content-Type": "",
	})

	app := &cli.App{
		Name:        "quire-cli",
		Description: "Lo Agency Quire CLI",
		Authors: []*cli.Author{
			{
				Name:  "Aien Saidi",
				Email: "aien@lo.agency",
			},
		},
		Version: "0.0.1",
		Commands: []*cli.Command{
			GitCommand,
			{
				Name:  "authorize",
				Usage: "Authorize the app to use Quire Boards",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "clientId",
						Usage: "Quire client id",
					},
					&cli.StringFlag{
						Name:  "clientSecret",
						Usage: "Quire client secret",
					},
				},
				Action: func(c *cli.Context) error {
					var err error
					var cmd *exec.Cmd

					var url = fmt.Sprintf("https://quire.io/oauth?client_id=%s&redirect_uri=http://localhost:1992", c.String("clientId"))

					switch runtime.GOOS {
					case "linux":
						cmd = exec.Command("xdg-open", url)
					case "windows":
						cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
					case "darwin":
						cmd = exec.Command("open", url)
					default:
						return fmt.Errorf("unsupported platform")
					}

					err = cmd.Start()
					if err != nil {
						return err
					}

					m := http.NewServeMux()
					s := http.Server{Addr: ":1992", Handler: m}

					m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						keys, ok := r.URL.Query()["code"]

						if !ok || len(keys[0]) < 1 {
							_, _ = fmt.Fprintf(w, "url param 'code' is missing")
							return
						}

						// Query()["key"] will return an array of items,
						// we only want the single item.
						key := keys[0]

						configuration := Configuration{
							QuireAuthToken: key,
							ClientSecret:   c.String("clientSecret"),
							ClientId:       c.String("clientId"),
						}

						err := SaveConfig(configuration)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							_, _ = w.Write([]byte(err.Error()))
							return
						}

						log.Printf("configurations are saved, you can now use the extension!\ntoken: %s\n", key)

						_ = s.Shutdown(context.Background())

						return
					})

					if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
