package server

import (
	"blog"
	"html/template"
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
	postwith('/login', {'pwd':pwd});
}
</script>
</body>
</html>
`

func initTemplate() error {
	tmpl, err := template.ParseFiles(
		filepath.Join(blog.Path("template"), "article.tmpl"),
		filepath.Join(blog.Path("template"), "home.tmpl"),
	)
	if err != nil {
		return err
	}
	templates = tmpl
	return nil
}
