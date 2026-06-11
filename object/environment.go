package object

type Environment struct {
	store map[string]Object
	outer *Environment // the origin of this extended environment (outer scope)
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Try to get object named "name".
//
// Returns object and true If it can get the object associated with the name.
//
// Else returns nil and false.
//
// Try to get object from inner scope this function called to outer scopes that enclose recursively.
func (e *Environment) Get(name string) (Object, bool) {
	// search in this inner scope first
	obj, ok := e.store[name]
	// or search in the outer scope recursively
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Set new object (val) with keyword (name).
// Store the key-value in this scope where Environment.Set() is called.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Create new an enclosed inner scope.
//   - 外側のスコープは内側のスコープを包み込む
//   - 内側のスコープは外側のスコープを拡張する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
