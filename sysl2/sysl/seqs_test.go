package main

import (
	"bytes"
	"flag"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/seqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type labeler struct {
}

func (l *labeler) LabelApp(appName, controls string, attrs map[string]*sysl.Attribute) string {
	return appName
}

func (l *labeler) LabelEndpoint(p *seqs.EndpointLabelerParam) string {
	return p.EndpointName
}

func TestGenerateSequenceDiag(t *testing.T) {
	m, _ := Parse("demo/simple/sysl-sd.sysl", "../../")
	l := &labeler{}
	p := &sequenceDiagParam{}
	p.endpoints = []string{"WebFrontend <- RequestProfile"}
	p.AppLabeler = l
	p.EndpointLabeler = l
	p.title = "Profile"
	r, err := generateSequenceDiag(m, p)

	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "WebFrontend" as _0
control "Api" as _1
database "Database" as _2
skinparam maxMessageSize 250
title Profile
== WebFrontend <- RequestProfile ==
[->_0 : RequestProfile
activate _0
 _0->_1 : GET /users/{user_id}/profile
 activate _1
  _1->_2 : QueryUser
  activate _2
  _1<--_2 : User
  deactivate _2
 _0<--_1 : UserProfile
 deactivate _1
[<--_0 : Profile Page
deactivate _0
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

func TestGenerateSequenceDiagramsToFormatNameAttributes(t *testing.T) {
	m, _ := Parse("sequence_diagram_name_format.sysl", "./tests/")
	al := seqs.MakeFormatParser(`%(@status?<color red>%(appname)</color>|%(appname))`)
	el := seqs.MakeFormatParser(`%(@status? <color green>%(epname)</color>|%(epname))`)
	p := &sequenceDiagParam{}
	p.endpoints = []string{"User <- Check Balance"}
	p.AppLabeler = al
	p.EndpointLabeler = el
	r, err := generateSequenceDiag(m, p)
	p.title = "Diagram"
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
actor "User" as _0
boundary "MobileApp" as _1
control "<color red>Server</color>" as _2
database "DB" as _3
skinparam maxMessageSize 250
== User <- Check Balance ==
 _0->_1 : Login
 activate _1
  _1->_2 : Login
  activate _2
  _2 -> _2 : do input validation
   _2->_3 :  <color green>Save</color>
  _1<--_2 : success or failure
  deactivate _2
 deactivate _1
 _0->_1 : Check Balance
 activate _1
  _1->_2 : Read User Balance
  activate _2
   _2->_3 :  <color green>Load</color>
  _1<--_2 : balance
  deactivate _2
 deactivate _1
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

func TestGenerateSequenceDiagramsToFormatComplexAttributes(t *testing.T) {
	m, _ := Parse("sequence_diagram_name_format.sysl", "./tests/")
	al := seqs.MakeFormatParser(`%(@status?<color red>%(appname)</color>|%(appname))`)
	el := seqs.MakeFormatParser(`%(@status? <color green>%(epname)</color>|%(epname))`)
	p := &sequenceDiagParam{}
	p.endpoints = []string{"User <- Check Balance"}
	p.AppLabeler = al
	p.EndpointLabeler = el
	r, err := generateSequenceDiag(m, p)
	p.title = "Diagram"
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
actor "User" as _0
boundary "MobileApp" as _1
control "<color red>Server</color>" as _2
database "DB" as _3
skinparam maxMessageSize 250
== User <- Check Balance ==
 _0->_1 : Login
 activate _1
  _1->_2 : Login
  activate _2
  _2 -> _2 : do input validation
   _2->_3 :  <color green>Save</color>
  _1<--_2 : success or failure
  deactivate _2
 deactivate _1
 _0->_1 : Check Balance
 activate _1
  _1->_2 : Read User Balance
  activate _2
   _2->_3 :  <color green>Load</color>
  _1<--_2 : balance
  deactivate _2
 deactivate _1
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

type loadAppArgs struct {
	root   string
	models string
}

func TestLoadAppReturnError(t *testing.T) {
	test := loadAppArgs{
		"../../demo/simple/", "",
	}
	mod := loadApp(test.root, test.models)
	assert.Nil(t, mod)
}

func TestLoadApp(t *testing.T) {
	test := loadAppArgs{
		"./tests/", "sequence_diagram_test.sysl",
	}
	mod := loadApp(test.root, test.models)
	assert.NotNil(t, mod)
	apps := mod.GetApps()
	app := apps["Database"]

	assert.Equal(t, []string{"Database"}, app.GetName().GetPart())

	appPatternsAttr := app.GetAttrs()["patterns"].GetA().GetElt()
	patterns := make([]string, 0, len(appPatternsAttr))
	for _, val := range appPatternsAttr {
		patterns = append(patterns, val.GetS())
	}
	assert.Equal(t, []string{"db"}, patterns)

	queryUserParams := app.GetEndpoints()["QueryUser"].GetParam()
	params := make([]string, 0, len(queryUserParams))
	for _, val := range queryUserParams {
		params = append(params, val.GetName())
	}
	assert.Equal(t, []string{"user_id"}, params)
}

type sdArgs struct {
	rootModel      string
	endpointFormat string
	appFormat      string
	title          string
	output         string
	endpoints      []string
	apps           []string
	modules        string
	blackboxes     [][]string
}

func TestDoConstructSequenceDiagramsNoSyslSdFiltersWithoutEndpoints(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_test.sysl",
	}
	expected := make(map[string]string)

	// When
	result := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes)

	// Then
	assert.Equal(t, expected, result, "unexpected outSet!")
}

func TestDoConstructSequenceDiagramsNoSyslSdFilters(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_test.sysl",
		endpoints: []string{"QueryUser"},
		output:    "_.png",
	}
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "" as _0
skinparam maxMessageSize 250
==  <-  ==
@enduml
`
	expected := map[string]string{
		"_.png": expectContent,
	}

	// When
	result := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes)

	// Then
	assert.Equal(t, expected, result, "unexpected outSet!")
}

func TestDoConstructSequenceDiagrams(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_project.sysl",
		output:    "%(epname).png",
		apps:      []string{"Project"},
	}
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "" as _0
control "" as _1
database "" as _2
skinparam maxMessageSize 250
title Profile
== WebFrontend <- RequestProfile ==
[->_0 : RequestProfile
activate _0
 _0->_1 :` + " " + `
 activate _1
  _1->_2 :` + " " + `
  activate _2
  _1<--_2 : User
  deactivate _2
 _0<--_1 : UserProfile
 deactivate _1
[<--_0 : Profile Page
deactivate _0
@enduml
`
	expected := map[string]string{
		"_.png": expectContent,
	}

	// When
	result := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes)

	// Then
	assert.Equal(t, expected, result, "unexpected content!")
}

func TestDoConstructSequenceDiagramsToFormatComplexName(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_complex_format.sysl",
		output:    "%(epname).png",
		apps:      []string{"Project"},
	}
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "//te//\n<color grey>Ex e</color>\n**User**" as _0
control "**MobileApp**" as _1
skinparam maxMessageSize 250
title Diagram
== User <- Check Balance ==
[->_0 : Check Balance
activate _0
 _0->_1 : //«hello»//** <color red>pat?</color>**aa Login
 deactivate _0
@enduml
`
	expected := map[string]string{
		"Seq.png": expectContent,
	}

	// When
	result := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes)

	// Then
	assert.Equal(t, expected, result, "unexpected content!")
}

func TestDoGenerateSequenceDiagrams(t *testing.T) {
	type args struct {
		flags *flag.FlagSet
		args  []string
	}
	argsData := []string{"sd"}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
	}{
		{
			"Case-Do generate sequence diagrams",
			args{
				flag.NewFlagSet(argsData[0], flag.PanicOnError),
				argsData,
			},
			"",
			"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			require.NoError(t, DoGenerateSequenceDiagrams(stdout, stderr, tt.args.args))
			assert.Equal(t, tt.wantStdout, stdout.String())
			assert.Equal(t, tt.wantStderr, stderr.String())
		})
	}
}
