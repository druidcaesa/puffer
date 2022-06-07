package cookie

import (
	"net/http"
)

type CookitFunc interface {
	Cookie(key string) (*http.Cookie, error)
}

type Cookie struct {
	r    *http.Request
	resp http.ResponseWriter
}

func (c *Cookie) SetReq(r *http.Request) {
	c.r = r
}

func (c *Cookie) SetResp(r http.ResponseWriter) {
	c.resp = r
}

// Cookie get cookies
func (c *Cookie) Cookie(key string) (*http.Cookie, error) {
	return c.r.Cookie(key)
}

/**
 * @author fanyanan
 * @description set cookie function
 * @date 14:11 2022/6/7
 * @param key cookie key
 * @param value cookie value
 * @param domain domain name
 * @param maxAge Maximum aging unit second
 * @param secure Can it be accessed via https
 * @param httpOnly Whether to allow js to get
 **/
func (c *Cookie) SetCookie(key, value, path, domain string, maxAge int, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.resp, cookie)
}
