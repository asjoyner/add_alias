// add_alias adds an alias, in the form of a Google Group, to a Google Apps user

// The intent is to allow you to use the Groups feature to have an unlimited
// number of email aliases which deliver mail to the same account.

package main

import (
  "fmt"

  "github.com/google/google-api-go-client/admin/directory_v1"
)

// TODO: import and use flag
var domain := "joyner.ws"
var user := "asjoyner"  // Get from env

//fmt.Println("usage: %v <alias> [<alias>..]", sys.argv[0])

func main() {
  client := getOauthClient()
}
