package view

import (
	"github.com/go-iris2/iris2/template/django"
)

type (
	// FilterFunction conversion for django.FilterFunction
	FilterFunction = django.FilterFunction
)

// DjangoAdaptor is the  adaptor for the Django engine.
// Read more about the Django Go Template at:
// https://github.com/flosch/pongo2
// and https://github.com/go-iris2/iris2/template/tree/master/django
type DjangoAdaptor struct {
	*Adaptor
	engine  *django.Engine
	filters map[string]FilterFunction
}

// Django returns a new kataras/go-template/django template engine
// with the same features as all iris' view engines have:
// Binary assets load (templates inside your executable with .go extension)
// Layout, Funcs, {{ url }} {{ urlpath}} for reverse routing and much more.
//
// Read more: https://github.com/flosch/pongo2
func Django(directory string, extension string) *DjangoAdaptor {
	e := django.New()
	return &DjangoAdaptor{
		Adaptor: NewAdaptor(directory, extension, e),
		engine:  e,
		filters: make(map[string]FilterFunction, 0),
	}
}

// Filters for pongo2, map[name of the filter] the filter function .
//
// Note, these Filters function overrides ALL the previous filters
// It SETS a new filter map based on the given 'filtersMap' parameter.
func (d *DjangoAdaptor) Filters(filtersMap map[string]FilterFunction) *DjangoAdaptor {
	if len(filtersMap) == 0 {
		return d
	}
	// configuration maps are never nil, because
	// they are initialized at each of the engine's New func
	// so we're just passing them inside it.
	d.engine.Config.Filters = filtersMap
	return d
}

// Globals share context fields between templates. https://github.com/flosch/pongo2/issues/35
func (d *DjangoAdaptor) Globals(globalsMap map[string]interface{}) *DjangoAdaptor {
	if len(globalsMap) == 0 {
		return d
	}

	for k, v := range globalsMap {
		d.engine.Config.Globals[k] = v
	}

	return d
}

// DebugTemplates enables template debugging.
// The verbose error messages will appear in browser instead of quiet passes with error code.
func (d *DjangoAdaptor) DebugTemplates(debug bool) *DjangoAdaptor {
	d.engine.Config.DebugTemplates = debug
	return d
}
