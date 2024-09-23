package object

// Environment represents the scope of a program.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new global environment.
func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

// NewEnclosedEnvironment creates a new enclosed environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: outer}
}

// Get gets the value associated with the identifier.
func (e *Environment) Get(identifier string) (Object, bool) {
	obj, ok := e.store[identifier]
	// If the identifier does not exist in the current environment, continously search all other nested environements
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(identifier)
	}
	return obj, ok
}

// Set stores a value in the environment for the given identifier.
func (e *Environment) Set(identifier string, value Object) Object {
	e.store[identifier] = value
	return value
}
