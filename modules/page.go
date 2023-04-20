package modules

import (
	"strings"
)

const PageTemplate = "<!doctype html>" +
	"<html lang=\"en\">" +
	"<head>" +
	"<meta charset=\"UTF-8\">" +
	"<meta name=\"viewport\" content=\"width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0\">" +
	"<meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">" +
	"<title>{{title}}</title>" +
	"<style>" +
	"html, body{" +
	"	padding: 0;" +
	"	margin: 0;" +
	"	font-size: 12px;" +
	"}" +
	".container{" +
	"	box-sizing: border-box;" +
	"	padding: 12px;" +
	"	height: 100vh;" +
	"	width: 100vw;" +
	"	display: flex;" +
	"	justify-content: center;" +
	"	align-items: center;" +
	"	flex-direction: column;" +
	"}" +
	".container p{" +
	"	font-size: 1.25rem;" +
	"   margin-bottom: 3rem;" +
	"}" +
	".container a{" +
	"	color: #000;" +
	"}" +
	"</style>" +
	"</head>" +
	"<body>" +
	"<div class=\"container\">" +
	"	<h1>{{title}}</h1>" +
	"	<p>{{desc}}</p>" +
	"   <a href=\"https://github.com/kccd/nkc-reverse-proxy\" target=\"_blank\">NKC-REVERSE-PROXY</a>" +
	"</div>" +
	"</body>" +
	"</html>"

func GetPageByTemplateContent(templateContent TemplateContent) []byte {
	page := strings.Replace(PageTemplate, "{{title}}", templateContent.Title, -1)
	page = strings.Replace(page, "{{desc}}", templateContent.Desc, -1)

	return []byte(page)
}
