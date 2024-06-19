package kind

type Kind string

const (
	Unknown    Kind = ""
	Executable Kind = "executable"
	Source     Kind = "source"
	Virtual    Kind = "virtual"
	Root       Kind = "root"
	Posix      Kind = "posix"
)
