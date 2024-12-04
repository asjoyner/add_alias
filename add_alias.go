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
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/asjoyner/googoauth"
	directory "google.golang.org/api/admin/directory/v1"
	groupssettings "google.golang.org/api/groupssettings/v1"
	"google.golang.org/api/option"
)

var (
	defaultDomain = "joyner.ws"
	id            = "832157669857-it75mreuujhlm24iktv9927gakvnm1ni.apps.googleusercontent.com"
	secret        = "i1wjHCLLhCXhxD9nZJbnT0Hu"
	scope         = []string{
		directory.AdminDirectoryGroupScope,     // to create a group, and add members
		groupssettings.AppsGroupsSettingsScope, // to set the isArchived bit
	}
)

func addAlias(c *http.Client, alias, domain, target, additional string) error {
	ctx := context.Background()
	directoryService, err := directory.NewService(ctx, option.WithHTTPClient(c))
	if err != nil {
		fmt.Printf("failed to make directoryService connection: %s", err)
		os.Exit(1)
	}
	groupName := alias + "@" + domain
	g := &directory.Group{
		Name:        alias,
		Description: alias + " created programatically",
		Email:       groupName,
	}

	// Create a group with the name of the alias
	group, err := directory.NewGroupsService(directoryService).Insert(g).Do()
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
	_, err = directory.NewMembersService(directoryService).Insert(group.Id, m).Do()
	if err != nil {
		return fmt.Errorf("failed to add %s to group %s: %s", target, alias, err)
	}

	// Configure the default settings for the group
	groupssettingsService, err := groupssettings.NewService(ctx, option.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to initialize groupssettingsService: %s", err)
	}
	preferredDefaults := &groupssettings.Groups{
		// Store conversation history.
		IsArchived: "true",
		// Turn off moderation.  (This is the default, but the docs don't seem to
		// be happy about it, so I set it out of an abundance of caution.)
		MessageModerationLevel: "MODERATE_NONE",
		// Turn off spam filtering at the group level (letting GMail do that).
		SpamModerationLevel: "ALLOW",
	}
	groupsService := groupssettings.NewGroupsService(groupssettingsService)
	if _, err := groupsService.Update(groupName, preferredDefaults).Do(); err != nil {
		return fmt.Errorf("failed to request the group (%s) be archived: %s", group.Id, err)
	}

	// Add any additional addresses
	m.Role = "MEMBER"
	for _, add := range strings.Split(additional, ",") {
		if add == "" {
			continue
		}
		m.Email = add
		_, err = directory.NewMembersService(directoryService).Insert(groupName, m).Do()
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

	for _, alias := range flag.Args() {
		if err := addAlias(c, alias, *domain, *target, *additional); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("added alias: ", alias)
	}
}
