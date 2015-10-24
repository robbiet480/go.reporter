# go.reporter
go.reporter is a Golang library for parsing [Reporter-App](http://www.reporter-app.com/) JSON files.

# Features
* Full support for all fields in all JSON versions.
* Supports both version of the JSON schema.
* Allows reading JSON from a string, the local filesystem, or Dropbox.

# Getting started
```
import "github.com/robbiet480/go.reporter"
```

Check out the examples in [`example_test.go`](example_test.go).

To use this library with Dropbox, you will need to make a [new Dropbox app](https://www.dropbox.com/developers-v1/apps/create).
You must set the permissions to allow full Dropbox access.
If you want, you can further limit access from the created app by only allowing it access to text files (JSON is covered in that category).
Once you have done this, [follow these instructions](https://www.dropbox.com/developers-v1/reference/oauthguide#testing-with-a-generated-access-token)
to generate an access token for your own account.

# Compatibility Notes
This library provides compatibility with both versions of the Reporter JSON schema. The differences that I have noticed are:

1. Timestamps are expressed as seconds since Apple Epoch (January 1st, 2001, 00:00:00 UTC)
2. There were no `uniqueIdentifiers` anywhere.
3. Some metadata variables are missing in version two (the latest), such as `dwellStatus` and `sync`.

# Other resources
* [The schema gist](https://gist.github.com/dbreunig/9315705)
* [Example data](https://reporter.zendesk.com/hc/en-us/articles/200273009-How-DropBox-works-with-Reporter)

# Tests
`go test`