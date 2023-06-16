package scribble

type ParseStateTransition struct {
	CurrentState string
	Token        string
	NextState    string
}
