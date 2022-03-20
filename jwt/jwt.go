package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/Ruixues/RWeb"
)

type JwtCore[Jwt any] struct {
	Jwt        Jwt   `json:"jwt"`
	CreateTime int64 `json:"createTime"`
}

type JwtEngine[Jwt any] struct {
	timeLimit   time.Duration
	jwtPool     sync.Pool
	checkers    []JwtChecker[Jwt]
	checkerLock sync.RWMutex
	// 在拦截器中使用，以获得Jwt
	JwtGetter func(context *RWeb.Context) string
}

// jwt为待检测的Jwt, symbol为API的标识符
type JwtChecker[Jwt any] func(jwt *Jwt, tags []string) error

func NewJwtEngine[Jwt any](timeLimit time.Duration) JwtEngine[Jwt] {
	return JwtEngine[Jwt]{
		timeLimit: timeLimit,
		jwtPool: sync.Pool{
			New: func() any {
				return new(JwtCore[Jwt])
			},
		},
		checkers: make([]JwtChecker[Jwt], 0),
	}
}
func (z *JwtEngine[Jwt]) AddJwtChecker(checker JwtChecker[Jwt]) {
	z.checkerLock.Lock()
	defer z.checkerLock.Unlock()
	z.checkers = append(z.checkers, checker)
}

// 不推荐使用
func (z *JwtEngine[Jwt]) RemoveJwtChecker(checker JwtChecker[Jwt]) {
	z.checkerLock.Lock()
	defer z.checkerLock.Unlock()
	checkerPointer := reflect.ValueOf(checker).Pointer()
	for i, v := range z.checkers {
		if reflect.ValueOf(v).Pointer() == checkerPointer {
			z.checkers = append(z.checkers[:i], z.checkers[i+1:]...)
			break
		}
	}
}
func (z *JwtEngine[Jwt]) GenJwt(data *Jwt) (string, error) {
	byteData, err := json.Marshal(JwtCore[Jwt]{
		Jwt:        *data,
		CreateTime: time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(byteData), nil
}
func (z *JwtEngine[Jwt]) Jwt(data string, tags []string) (*JwtCore[Jwt], error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	jwt := z.jwtPool.Get().(*JwtCore[Jwt])
	if err := json.Unmarshal(decoded, &jwt); err != nil {
		return nil, err
	}
	if time.Now().Unix()-jwt.CreateTime > int64(z.timeLimit) {
		return nil, errors.New("retired jwt")
	}
	z.checkerLock.RLock()
	defer z.checkerLock.RUnlock()
	for _, checker := range z.checkers {
		if err := checker(&jwt.Jwt, tags); err != nil {
			return nil, err
		}
	}
	return jwt, nil
}

// 会在context中以关键字jwt储存 类型为用户自定义的Jwt
func (z *JwtEngine[Jwt]) Interceptor(context *RWeb.Context, handler *RWeb.RouterHandler) bool {
	var strJwt string
	if z.JwtGetter == nil { // 那就按照默认,读取body的Jwt
		strJwt = string(context.FormValue("jwt"))
	} else {
		strJwt = z.JwtGetter(context)
	}
	jwt, err := z.Jwt(strJwt, handler.Tags)
	if err != nil {
		return false
	}
	context.StoreValue("jwt", jwt.Jwt)
	return true
}
