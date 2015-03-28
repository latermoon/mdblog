package controller

import (
	"blog"
	"crypto/md5"
	"fmt"
	"github.com/martini-contrib/sessions"
	"io"
	"log"
	"net/http"
	"strings"
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

func AuthHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	if strings.HasPrefix(r.URL.Path, "/private/") {
		salt := blog.Config().Password
		auth := session.Get(blog.Config().AuthKey)
		if auth != salt {
			io.WriteString(w, authFormStirng)
		}
	}
}

func LoginAction(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	pwd := r.PostFormValue("pwd")
	salt := fmt.Sprintf("%x", md5.Sum([]byte(blog.Config().AuthKey+pwd)))
	if salt == blog.Config().Password {
		session.Set(blog.Config().AuthKey, salt)
	} else {
		log.Println("login fail", pwd)
	}
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func LogoutAction(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	session.Delete(blog.Config().AuthKey)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
