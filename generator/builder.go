// +build generate
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
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
	}
	flag.StringVar(&opt.Input, "i", "", "input file")
	flag.StringVar(&opt.Output, "o", "", "output file")
	flag.BoolVar(&opt.Comment, "comment", false, "comment file")
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

	for _, api := range apis {
		api.WriteHeader()
		api.WriteService()
		api.WriteTrailer()
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
	fmt.Fprintf(out, "package %s\n\n", os.Getenv("GOPACKAGE"))
	if api.Comment {
		fmt.Fprintf(out, "/*\n")
	}
	fmt.Fprintf(out, "import (\n")
	fmt.Fprintf(out, "\t\"encoding/json\"\n")
	fmt.Fprintf(out, "\t\"fmt\"\n")
	fmt.Fprintf(out, "\t\"log\"\n")
	fmt.Fprintf(out, "\t\"net/http\"\n")
	fmt.Fprintf(out, "\t\"net/http/httputil\"\n")
	fmt.Fprintf(out, "\t\"net/url\"\n")
	fmt.Fprintf(out, "\t\"strings\"\n")
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "\t\"github.com/olivere/elastic\"\n")
	fmt.Fprintf(out, "\t\"github.com/olivere/elastic/uritemplates\"\n")
	fmt.Fprintf(out, ")\n\n")

	fmt.Fprintf(out, "var (\n")
	fmt.Fprintf(out, "\t_ = fmt.Print\n")
	fmt.Fprintf(out, "\t_ = log.Print\n")
	fmt.Fprintf(out, "\t_ = strings.Index\n")
	fmt.Fprintf(out, "\t_ = uritemplates.Expand\n")
	fmt.Fprintf(out, "\t_ = url.Parse\n")
	fmt.Fprintf(out, ")\n\n")
}

func (api *Api) WriteTrailer() {
	if api.Comment {
		fmt.Fprintf(out, "*/\n")
	}
}

func (api *Api) WriteService() {
	if api.Documentation != "" {
		fmt.Fprintf(out, "// %s is documented at %s.\n", api.ServiceName(), api.Documentation)
	}
	fmt.Fprintf(out, "type %s struct {\n", api.ServiceName())
	fmt.Fprintf(out, "\tclient\t*elastic.Client\n")
	fmt.Fprintf(out, "\tdebug\tbool\n")
	fmt.Fprintf(out, "\tpretty\tbool\n")

	// Write variables in struct
	fieldsWritten := make(map[string]bool)
	for name, p := range api.URL.Parts {
		fmt.Fprintf(out, "\t%s\t%s\n", p.VariableName(), p.TypeName())
		fieldsWritten[name] = true
	}
	for name, p := range api.URL.Params {
		if found, _ := fieldsWritten[name]; !found {
			fmt.Fprintf(out, "\t%s\t%s\n", p.VariableName(), p.TypeName())
			fieldsWritten[name] = true
		}
	}
	// Write body variable
	if api.Body != nil {
		fmt.Fprintf(out, "\tbodyJson\tinterface{}\n")
		fmt.Fprintf(out, "\tbodyString\tstring\n")
	}
	fmt.Fprintf(out, "}\n\n")

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
	// Write body setter
	if api.Body != nil {
		if api.Body.Description != "" {
			fmt.Fprintf(out, "// BodyJson is documented as: %s.\n", api.Body.Description)
		}
		fmt.Fprintf(out, "func (s *%s) BodyJson(body interface{}) *%s {\n",
			api.ServiceName(),
			api.ServiceName(),
		)
		fmt.Fprintf(out, "\ts.bodyJson = body\n")
		fmt.Fprintf(out, "\treturn s\n")
		fmt.Fprintf(out, "}\n\n")
		if api.Body.Description != "" {
			fmt.Fprintf(out, "// BodyString is documented as: %s.\n", api.Body.Description)
		}
		fmt.Fprintf(out, "func (s *%s) BodyString(body string) *%s {\n",
			api.ServiceName(),
			api.ServiceName(),
		)
		fmt.Fprintf(out, "\ts.bodyString = body\n")
		fmt.Fprintf(out, "\treturn s\n")
		fmt.Fprintf(out, "}\n\n")
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
	fmt.Fprintf(out, "// %s\n", p.Documentation())
	fmt.Fprintf(out, "func (s *%s) %s(%s %s) *%s {\n",
		api.ServiceName(),
		p.SetterName(),
		p.VariableName(),
		p.TypeNameForSetter(),
		api.ServiceName(),
	)
	fmt.Fprintf(out, "\ts.%s = %s\n", p.VariableName(), p.VariableNameForSetter())
	fmt.Fprintf(out, "\treturn s\n")
	fmt.Fprintf(out, "}\n\n")
}

func (p *ApiParam) writeSetter(api *Api) {
	fmt.Fprintf(out, "// %s\n", p.Documentation())
	fmt.Fprintf(out, "func (s *%s) %s(%s %s) *%s {\n",
		api.ServiceName(),
		p.SetterName(),
		p.VariableName(),
		p.TypeNameForSetter(),
		api.ServiceName(),
	)
	fmt.Fprintf(out, "\ts.%s = %s\n", p.VariableName(), p.VariableNameForSetter())
	fmt.Fprintf(out, "\treturn s\n")
	fmt.Fprintf(out, "}\n\n")
}

func (api *Api) writeValidate() {
	emptyBody := true

	fmt.Fprintf(out, "// Validate checks if the operation is valid.\n")
	fmt.Fprintf(out, "func (s *%s) Validate() error {\n", api.ServiceName())
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
		fmt.Fprintf(out, "\treturn nil\n")
		fmt.Fprintf(out, "}\n\n")
		return
	}

	fmt.Fprintf(out, "\tvar invalid []string\n")
	for _, p := range api.URL.Parts {
		if p.Required {
			switch p.Type {
			case "list":
				fmt.Fprintf(out, "\tif len(s.%s) == 0 {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tinvalid = append(invalid, \"%s\")\n", p.SetterName())
				fmt.Fprintf(out, "\t}\n")
			case "string", "text", "enum":
				fmt.Fprintf(out, "\tif s.%s == \"\" {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tinvalid = append(invalid, \"%s\")\n", p.SetterName())
				fmt.Fprintf(out, "\t}\n")
			}
		}
	}
	for _, p := range api.URL.Params {
		if p.Required {
			switch p.Type {
			case "list":
				fmt.Fprintf(out, "\tif len(s.%s) == 0 {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tinvalid = append(invalid, \"%s\")\n", p.SetterName())
				fmt.Fprintf(out, "\t}\n")
			case "string", "text", "enum":
				fmt.Fprintf(out, "\tif s.%s == \"\" {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tinvalid = append(invalid, \"%s\")\n", p.SetterName())
				fmt.Fprintf(out, "\t}\n")
			}
		}
	}
	if api.Body != nil && api.Body.Required {
		fmt.Fprintf(out, "\tif s.bodyString == \"\" && s.bodyJson == nil {\n")
		fmt.Fprintf(out, "\t\tinvalid = append(invalid, \"BodyJson\")\n")
		fmt.Fprintf(out, "\t}\n")
	}
	fmt.Fprintf(out, "\tif len(invalid) > 0 {\n")
	fmt.Fprintf(out, "\t\treturn fmt.Errorf(\"missing required fields: %%v\", invalid)\n")
	fmt.Fprintf(out, "\t}\n")
	fmt.Fprintf(out, "\treturn nil\n")
	fmt.Fprintf(out, "}\n\n")
}

