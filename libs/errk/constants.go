package errk

const (
	pkgNamespace = "errk"
)

var DuplicateFallbackError = NewError("E_COM_0", "Cannot create new Error that has same code with Fallback Error",
	WithNamespace(pkgNamespace))
