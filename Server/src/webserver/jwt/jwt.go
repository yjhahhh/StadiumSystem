package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	UserRole  = 1
	AdminRole = 2
)

//用户信息类，作为生成token的参数
type UserClaims struct {
	Number string
	Name   string
	Role   uint
	Level  uint
	//jwt-go提供的标准claim
	jwt.StandardClaims
}

var (
	//自定义的token秘钥
	secret = []byte("yjh666")
	noVerify = []string{"/user/login", "/user/register", "/admin/login", "/admin/register"}
	//token有效时间（纳秒）
	effectTime = 2 * time.Hour
)
	
// 生成token
func GenerateToken(claims *UserClaims) string {
	//设置token有效期
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return sign
}

//验证token
func JwtVerify(c *gin.Context) {
	uri := c.Request.RequestURI
	for _, r := range noVerify {
		if uri == r {
			return
		}
	}
	token := c.GetHeader("Token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status" : http.StatusUnauthorized, "msg" : "401 Unauthorized", "data" : nil})
		return
	}
	//验证token，并存储在请求中
	claims, err := parseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status" : http.StatusUnauthorized, "msg" : "401 Unauthorized", "data" : nil})
		return
	}
	c.Set("number", claims.Number)
	c.Set("name", claims.Name)
	c.Set("role", claims.Role)
	c.Set("level", claims.Level)
}

// 解析Token
func parseToken(tokenString string) (*UserClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("token is valid")
	}
	return claims, nil
}

// 更新token
func Refresh(tokenString string) string {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("token is valid")
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
	return GenerateToken(claims)
}