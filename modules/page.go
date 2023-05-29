package modules

import (
	_ "embed"
	"strings"
)

//go:embed ..\html\template.html
var PageTemplate string

//go:embed ..\version.txt
var CodeVersion string

func GetPageByTemplateContent(templateContent TemplateContent) []byte {
	page := strings.Replace(PageTemplate, "{{title}}", templateContent.Title, -1)
	page = strings.Replace(page, "{{desc}}", templateContent.Desc, -1)
	page = strings.Replace(page, "{{version}}", CodeVersion, -1)

	return []byte(page)
}
