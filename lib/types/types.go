package types
//package libs/types/ defines the type keywords of the lang, as well as defines the struct used for tracking Types

//TODO: eventually algorithmically fill the type keywords to allow user defined types
//	FAR stretch goal though

var TypeKeywords = map[string]struct{} { //Go doesn't have a set, so they are usually done like this. :(
	"int32" : {},
	"int64" : {},
}

type TypeStruct struct {
	Type string
}
