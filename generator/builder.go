// +build generate
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var (
	in  io.ReadCloser
	out io.WriteCloser
)

func main() {
	var opt struct {
		Input   string
		Output  string
		Comment bool
		Format  bool
	}
	flag.StringVar(&opt.Input, "i", "", "input file")
	flag.StringVar(&opt.Output, "o", "", "output file")
	flag.BoolVar(&opt.Comment, "comment", false, "comment file")
	flag.BoolVar(&opt.Format, "format", true, "go format output file")
	flag.Parse()

	var err error

	in, err = os.Open(opt.Input)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer in.Close()

	var apis Apis
	if err := json.NewDecoder(in).Decode(&apis); err != nil {
		log.Fatalf("%v", err)
	}
	for name, api := range apis {
		api.Name = name
		api.Comment = opt.Comment

		for name, p := range api.URL.Parts {
			p.Name = name
		}
		for name, p := range api.URL.Params {
			p.Name = name
		}
	}

	if opt.Output == "" {
		f := path.Clean(path.Base(opt.Input))
		if f == "/" || f == "." {
			log.Fatalf("invalid input file: %s", opt.Input)
		}
		parts := strings.SplitN(f, ".", -1)
		output := strings.Join(parts[0:len(parts)-1], ".") + ".gen.go"
		opt.Output = output
	}

	out, err = os.Create(opt.Output)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer out.Close()

	var buf bytes.Buffer

	for _, api := range apis {
		api.p = func(format string, args ...interface{}) {
			_, err := fmt.Fprintf(&buf, format, args...)
			if err != nil {
				panic(err)
			}
		}
		api.pn = func(format string, args ...interface{}) {
			api.p(format+"\n", args...)
		}

		api.WriteHeader()
		api.WriteService()
		api.WriteTrailer()
	}

	if opt.Format {
		// Go fmt output before writing
		clean, err := format.Source(buf.Bytes())
		if err != nil {
			// gofmt failed, so write without formatting
			_, err := out.Write(buf.Bytes())
			if err != nil {
				panic(err)
			}
			return
		}
		_, err = out.Write(clean)
		if err != nil {
			panic(err)
		}
	} else {
		// Do not go fmt
		_, err := out.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

type Apis map[string]*Api

type Api struct {
	Name    string `json:"-"`
	Comment bool   `json:"-"`

	Documentation string   `json:"documentation"`
	Methods       []string `json:"methods"`
	URL           *ApiURL  `json:"url"`
	Body          *ApiBody `json:"body"`

	p  func(format string, args ...interface{}) // raw print
	pn func(format string, args ...interface{}) // print with indent and newline
}

type ApiURL struct {
	Path   string               `json:"path"`
	Paths  []string             `json:"paths"`
	Parts  map[string]*ApiPart  `json:"parts"`
	Params map[string]*ApiParam `json:"params"`
}

type ApiPart struct {
	Name string `json:"-"`

	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type ApiParam struct {
	Name string `json:"-"`

	Type        string      `json:"type"`
	Description string      `json:"description"`
	Options     []string    `json:"options"`
	Default     interface{} `json:"default"`
	Required    bool        `json:"required"`
}

type ApiBody struct {
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Serialize   string `json:"serialize"`
}

func (api *Api) CamelCaseName() string {
	return CamelCase(strings.Replace(api.Name, ".", "_", -1))
}

func (api *Api) ServiceName() string {
	return api.CamelCaseName() + "Service"
}

func (api *Api) ResponseTypeName() string {
	return api.CamelCaseName() + "Response"
}

func (api *Api) HasField(fieldName string) bool {
	for name, _ := range api.URL.Parts {
		if name == fieldName {
			return true
		}
	}
	for name, _ := range api.URL.Params {
		if name == fieldName {
			return true
		}
	}
	return false
}

func (api *Api) WriteHeader() {
	_, pn := api.p, api.pn

	banner := fmt.Sprintf(`// Copyright 2012-%d Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.
`, time.Now().Year())
	pn(banner)

	pn("package %s\n", "elastic") // os.Getenv("GOPACKAGE"))
	if api.Comment {
		pn("/*")
	}
	pn("import (")
	pn("\t\"encoding/json\"")
	pn("\t\"fmt\"")
	pn("\t\"log\"")
	pn("\t\"net/http\"")
	pn("\t\"net/http/httputil\"")
	pn("\t\"net/url\"")
	pn("\t\"strings\"")
	pn("")
	pn("\t\"github.com/olivere/elastic/uritemplates\"")
	pn(")\n")

	pn("var (")
	pn("\t_ = fmt.Print")
	pn("\t_ = http.Request")
	pn("\t_ = httputil.DumpRequest")
	pn("\t_ = log.Print")
	pn("\t_ = strings.Index")
	pn("\t_ = uritemplates.Expand")
	pn("\t_ = url.Parse")
	pn(")\n")
}

func (api *Api) WriteTrailer() {
	_, pn := api.p, api.pn

	if api.Comment {
		pn("*/")
	}
}

func (api *Api) WriteService() {
	_, pn := api.p, api.pn

	if api.Documentation != "" {
		pn("// %s is documented at %s.", api.ServiceName(), api.Documentation)
	}
	pn("type %s struct {", api.ServiceName())
	pn("\tclient\t*Client")
	pn("\tpretty\tbool")

	// Write variables in struct
	fieldsWritten := make(map[string]bool)
	for name, p := range api.URL.Parts {
		pn("\t%s\t%s", p.VariableName(), p.TypeName())
		fieldsWritten[name] = true
	}
	for name, p := range api.URL.Params {
		if found, _ := fieldsWritten[name]; !found {
			pn("\t%s\t%s", p.VariableName(), p.TypeName())
			fieldsWritten[name] = true
		}
	}
	// Write body variable
	if api.Body != nil {
		pn("\tbodyJson\tinterface{}")
		pn("\tbodyString\tstring")
	}
	pn("}\n\n")

	// Write New... method
	pn("// New%s creates a new %s.", api.ServiceName(), api.ServiceName())
	pn("func New%s(client *Client) *%s {", api.ServiceName(), api.ServiceName())
	pn("\treturn &%s{", api.ServiceName())
	pn("\t\tclient: client,")
	fieldsWritten = make(map[string]bool)
	for name, p := range api.URL.Parts {
		if p.IsSlice() {
			pn("\t\t%s:\tmake(%s, 0),", p.VariableName(), p.TypeName())
			fieldsWritten[name] = true
		}
	}
	for name, p := range api.URL.Params {
		if found, _ := fieldsWritten[name]; !found {
			if p.IsSlice() {
				pn("\t\t%s:\tmake(%s, 0),", p.VariableName(), p.TypeName())
				fieldsWritten[name] = true
			}
		}
	}
	pn("\t}")
	pn("}\n")

	// Write setters
	settersWritten := make(map[string]bool)
	for name, p := range api.URL.Parts {
		p.writeSetter(api)
		settersWritten[name] = true
	}
	for name, p := range api.URL.Params {
		if found, _ := settersWritten[name]; !found {
			p.writeSetter(api)
			settersWritten[name] = true
		}
	}
	// Write setter for pretty parameter
	pn("// Pretty indicates that the JSON response be indented and human readable.")
	pn("func (s *%s) Pretty(pretty bool) *%s {", api.ServiceName(), api.ServiceName())
	pn("\ts.pretty = pretty")
	pn("\treturn s")
	pn("}\n")

	// Write body setter
	if api.Body != nil {
		if api.Body.Description != "" {
			pn("// BodyJson is documented as: %s.", api.Body.Description)
		}
		pn("func (s *%s) BodyJson(body interface{}) *%s {",
			api.ServiceName(),
			api.ServiceName(),
		)
		pn("\ts.bodyJson = body")
		pn("\treturn s")
		pn("}\n")
		if api.Body.Description != "" {
			pn("// BodyString is documented as: %s.", api.Body.Description)
		}
		pn("func (s *%s) BodyString(body string) *%s {",
			api.ServiceName(),
			api.ServiceName(),
		)
		pn("\ts.bodyString = body")
		pn("\treturn s")
		pn("}\n")
	}

	// Write buildURL func
	api.writeBuildURL()

	// Write Validate func
	api.writeValidate()

	// Write Do func
	api.writeDo()

	// Write response type
	api.writeResponseType()
}

func (p *ApiPart) writeSetter(api *Api) {
	_, pn := api.p, api.pn

	pn("// %s", p.Documentation())
	pn("func (s *%s) %s(%s %s) *%s {",
		api.ServiceName(),
		p.SetterName(),
		p.VariableName(),
		p.TypeNameForSetter(),
		api.ServiceName(),
	)
	pn("\ts.%s = %s", p.VariableName(), p.VariableNameForSetter())
	pn("\treturn s")
	pn("}\n")
}

func (p *ApiParam) writeSetter(api *Api) {
	_, pn := api.p, api.pn

	pn("// %s", p.Documentation())
	pn("func (s *%s) %s(%s %s) *%s {",
		api.ServiceName(),
		p.SetterName(),
		p.VariableName(),
		p.TypeNameForSetter(),
		api.ServiceName(),
	)
	pn("\ts.%s = %s", p.VariableName(), p.VariableNameForSetter())
	pn("\treturn s")
	pn("}\n")
}

func (api *Api) writeValidate() {
	_, pn := api.p, api.pn

	emptyBody := true

	pn("// Validate checks if the operation is valid.")
	pn("func (s *%s) Validate() error {", api.ServiceName())
	for _, p := range api.URL.Parts {
		if p.Required {
			emptyBody = false
			break
		}
	}
	for _, p := range api.URL.Params {
		if p.Required {
			emptyBody = false
			break
		}
	}
	if emptyBody {
		pn("\treturn nil")
		pn("}\n")
		return
	}

	pn("\tvar invalid []string")
	for _, p := range api.URL.Parts {
		if p.Required {
			switch p.Type {
			case "list":
				pn("\tif len(s.%s) == 0 {", p.VariableName())
				pn("\t\tinvalid = append(invalid, \"%s\")", p.SetterName())
				pn("\t}")
			case "string", "text", "enum":
				pn("\tif s.%s == \"\" {", p.VariableName())
				pn("\t\tinvalid = append(invalid, \"%s\")", p.SetterName())
				pn("\t}")
			}
		}
	}
	for _, p := range api.URL.Params {
		if p.Required {
			switch p.Type {
			case "list":
				pn("\tif len(s.%s) == 0 {", p.VariableName())
				pn("\t\tinvalid = append(invalid, \"%s\")", p.SetterName())
				pn("\t}")
			case "string", "text", "enum":
				pn("\tif s.%s == \"\" {", p.VariableName())
				pn("\t\tinvalid = append(invalid, \"%s\")", p.SetterName())
				pn("\t}")
			}
		}
	}
	if api.Body != nil && api.Body.Required {
		pn("\tif s.bodyString == \"\" && s.bodyJson == nil {")
		pn("\t\tinvalid = append(invalid, \"BodyJson\")")
		pn("\t}")
	}
	pn("\tif len(invalid) > 0 {")
	pn("\t\treturn fmt.Errorf(\"missing required fields: %%v\", invalid)")
	pn("\t}")
	pn("\treturn nil")
	pn("}\n")
}

func (api *Api) writeBuildURL() {
	_, pn := api.p, api.pn

	pn("// buildURL builds the URL for the operation.")
	pn("func (s *%s) buildURL() (string, url.Values, error) {", api.ServiceName())

	/*
		if len(api.URL.Paths) == 0 {
			pn("\tpath := `%s`", api.URL.Path)
		} else {
			pn("\tpath := `/`", api.URL.Paths[len(api.URL.Paths)-1])
		}
	*/

	if len(api.URL.Parts) > 0 {
		pn("\t// Build URL")
		pn("\tpath, err := uritemplates.Expand(\"%s\", map[string]string{", api.URL.Path)
		for _, p := range api.URL.Parts {
			switch p.Type {
			case "list":
				pn("\t\t\"%s\":\tstrings.Join(s.%s, \",\"),", p.Name, p.VariableName())
			default:
				pn("\t\t\"%s\":\ts.%s,", p.Name, p.VariableName())
			}
		}
		pn("\t})")
		pn("\tif err != nil {")
		pn("\t\treturn \"\", url.Values{}, err")
		pn("\t}\n")
	} else {
		pn("\t// Build URL path")
		pn("\tpath := \"%s\"\n", api.URL.Path)
	}

	pn("\t// Add query string parameters")
	pn("\tparams := url.Values{}")
	pn("\tif s.pretty {")
	pn("\t\tparams.Set(\"pretty\", \"1\")")
	pn("\t}")

	if len(api.URL.Params) > 0 {
		for _, p := range api.URL.Params {
			switch p.Type {
			case "boolean":
				if !p.Required {
					pn("\tif s.%s != nil {", p.VariableName())
					pn("\t\tparams.Set(\"%s\", fmt.Sprintf(\"%%v\", *s.%s))", p.Name, p.VariableName())
					pn("\t}")
				} else {
					pn("\tif s.%s {", p.VariableName())
					pn("\t\tparams.Set(\"%s\", \"1\")", p.Name)
					pn("\t} else {")
					pn("\t\tparams.Set(\"%s\", \"0\")", p.Name)
					pn("\t}")
				}
			case "enum":
				pn("\tif s.%s != \"\" {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", s.%s)", p.Name, p.VariableName())
				pn("\t}")
			case "time", "duration":
				pn("\tif s.%s != \"\" {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", s.%s)", p.Name, p.VariableName())
				pn("\t}")
			case "number":
				pn("\tif s.%s != nil {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", fmt.Sprintf(\"%%v\", s.%s))", p.Name, p.VariableName())
				pn("\t}")
			case "string", "text":
				pn("\tif s.%s != \"\" {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", s.%s)", p.Name, p.VariableName())
				pn("\t}")
			case "list":
				pn("\tif len(s.%s) > 0 {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", strings.Join(s.%s, \",\"))", p.Name, p.VariableName())
				pn("\t}")
			default:
				pn("\tif s.%s != nil {", p.VariableName())
				pn("\t\tparams.Set(\"%s\", s.%s)", p.Name, p.VariableName())
				pn("\t}")
			}
		}
	}
	pn("\treturn path, params, nil")
	pn("}\n")
}

func (api *Api) writeDo() {
	_, pn := api.p, api.pn

	pn("// Do executes the operation.")
	pn("func (s *%s) Do() (*%s, error) {",
		api.ServiceName(),
		api.ResponseTypeName(),
	)
	pn("\t// Check pre-conditions")
	pn("\tif err := s.Validate(); err != nil {")
	pn("\t\treturn nil, err")
	pn("\t}\n")

	pn("\t// Get URL for request")
	pn("\tpath, params, err := s.buildURL()")
	pn("\tif err != nil {")
	pn("\t\treturn nil, err")
	pn("\t}\n")

	// Request
	if len(api.Methods) == 0 {
		log.Fatalf("no HTTP methods found")
	}

	// Set body
	if api.Body != nil {
		pn("\t// Setup HTTP request body")
		pn("\tvar body interface{}")
		pn("\tif s.bodyJson != nil {")
		pn("\t\tbody = s.bodyJson")
		pn("\t} else {")
		pn("\t\tbody = s.bodyString")
		pn("\t}\n")
		// Get response
		pn(`
		// Get HTTP response
		res, err := s.client.PerformRequest(method, path, params, body)
		if err != nil {
			return nil, err
		}
`)
	} else {
		// Get response for API without body
		var method = ""
		if len(api.Methods) != 1 {
			pn("\t// TODO: Add method here")
			pn("\tmethod := ...")
		} else {
			method = api.Methods[0]
		}
		pn(`
	// Get HTTP response
	res, err := s.client.PerformRequest("` + strings.ToUpper(method) + `", path, params, nil)
	if err != nil {
		return nil, err
	}
`)
	}

	// Response
	pn("\t// Return operation response")
	pn("\tret := new(%s)", api.ResponseTypeName())
	pn("\tif err := json.Unmarshal(res.Body, ret); err != nil {")
	pn("\t\treturn nil, err")
	pn("\t}")
	pn("\treturn ret, nil")
	pn("}\n")
}

func (api *Api) writeResponseType() {
	_, pn := api.p, api.pn

	pn("// %s is the response of %s.%s.", api.ResponseTypeName(), api.ServiceName(), "Do")
	pn("type %s struct {", api.ResponseTypeName())
	pn("}\n")
}

func (p *ApiPart) VariableName() string {
	s := CamelCase(p.Name)
	if s != "" {
		s = strings.ToLower(s[0:1]) + s[1:]
	}
	if s == "type" {
		s = "typ"
	}
	return s
}

func (p *ApiPart) Documentation() string {
	if p.Description != "" {
		if strings.HasPrefix(p.Description, p.SetterName()) {
			return fmt.Sprintf("%s.", p.Description)
		} else {
			return fmt.Sprintf("%s is documented as: %s.", p.SetterName(), p.Description)
		}
	}
	return fmt.Sprintf("%s is not documented yet.", p.SetterName())
}

func (p *ApiParam) Documentation() string {
	if p.Description != "" {
		if strings.HasPrefix(p.Description, p.SetterName()) {
			return fmt.Sprintf("%s.", p.Description)
		} else {
			return fmt.Sprintf("%s is documented as: %s.", p.SetterName(), p.Description)
		}
	}
	return fmt.Sprintf("%s is not documented yet.", p.SetterName())
}

func (p *ApiParam) VariableName() string {
	s := CamelCase(p.Name)
	if s != "" {
		s = strings.ToLower(s[0:1]) + s[1:]
	}
	if s == "type" {
		s = "typ"
	}
	return s
}

func (p *ApiPart) VariableNameForSetter() string {
	v := p.VariableName()
	switch p.Type {
	case "boolean":
		if !p.Required {
			v = "&" + v
		}
	}
	return v
}

func (p *ApiParam) VariableNameForSetter() string {
	v := p.VariableName()
	switch p.Type {
	case "boolean":
		if !p.Required {
			v = "&" + v
		}
	}
	return v
}

func (p *ApiPart) SetterName() string {
	return CamelCase(p.Name)
}

func (p *ApiParam) SetterName() string {
	return CamelCase(p.Name)
}

func (p *ApiPart) TypeName() string {
	switch p.Type {
	case "boolean":
		return "bool"
	case "enum":
		return "string"
	case "time", "duration":
		return "string"
	case "number":
		return "interface{}"
	case "string", "text":
		return "string"
	case "list":
		return "[]string"
	default:
		return "interface{}"
	}
}

func (p *ApiParam) TypeName() string {
	switch p.Type {
	case "boolean":
		if p.Required {
			return "bool"
		}
		return "*bool"
	case "enum":
		return "string"
	case "time", "duration":
		return "string"
	case "number":
		return "interface{}"
	case "string", "text":
		return "string"
	case "list":
		return "[]string"
	default:
		return "interface{}"
	}
}

func (p *ApiPart) TypeNameForSetter() string {
	switch p.Type {
	case "boolean":
		return "bool"
	case "enum":
		return "string"
	case "time", "duration":
		return "string"
	case "number":
		return "interface{}"
	case "string", "text":
		return "string"
	case "list":
		return "[]string"
	default:
		return "interface{}"
	}
}

func (p *ApiParam) TypeNameForSetter() string {
	switch p.Type {
	case "boolean":
		return "bool"
	case "enum":
		return "string"
	case "time", "duration":
		return "string"
	case "number":
		return "interface{}"
	case "string", "text":
		return "string"
	case "list":
		return "[]string"
	default:
		return "interface{}"
	}
}

func (p *ApiPart) IsSlice() bool {
	return p.Type == "list"
}

func (p *ApiParam) IsSlice() bool {
	return p.Type == "list"
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}
