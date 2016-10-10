[![CircleCI](https://circleci.com/gh/ikaros/go-reddit.svg?style=svg)](https://circleci.com/gh/ikaros/go-reddit) [![Build Status](https://drone.io/github.com/ikaros/go-reddit/status.png)](https://drone.io/github.com/ikaros/go-reddit/latest) [![Build Status](https://travis-ci.org/ikaros/go-reddit.svg?branch=master)](https://travis-ci.org/ikaros/go-reddit) [![GoDoc](https://godoc.org/github.com/ikaros/go-reddit/reddit?status.svg)](https://godoc.org/github.com/ikaros/go-reddit/reddit) [![Go Report Card](https://goreportcard.com/badge/github.com/ikaros/go-reddit)](https://goreportcard.com/report/github.com/ikaros/go-reddit) [![Coverage Status](https://coveralls.io/repos/github/ikaros/go-reddit/badge.svg?branch=master)](https://coveralls.io/github/ikaros/go-reddit?branch=master)
# ![go-reddit go-reddit - Unofficial reddit API client for Go](docs/static/goredditlogo_header.png "go-reddit - Unofficial reddit API client for Go")
**go-reddit** is a **unofficial Go client package** for the reddit API.
It covers JSON encoding/decoding, rate limiting and errors.
OAUTH authentication should be done with the [oauth2 package](https://github.com/golang/oauth2).

**Currently the Go API is not finished and quite unstable**

# Goals
- Cover full API surface and functionality of the reddit API
- Handle parsing of responses, so there is no need to deal with raw JSON
- Implement circuit breaking for reddit's rate limiting

# Non-goals
- Reimplement OAUTH authentication

# Additional reading
- [Documentation](https://www.reddit.com/dev/api/) – reddit API Documentation
- [Rules](https://github.com/reddit/reddit/wiki/API) – reddit API Rules
