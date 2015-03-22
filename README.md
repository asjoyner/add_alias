# add\_alias
Add an alias, in the form of a Google Group, to a Google Apps user

The intent is to allow you to use the Groups feature to have an unlimited
number of email aliases which deliver mail to the same account.

usage: ./add\_alias <alias> [<alias>..]
  -add="": comma separated list of additional email address to subscribe to the group
  -authport="12345": HTTP Server port.  Only needed for the first run, your browser will send credentials here.  Must be accessible to your browser, and authorized in the developer console.
  -debug.http=false: show HTTP traffic
  -domain="joyner.ws": domain to create the group in
  -target="asjoyner@joyner.ws": email address to add to the group

The first time you run the command, it will load a URL which will prompt you to login if necessary, and authorize the app to access the Groups service for your domain.  The Google webpage responds by redirecting you with the authorization code to http://localhost:12345.  add\_alias will listen on that port, but if you are running add\_alias on a remote server, you will need to port-forward 12345 from your local host to 12345 on the remote server.  add\_alias saves the authentication credentials to disk, so that subsequent invocations do not require you to reauthenticate.
