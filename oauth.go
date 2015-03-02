// oauth contains helper functions for creating an oauth session

// Originally based on the oauth example:
// http://godoc.org/golang.org/x/oauth2#example-Config

package main

import "golang.org/x/oauth2"

var defaultId :=     "832157669857-it75mreuujhlm24iktv9927gakvnm1ni.apps.googleusercontent.com"
var defaultSecret := "i1wjHCLLhCXhxD9nZJbnT0Hu"


func getOauthClient() {
  conf := &oauth2.Config{
      ClientID:     clientId,
      ClientSecret: clientSecret,
      Scopes:       []string{"SCOPE1", "SCOPE2"},
      Endpoint: oauth2.Endpoint{
          AuthURL:  "https://provider.com/o/oauth2/auth",
          TokenURL: "https://provider.com/o/oauth2/token",
      },
  }

  // Redirect user to consent page to ask for permission
  // for the scopes specified above.
  url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
  fmt.Printf("Visit the URL for the auth dialog: %v", url)

  // Use the authorization code that is pushed to the redirect URL.
  // NewTransportWithCode will do the handshake to retrieve
  // an access token and initiate a Transport that is
  // authorized and authenticated by the retrieved token.
  var code string
  if _, err := fmt.Scan(&code); err != nil {
      log.Fatal(err)
  }
  tok, err := conf.Exchange(oauth2.NoContext, code)
  if err != nil {
      log.Fatal(err)
  }

  client := conf.Client(oauth2.NoContext, tok)
  client.Get("...")
}
