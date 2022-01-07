package parser

import (
	"go/format"
	"strings"

	template "github.com/timmypotts/gk/templates"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"
)

type ParsedSrc interface {
	String() string
}
type File struct {
	Comment    string
	Package    string
	Imports    []NamedTypeValue
	Constants  []NamedTypeValue
	Vars       []NamedTypeValue
	Interfaces []Interface
	Structs    []Struct
	Methods    []Method
}

type Struct struct {
	Name    string
	Comment string
	Vars    []NamedTypeValue
}
type Interface struct {
	Name    string
	Comment string
	Methods []Method
}
type Method struct {
	Comment    string
	Name       string
	Struct     NamedTypeValue
	Body       string
	Parameters []NamedTypeValue
	Results    []NamedTypeValue
}
type NamedTypeValue struct {
	Name     string
	Type     string
	Value    string
	HasValue bool
}

func (f *File) String() string {
	s, err := template.NewEngine().Execute("file", f)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := imports.Process(f.Package, []byte(s), nil)
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func (m *Method) String() string {
	str := ""
	if m.Struct.Name != "" {
		s, err := template.NewEngine().ExecuteString("{{template \"struct_function\" .}}", m)
		if err != nil {
			logrus.Panic(err)
		}
		str = s
	} else {
		s, err := template.NewEngine().ExecuteString("{{template \"func\" .}}", m)
		if err != nil {
			logrus.Panic(err)
		}
		str = s
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func (s *Struct) String() string {
	str, err := template.NewEngine().ExecuteString("{{template \"struct\" .}}", s)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}
func (i *Interface) String() string {
	str, err := template.NewEngine().ExecuteString("{{template \"interface\" .}}", i)
	if err != nil {
		logrus.Panic(err)
	}
	dt, err := format.Source([]byte(str))
	if err != nil {
		logrus.Panic(err)
	}
	return string(dt)
}

func NewNameType(name string, tp string) NamedTypeValue {
	return NamedTypeValue{
		Name:     name,
		Type:     tp,
		HasValue: false,
	}
}
func NewNameTypeValue(name string, tp string, vl string) NamedTypeValue {
	return NamedTypeValue{
		Name:     name,
		Type:     tp,
		HasValue: true,
		Value:    vl,
	}
}

func NewMethod(name string, str NamedTypeValue, body string, parameters, results []NamedTypeValue) Method {
	return Method{
		Name:       name,
		Comment:    "",
		Struct:     str,
		Body:       body,
		Parameters: parameters,
		Results:    results,
	}
}
func NewMethodWithComment(name string, comment string, str NamedTypeValue, body string, parameters, results []NamedTypeValue) Method {
	m := NewMethod(name, str, body, parameters, results)
	m.Comment = prepareComments(comment)
	return m
}
func NewInterface(name string, methods []Method) Interface {
	return Interface{
		Name:    name,
		Comment: "",
		Methods: methods,
	}
}
func NewInterfaceWithComment(name string, comment string, methods []Method) Interface {
	i := NewInterface(name, methods)
	i.Comment = prepareComments(comment)
	return i
}
func prepareComments(comment string) string {
	commentList := strings.Split(comment, "\n")
	comment = ""
	for _, v := range commentList {
		comment += "// " + strings.TrimSpace(v) + "\n"
	}
	return comment
}
func NewStruct(name string, vars []NamedTypeValue) Struct {
	return Struct{
		Name:    name,
		Comment: "",
		Vars:    vars,
	}
}
func NewStructWithComment(name string, comment string, vars []NamedTypeValue) Struct {
	s := NewStruct(name, vars)
	s.Comment = prepareComments(comment)
	return s
}

func NewFile() File {
	return File{
		Interfaces: []Interface{},
		Imports:    []NamedTypeValue{},
		Structs:    []Struct{},
		Vars:       []NamedTypeValue{},
		Constants:  []NamedTypeValue{},
		Methods:    []Method{},
	}
}
