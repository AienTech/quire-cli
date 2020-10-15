package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	urllib "net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/ddliu/go-httpclient"
	"github.com/urfave/cli/v2"
)

type Configuration struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

var ClientId string
var ClientSecret string

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
		Version: "0.0.2",
		Commands: []*cli.Command{
			GitCommand,
			{
				Name:  "authorize",
				Usage: "Authorize the app to use Quire Boards",
				Action: func(c *cli.Context) error {
					var err error
					var cmd *exec.Cmd

					var url = fmt.Sprintf("https://quire.io/oauth?client_id=%s&redirect_uri=http://localhost:1992", ClientId)

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

						formData := urllib.Values{
							"grant_type":    {"authorization_code"},
							"code":          {key},
							"client_id":     {ClientId},
							"client_secret": {ClientSecret},
						}

						resp, err := http.PostForm("https://quire.io/oauth/token", formData)
						if err != nil {
							log.Fatalln(err)
						}

						var result Configuration

						/*
							{
								"access_token":"ACCESS_TOKEN",
								"token_type":"bearer",
								"expires_in":2592000,
								"refresh_token":"REFRESH_TOKEN"
							}
						*/
						json.NewDecoder(resp.Body).Decode(&result)

						err = SaveConfig(result)
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
