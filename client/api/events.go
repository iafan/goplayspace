package api

// CompileEvent is the individual event structure
// as returned by play.golang.org
type CompileEvent struct {
	Message string
	Kind    string
	Delay   int
}

// CompileResponse is the entire /compile response payload structure
// as returned by play.golang.org
type CompileResponse struct {
	Body   *string
	Events []*CompileEvent
	Errors string
}
