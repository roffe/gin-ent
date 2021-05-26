package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roffe/gin-ent/ent"
	"github.com/roffe/gin-ent/pkg/viewer"
	"golang.org/x/crypto/bcrypt"
)

type userinfo struct {
	Username string `form:"username" json:"username" binding:"required,alphanumunicode"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=32"`
}

func register(c *gin.Context) {
	var input userinfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user, err := client.User.Create().
		SetUsername(input.Username).
		SetPassword(hash).
		Save(c.Request.Context())
	if err != nil {
		switch {
		case ent.IsConstraintError(err):
			c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
			return
		default:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
	}
	output := fmt.Sprintf("created user %q with id %d", user.Username, user.ID)
	c.JSON(http.StatusOK, gin.H{"message": output})
}

type todoinput struct {
	Text string `json:"text"`
}

func createTodo(c *gin.Context) {
	var input todoinput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx := viewer.FromGinContext(c)
	view := viewer.FromContext(ctx)

	t, err := client.Todo.Create().
		SetOwnerID(view.GetID()).
		SetText(input.Text).
		Save(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save todo"})
		return
	}

	c.JSON(http.StatusOK, t)
}

func getTodos(c *gin.Context) {
	ctx := viewer.FromGinContext(c)
	t, err := client.Todo.Query().All(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load todos"})
		return
	}
	c.JSON(http.StatusOK, t)
}
