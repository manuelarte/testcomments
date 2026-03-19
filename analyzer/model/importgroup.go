package model

import "go/ast"

// ImportGroup contains the imports that are important for the test.
type ImportGroup struct {
	// goCmp "go-cmp" import spec. Nil if "go-cmp" is not imported.
	goCmp *ast.ImportSpec
	// reflect "reflect" import spec. Nil if "reflect" is not imported.
	reflect *ast.ImportSpec
	// testing "testing" import spec.
	testing *ast.ImportSpec
}

func (i ImportGroup) NewWithImportSpec(is *ast.ImportSpec) ImportGroup {
	if is == nil {
		return i
	}

	switch is.Path.Value {
	case "\"github.com/google/go-cmp/cmp\"":
		i.goCmp = is
	case "\"reflect\"":
		i.reflect = is
	case "\"testing\"":
		i.testing = is
	}

	return i
}

func (i ImportGroup) GoCmpImportName() (string, bool) {
	if i.goCmp == nil {
		return "", false
	}

	return importName(i.goCmp), true
}

func (i ImportGroup) ReflectImportName() (string, bool) {
	if i.reflect == nil {
		return "", false
	}

	return importName(i.reflect), true
}

func (i ImportGroup) TestingImportName() (string, bool) {
	if i.testing == nil {
		return "", false
	}

	return importName(i.testing), true
}