func (api *Api) writeBuildURL() {
	fmt.Fprintf(out, "// buildURL builds the URL for the operation.\n")
	fmt.Fprintf(out, "func (s *%s) buildURL() (string, error) {\n", api.ServiceName())

	/*
		if len(api.URL.Paths) == 0 {
			fmt.Fprintf(out, "\turls := `%s`\n", api.URL.Path)
		} else {
			fmt.Fprintf(out, "\turls := `/`\n", api.URL.Paths[len(api.URL.Paths)-1])
		}
	*/

	if len(api.URL.Parts) > 0 {
		fmt.Fprintf(out, "\t// Build URL\n")
		fmt.Fprintf(out, "\turls, err := uritemplates.Expand(\"%s\", map[string]string{\n", api.URL.Path)
		for _, p := range api.URL.Parts {
			switch p.Type {
			case "list":
				fmt.Fprintf(out, "\t\t\"%s\":\tstrings.Join(s.%s, \",\"),\n", p.Name, p.VariableName())
			default:
				fmt.Fprintf(out, "\t\t\"%s\":\ts.%s,\n", p.Name, p.VariableName())
			}
		}
		fmt.Fprintf(out, "\t})\n")
		fmt.Fprintf(out, "\tif err != nil {\n")
		fmt.Fprintf(out, "\t\treturn \"\", err\n")
		fmt.Fprintf(out, "\t}\n\n")
	} else {
		fmt.Fprintf(out, "\t// Build URL\n")
		fmt.Fprintf(out, "\turls := \"%s\"\n\n", api.URL.Path)
	}

	if len(api.URL.Params) > 0 {
		fmt.Fprintf(out, "\t// Add query string parameters\n")
		fmt.Fprintf(out, "\tparams := url.Values{}\n")
		for _, p := range api.URL.Params {
			switch p.Type {
			case "boolean":
				if !p.Required {
					fmt.Fprintf(out, "\tif s.%s != nil {\n", p.VariableName())
					fmt.Fprintf(out, "\t\tparams.Set(\"%s\", fmt.Sprintf(\"%%v\", *s.%s))\n", p.VariableName(), p.VariableName())
					fmt.Fprintf(out, "\t}\n")
				} else {
					fmt.Fprintf(out, "\tif s.%s {\n", p.VariableName())
					fmt.Fprintf(out, "\t\tparams.Set(\"%s\", \"1\")\n", p.VariableName())
					fmt.Fprintf(out, "\t} else {\n")
					fmt.Fprintf(out, "\t\tparams.Set(\"%s\", \"0\")\n", p.VariableName())
					fmt.Fprintf(out, "\t}\n")
				}
			case "enum":
				fmt.Fprintf(out, "\tif s.%s != \"\" {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", s.%s)\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			case "time", "duration":
				fmt.Fprintf(out, "\tif s.%s != \"\" {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", s.%s)\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			case "number":
				fmt.Fprintf(out, "\tif s.%s != nil {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", fmt.Sprintf(\"%%v\", s.%s))\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			case "string", "text":
				fmt.Fprintf(out, "\tif s.%s != \"\" {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", s.%s)\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			case "list":
				fmt.Fprintf(out, "\tif len(s.%s) > 0 {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", strings.Join(s.%s, \",\"))\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			default:
				fmt.Fprintf(out, "\tif s.%s != nil {\n", p.VariableName())
				fmt.Fprintf(out, "\t\tparams.Set(\"%s\", s.%s)\n", p.VariableName(), p.VariableName())
				fmt.Fprintf(out, "\t}\n")
			}
		}
		fmt.Fprintf(out, "\tif len(params) > 0 {\n")
		fmt.Fprintf(out, "\t\turls += \"?\" + params.Encode()\n")
		fmt.Fprintf(out, "\t}\n\n")
	}
	fmt.Fprintf(out, "\treturn urls, nil\n")
	fmt.Fprintf(out, "}\n\n")
}

func (api *Api) writeDo() {
	fmt.Fprintf(out, "// Do executes the operation.\n")
	fmt.Fprintf(out, "func (s *%s) Do() (*%s, error) {\n",
		api.ServiceName(),
		api.ResponseTypeName(),
	)
	fmt.Fprintf(out, "\t// Check pre-conditions\n")
	fmt.Fprintf(out, "\tif err := s.Validate(); err != nil {\n")
	fmt.Fprintf(out, "\t\treturn nil, err\n")
	fmt.Fprintf(out, "\t}\n\n")

	fmt.Fprintf(out, "\t// Get URL for request\n")
	fmt.Fprintf(out, "\turls, err := s.buildURL()\n")
	fmt.Fprintf(out, "\tif err != nil {\n")
	fmt.Fprintf(out, "\t\treturn nil, err\n")
	fmt.Fprintf(out, "\t}\n\n")

	// Request
	if len(api.Methods) == 0 {
		log.Fatalf("no HTTP methods found")
	}
	fmt.Fprintf(out, "\t// Setup HTTP request\n")
	fmt.Fprintf(out, "\treq, err := s.client.NewRequest(\"%s\", urls)\n", api.Methods[0])
	fmt.Fprintf(out, "\tif err != nil {\n")
	fmt.Fprintf(out, "\t\treturn nil, err\n")
	fmt.Fprintf(out, "\t}\n\n")

	// Set body
	if api.Body != nil {
		fmt.Fprintf(out, "\t// Setup HTTP request body\n")
		fmt.Fprintf(out, "\tif s.bodyJson != nil {\n")
		fmt.Fprintf(out, "\t\treq.SetBodyJson(s.bodyJson)\n")
		fmt.Fprintf(out, "\t} else {\n")
		fmt.Fprintf(out, "\t\treq.SetBodyString(s.bodyString)\n")
		fmt.Fprintf(out, "\t}\n\n")
	}

	// Debug
	fmt.Fprintf(out, "\t// Debug output?\n")
	fmt.Fprintf(out, "\tif s.debug {\n")
	fmt.Fprintf(out, "\t\tout, err := httputil.DumpRequestOut((*http.Request)(req), true)\n")
	fmt.Fprintf(out, "\t\tlog.Printf(\"%%s\\n\", string(out))\n")
	fmt.Fprintf(out, "\t}\n")

	// Get response
	fmt.Fprintf(out, `
	// Get HTTP response
	res, err := s.client.c.Do((*http.Request)(req))
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Debug output?
	if s.debug {
		out, _ := httputil.DumpResponse(res, true)
		log.Printf("%%s\n", string(out))
	}
`)

	// Response
	fmt.Fprintf(out, "\t// Return operation response\n")
	fmt.Fprintf(out, "\tresp := new(%s)\n", api.ResponseTypeName())
	fmt.Fprintf(out, `
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
`)
	fmt.Fprintf(out, "}\n\n")
}

func (api *Api) writeResponseType() {
	fmt.Fprintf(out, "// %s is the response of %s.%s.\n", api.ResponseTypeName(), api.ServiceName(), "Do")
	fmt.Fprintf(out, "type %s struct {\n", api.ResponseTypeName())
	fmt.Fprintf(out, "}\n\n")
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
