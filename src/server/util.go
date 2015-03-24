package server

import (
	"io/ioutil"
	"path/filepath"
)

const authFormStirng = `
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
postwith('/auth', {'pwd':pwd});
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
