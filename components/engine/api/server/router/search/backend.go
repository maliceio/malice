package search

import "github.com/maliceio/engine/api/types/search"

// Backend for Plugin
type Backend interface {
	Search(hash string, config *search.Config) (search.Result, error)
}
