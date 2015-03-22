// add_alias adds an email alias to a Google Apps user

// The intent is to allow you to use the Groups feature to have an unlimited
// number of email aliases which deliver mail to the same account.
//
// The defaults attempt to be comprehensive, as I typically use it from my
// phone via ConnectBot.  I use an alias of `aa` for faster input, eg:
// aa <newaddr>
// ... and voila!
//
// In the event you want to conveniently allow multiple users on the same
// machine to use it to add aliases to the same domain, you may wish to change
// the defaultDomain constant.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/asjoyner/googoauth"
	directory "github.com/google/google-api-go-client/admin/directory_v1"
)

var (
	defaultDomain = "joyner.ws"
	id     = "832157669857-it75mreuujhlm24iktv9927gakvnm1ni.apps.googleusercontent.com"
	secret = "i1wjHCLLhCXhxD9nZJbnT0Hu"
	scope  = []string{directory.AdminDirectoryGroupScope}
)

func addAlias(service *directory.Service, alias, domain, target, additional string) error {
	groupName := alias + "@" + domain
	g := &directory.Group{
		Name:        alias,
		Description: alias + " created programatically",
		Email:       groupName,
	}

	// Create a group with the name of the alias
	_, err := directory.NewGroupsService(service).Insert(g).Do()
	if err != nil {
		return fmt.Errorf("failed to create group %s: %s", alias, err)
	}

	// Add a user to it
	m := &directory.Member{
		Kind:  "admin#directory#member",
		Role:  "OWNER",
		Type:  "USER",
		Email: target,
	}
	_, err = directory.NewMembersService(service).Insert(groupName, m).Do()
	if err != nil {
		return fmt.Errorf("failed to add %s to group %s: %s", target, alias, err)
	}

	// Add any additional addresses
	m.Role = "MEMBER"
	for _, add := range strings.Split(additional, ",") {
		if add == "" {
			continue
		}
		m.Email = add
		_, err = directory.NewMembersService(service).Insert(groupName, m).Do()
		if err != nil {
			return fmt.Errorf("failed to add %s to group %s: %s", add, alias, err)
		}
	}
	return nil
}

func usage() {
	fmt.Printf("usage: %v <alias> [<alias>..]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var username string
	if u, err := user.Current(); err == nil {
		username = u.Username + "@" + defaultDomain
	}
	target := flag.String("target", username, "email address to add to the group")
	additional := flag.String("add", "", "comma separated list of additional email address to subscribe to the group")
	domain := flag.String("domain", defaultDomain, "domain to create the group in")

	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	// Setup the oauth connection to the admin/directory service
	c := googoauth.Client(id, secret, scope)
	service, err := directory.New(c)
	if err != nil {
		fmt.Printf("failed to make directory service connection: %s", err)
		os.Exit(1)
	}

	for _, alias := range flag.Args() {
		if err := addAlias(service, alias, *domain, *target, *additional); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("added alias: ", alias)
	}
}
