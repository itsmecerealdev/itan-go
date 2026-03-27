package types

var TypeKeywords = map[string]struct{} { //Go doesn't have a set, so they are usually done like this. :(
	"int32",
	"int64",
}

type TypeStruct struct {
	Type string
}
