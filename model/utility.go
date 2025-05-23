package model

type ReadAndConvert struct {
	EncodedContent string
	Err            error
}

type ChromeHandler struct {
	Status string
	Err    error
}
