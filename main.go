package main

import (
	"context"
	"errors"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/roffe/gin-ent/ent"
	_ "github.com/roffe/gin-ent/ent/runtime"
	"github.com/roffe/gin-ent/ent/user"
	"github.com/roffe/gin-ent/pkg/viewer"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

var client *ent.Client

func setupDB() {
	var err error
	client, err = ent.Open("sqlite3", "file:ent.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

}

func main() {
	setupDB()
	defer client.Close()

	authMiddleware, err := newAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", authMiddleware.LoginHandler)
	r.Use(authMiddleware.MiddlewareFunc())
	{
		r.POST("/todo", createTodo)
		r.GET("/todos", getTodos)
		r.DELETE("/todo/:id", deleteTodo)
		r.PATCH("/todo/:id", updateTodo)
	}
	r.Run()
}

var identityKey = "id"

const ctxKey = "viewer"

func newAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*ent.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"username":  v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &ent.User{
				ID:       int(claims[identityKey].(float64)),
				Username: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var input userinfo
			if err := c.ShouldBind(&input); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			user, err := client.User.Query().
				Where(user.UsernameEqualFold(input.Username)).
				First(c.Request.Context())
			if err != nil {
				switch {
				case ent.IsNotFound(err):
					return nil, jwt.ErrFailedAuthentication
				default:
					return nil, errors.New("failed to load user")
				}
			}
			if err := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if user, ok := data.(*ent.User); ok {
				c.Set(ctxKey, viewer.UserViewer{Role: viewer.View, ID: user.ID})
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
