package scribble

type MatchTypes string

const (
	MATCH_NONE     MatchTypes = "none"
	MATCH_BOUNDARY MatchTypes = "bounds"
	MATCH_TOKEN    MatchTypes = "token"
)
