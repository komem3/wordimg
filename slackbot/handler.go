package slackbot

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/compute/metadata"
	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	"github.com/slack-go/slack"
)

func init() {
	ctx := context.Background()
	pnum, err := metadata.NumericProjectID()
	if err != nil {
		panic(err)
	}
	projectNum = pnum

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	secretClient = client
}

func SlackEmojiGen(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wordimgURL := os.Getenv("URL")
	slackSercret, err := getSecret(ctx, "SLACK_SECRET")
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slackAPIKey, err := getSecret(ctx, "SLACK_APIKEY")
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	verifier, err := slack.NewSecretsVerifier(r.Header, slackSercret)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = verifier.Ensure(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	url, err := url.Parse(wordimgURL)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	query := url.Query()
	query.Add("text", s.Text)

	switch s.Command {
	case "/wordimg1":
		query.Add("size", "j1")
	case "/wordimg2":
		query.Add("size", "j2")
	case "/wordimg3":
		query.Add("size", "j3")
	case "/wordimg":
		// pass
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] command is %s\n", s.Command)
		return
	}
	log.Printf("[DEBUG] query %#v, command %#v\n", query, s)
	url.RawQuery = query.Encode()

	api := slack.New(slackAPIKey)
	channelID, _, err := api.PostMessage(
		s.ChannelID,
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionAttachments(slack.Attachment{
			Text:     "絵文字を生成しました。",
			ImageURL: url.String(),
		}),
	)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("[INFO] write channel %s, text %s\n", channelID, url)
	w.WriteHeader(http.StatusNoContent)
}
