package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// Emoji ...
type Emoji struct {
	Emoji          string   `json:"emoji,omitempty"`
	Description    string   `json:"description,omitempty"`
	Category       string   `json:"category,omitempty"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version,omitempty"`
	IosVersion     string   `json:"ios_version,omitempty"`
}

func handleError(err error) {
	if err != nil {
		log.Println("ERR: ", err)
		os.Exit(1)
	}
}

var tmpl = template.Must(template.New("emojis").Funcs(template.FuncMap{
	"title":     strings.Title,
	"getalias":  getAlias,
	"lower":     strings.ToLower,
	"aliasList": aliasList,
}).Parse(`package gomoji

import "strings"

// Emoji constants
const (
{{range .}}
	{{- $this := .}}
	{{- with $alias := getalias .Aliases}}
		{{if ne $alias ""}} Emoji{{$alias}} = "{{$this.Emoji}}" {{end}}
	{{- end}}
{{- end}}
)

// Emoji returns an emoji from multiple possible aliases
func Emoji(name string) string {
	switch strings.ToLower(name) {
{{range .}}
		{{- $alias := getalias .Aliases -}}
		{{- $this := .}}
		{{- with $alias := getalias .Aliases}}
			{{- if ne $alias ""}} 
			case {{aliasList $this.Aliases}}: 
				return Emoji{{getalias $this.Aliases}}
			{{- end}}
		{{- end}}
{{- end}}
	}
	return ""
}
`))

func main() {
	var output bytes.Buffer

	emojiFile, err := os.Open("gen/emojis.json")
	defer emojiFile.Close()
	handleError(err)

	var emojis []Emoji
	handleError(json.NewDecoder(emojiFile).Decode(&emojis))

	handleError(tmpl.Execute(&output, emojis))
	formatted, err := format.Source(output.Bytes())
	handleError(err)

	handleError(ioutil.WriteFile("generated.go", formatted, 0600))
}

func getAlias(aliases []string) string {
	for _, v := range aliases {
		if !regexp.MustCompile("^[a-zA-Z_$0-9]*$").MatchString(v) || strings.Contains(v, "-") {
			continue
		}
		return camelCase(v)
	}
	return ""
}

func camelCase(name string) string {
	return strings.Replace(strings.Title(strings.Replace(name, "_", " ", -1)), " ", "", -1)
}

func aliasList(aliases []string) (retval string) {
	for i, v := range aliases {
		retval += "\"" + strings.ToLower(v) + "\""
		if i != len(aliases)-1 {
			retval += ","
		}
	}
	return
}
