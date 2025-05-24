package main

type Server struct {
	readabilityParser *ReadabilityParser
}

func newServer() (*Server, error) {
	parser, err := NewReadabilityParser()
	if err != nil {
		return nil, err
	}
	return &Server{
		readabilityParser: parser,
	}, nil
}
