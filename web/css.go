package main
import "web"

func GetCSS(ctx *web.Context) (css string, ok bool) {
	css, ok = ctx.GetSecureCookie("css")
	return 
}

func SetCSS(ctx *web.Context, css string) {
	ctx.SetSecureCookie("css", css, 31556926 )
}
