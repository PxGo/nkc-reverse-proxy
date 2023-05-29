package modules

import (
	_ "embed"
	"nkc-reverse-proxy/conf"
	"nkc-reverse-proxy/html"
	"strings"
)

func GetPageByTemplateContent(templateContent TemplateContent) []byte {
	page := strings.Replace(html.PageTemplate, "{{title}}", templateContent.Title, -1)
	page = strings.Replace(page, "{{desc}}", templateContent.Desc, -1)
	page = strings.Replace(page, "{{version}}", conf.CodeVersion, -1)

	return []byte(page)
}
