package https

import (
	"github.com/valyala/fasthttp"
)

func Response(ctx *fasthttp.RequestCtx, send_s3 []byte, contentType string) {
	ctx.SetContentType(contentType)
	ctx.SetBody([]byte(send_s3))
	ctx.Response.Header.Set("X-My-Header", "my-header-value")
	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

func DoRequest(url string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)
	return resp
}
