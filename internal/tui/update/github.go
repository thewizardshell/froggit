package update

import "froggit/internal/gh"

var ghClient = gh.NewGhClient()

func GetGhClient() *gh.GhClient {
	return ghClient
}
