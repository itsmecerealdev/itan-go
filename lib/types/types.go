package types
//package libs/types/ defines the type keywords of the lang, as well as defines the struct used for tracking Types

var TypeKeywords = map[string]struct{} { //Go doesn't have a set, so they are usually done like this. :(
	"int32" : {},
	"int64" : {},
}

type TypeStruct struct {
	Type string
}
