package twitter_search

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// SearchResult twitter search api result
type SearchResult struct {
	Data []struct {
		AuthorID  string `json:"author_id"`
		ID        string `json:"id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
}

// Search for twitter
func Search(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	query := "#ストVラウンジ募集 -is:retweet"
	URL := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s&tweet.fields=text,created_at&expansions=author_id&user.fields=username&max_results=50", url.QueryEscape(query))

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)

	if err != nil {
		log.Fatalf("failed to initialize request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getSecret()))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code is %d.", resp.StatusCode)
	}

	b := []byte{}
	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("failed to read response: %s", err.Error())
	}

	result := SearchResult{}
	err = json.Unmarshal(b, &result)

	if err != nil {
		log.Fatalf("failed to unmarshal: %s", err.Error())
	}
}

func getSecret() string {
	name := fmt.Sprintf("projects/%s/secrets/%s", os.Getenv("PROJECT_ID"), os.Getenv("SECRET_NAME"))

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create secretmanager client: %v", err)
	}

	// Build the request.
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetSecret(ctx, req)
	if err != nil {
		log.Fatalf("failed to get secret: %v", err)
	}
	return result.Name
}
