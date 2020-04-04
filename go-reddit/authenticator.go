package reddit

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// Authenticator provides functions for authenticating a user via OAuth2 and generating a client that can be used to access authorized API endpoints.
type Authenticator struct {
	config                *oauth2.Config
	state                 string
	userAgent             string
	RequestPermanentToken bool
}

const (
	// ScopeIdentity allows access to account information.
	ScopeIdentity = "identity"
	// ScopeEdit allows modification and deletion of comments and submissions.
	ScopeEdit = "edit"
	// ScopeFlair allows modification of user link flair on submissions.
	ScopeFlair = "flair"
	// ScopeHistory allows access to user voting history on comments and submissions
	ScopeHistory = "history"
	// ScopeModConfig allows management of configuration, sidebar, and CSS of user managed subreddits.
	ScopeModConfig = "modconfig"
	// ScopeModFlair allows management and assignment of user moderated subreddits.
	ScopeModFlair = "modflair"
	// ScopeModLog allows access to moderation log for user moderated subreddits.
	ScopeModLog = "modlog"
	// ScopeModWiki allows changing of editors and visibility of wiki pages in user moderated subreddits.
	ScopeModWiki = "modwiki"
	// ScopeMySubreddits allows access to the list of subreddits user moderates, contributes to, and is subscribed to.
	ScopeMySubreddits = "mysubreddits"
	// ScopePrivateMessages allows access to user inbox and the sending of private messages to other users.
	ScopePrivateMessages = "privatemessages"
	// ScopeRead allows access to user posts and comments.
	ScopeRead = "read"
	// ScopeReport allows reporting of content for rules violations.
	ScopeReport = "report"
	// ScopeSave allows saving and unsaving of user comments and submissions.
	ScopeSave = "save"
	// ScopeSubmit allows user submission of links and comments.
	ScopeSubmit = "submit"
	// ScopeSubscribe allows management of user subreddit subscriptions and friends.
	ScopeSubscribe = "subscribe"
	// ScopeVote allows user submission and changing of votes on comments and submissions.
	ScopeVote = "vote"
	// ScopeWikiEdit allows user editing of wiki pages.
	ScopeWikiEdit = "wikiedit"
	// ScopeWikiRead allow user viewing of wiki pages.
	ScopeWikiRead = "wikiread"

	authURL  = "https://www.reddit.com/api/v1/authorize"
	tokenURL = "https://www.reddit.com/api/v1/access_token"
)

// NewAuthenticator generates a new authenticator with the supplied client, state, and requested scopes.
func NewAuthenticator(clientID string, clientSecret string, redirectURL string, userAgent string, state string, scopes ...string) *Authenticator {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURL,
	}

	return &Authenticator{
		config:    config,
		state:     state,
		userAgent: userAgent,
	}
}

// GetAuthenticationURL retrieves the URL used to direct the authenticating user to Reddit for permissions approval.
func (a *Authenticator) GetAuthenticationURL() string {
	url := a.config.AuthCodeURL(a.state)
	if a.RequestPermanentToken {
		url += "&duration=permanent"
	}
	return url
}

type uaSetterTransport struct {
	config    *oauth2.Config
	userAgent string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (t *uaSetterTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	// set a non-standard Authorization header because reddit demands it
	// https://github.com/reddit/reddit/wiki/OAuth2#retrieving-the-access-token
	req.Header.Set("Authorization", basicAuth(t.config.ClientID, t.config.ClientSecret))
	return http.DefaultTransport.RoundTrip(req)
}

// GetToken exchanges an authorization code for an access token.
func (a *Authenticator) GetToken(state string, code string) (*oauth2.Token, error) {
	if state != a.state {
		return nil, errors.New("Invalid state")
	}

	// Construct a custom http client that forces the user-agent and attach it
	// to the oauth2 context. https://github.com/golang/oauth2/issues/179
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: a.config.TokenSource(oauth2.NoContext, &oauth2.Token{
				AccessToken: code,
			}),
			Base: &uaSetterTransport{
				config:    a.config,
				userAgent: a.userAgent,
			},
		},
	}
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, client)

	return a.config.Exchange(ctx, code)
}

// GetAuthClient generates a new authenticated client using the supplied access token.
func (a *Authenticator) GetAuthClient(token *oauth2.Token) *Client {
	return &Client{
		http:      a.config.Client(oauth2.NoContext, token),
		userAgent: a.userAgent,
	}
}
