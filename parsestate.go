package scribble

type ParseState struct {
	CurrentState string
	Trigger      string
	NextState    string
}
