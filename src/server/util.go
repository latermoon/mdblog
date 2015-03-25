package server

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
)

const authFormStirng = `
<!doctype html>
<html>
<head></head>
<body>
<script>
function postwith(to, vals) {
	var form = document.createElement("form");
	form.method = "post";
	form.action = to;
	document.body.appendChild(form);
	for (var k in vals) {
		var input = document.createElement("input");
		input.setAttribute("name", k);
		input.setAttribute("value", vals[k]);
		form.appendChild(input);
	}
	form.submit();
	document.body.removeChild(form);
}

var pwd = prompt('Your password?');
if (!pwd) {
	location.href = '/';
} else {
	postwith('/auth', {'pwd':pwd});
}
</script>
</body>
</html>
`

func currentPassword() string {
	b, err := ioutil.ReadFile(filepath.Join(Workspace, "private", "password.txt"))
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func initTemplate() error {
	tmpl, err := template.ParseFiles(
		filepath.Join(Workspace, "template", "article.tmpl"),
		filepath.Join(Workspace, "template", "home.tmpl"),
	)
	if err != nil {
		return err
	}
	templates = tmpl
	return nil
}
