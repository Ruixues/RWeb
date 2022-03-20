package RWeb

import (
	"mime/multipart"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Context struct {
	RequestUri  string
	Method      int
	RawCtx      *fasthttp.RequestCtx
	keyValueMap map[string]any
}

var json = jsoniter.ConfigFastest

func (z *Context) FormValue(key string) []byte {
	return z.RawCtx.FormValue(key)
}
func (z *Context) Write(data []byte) error {
	_, err := z.RawCtx.Write(data)
	return err
}
func (z *Context) Json(data interface{}) error {
	byte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return z.Write(byte)
}
func (z *Context) String(data string) error {
	_, err := z.RawCtx.WriteString(data)
	return err
}
func (z *Context) Query(key string) string {
	return string(z.RawCtx.QueryArgs().Peek(key))
}
func (z *Context) DefaultQuery(key string, defaultValue string) string {
	if z.RawCtx.QueryArgs().Has(key) {
		return string(z.RawCtx.QueryArgs().Peek(key))
	}
	return defaultValue
}
func (z *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return z.RawCtx.FormFile(key)
}
func (z *Context) MultipartForm() (*multipart.Form, error) {
	return z.RawCtx.MultipartForm()
}
func (z *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 64)
	_, err = f.Read(buf)
	if err != nil {
		return err
	}
	os.Remove(dst)
	ft, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer ft.Close()
	_, err = ft.Write(buf)
	return err
}
func (z *Context) Redirect(uri string, statusCode int) {
	z.RawCtx.Redirect(uri, statusCode)
}
func (z *Context) Cookie(key string) []byte {
	return z.RawCtx.Request.Header.Cookie(key)
}
func (z *Context) Value(key string) any {
	return z.keyValueMap[key]
}
func (z *Context) StoreValue(key string, value any) {
	z.keyValueMap[key] = value
}

var contextPool = &sync.Pool{
	New: func() interface{} {
		return &Context{
			keyValueMap: make(map[string]any),
		}
	},
}

func NewContext() (ret *Context) {
	ret = contextPool.Get().(*Context)
	return ret
}
func RemoveContext(context *Context) {
	contextPool.Put(context)
}
