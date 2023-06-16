package scribble

type TokenizerStateTypes string

const (
	STATE_NONE            TokenizerStateTypes = "none"
	STATE_START_OF_STREAM TokenizerStateTypes = "start_stream"
	STATE_IN_TOKEN        TokenizerStateTypes = "in_token"
	STATE_IN_BOUNDARY     TokenizerStateTypes = "in_bounds"
	STATE_WRITE_TOKEN     TokenizerStateTypes = "write_token"
	STATE_END_OF_STREAM   TokenizerStateTypes = "end_stream"
)
