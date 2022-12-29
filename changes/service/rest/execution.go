package service

func valueFromPath(name string, pathVariables map[string]string) (string, bool) {
	value, ok := pathVariables[name]
	return value, ok
}
