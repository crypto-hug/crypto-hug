package print

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"

	"github.com/v-braun/go-must"
	"gopkg.in/workanator/go-ataman.v1"
)

var re = regexp.MustCompile(`\{\{[^}]*}}`)

type Fields map[string]interface{}

func LineTpl(message string, data interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	message = re.ReplaceAllStringFunc(message, func(val string) string {
		return "{-+white}" + val + "{-+light+cyan}"
	})
	message = "{-+light+cyan}" + message + "{-}"

	tmpl, err := template.New("").Parse(message)
	must.NoError(err, "%s", err)

	var newMsg bytes.Buffer
	err = tmpl.Execute(&newMsg, data)

	rndr := ataman.NewRenderer(ataman.CurlyStyle())
	result := rndr.MustRender(newMsg.String())

	fmt.Println(result)
}
