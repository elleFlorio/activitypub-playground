package server

import (
	"elleFlorio/activitypub-playground/app"
	"elleFlorio/activitypub-playground/app/model"
	"elleFlorio/activitypub-playground/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	router.GET(app.WebFingerEndpoint, wellKnown)

	router.POST("/users", createUser)
	router.POST("/users/:username/outbox", postToOutbox)
	router.POST("/users/:username/inbox", postToInbox)
	router.GET("/search", searchUser)
	router.GET("/users/:username", getUser)

	port := config.Port
	router.Run("0.0.0.0:" + port)

}

func wellKnown(c *gin.Context) {
	resource := c.Query("resource")

	account := app.ResolveWellKnown(resource)
	c.Header("Content-Type", "application/jrd+json")
	c.IndentedJSON(http.StatusOK, account)
}

func createUser(c *gin.Context) {
	var user struct {
		Username string
		Name     string
	}

	if err := c.BindJSON(&user); err != nil {
		return
	}
	actor := app.CreateUser(user.Username, user.Name)
	c.IndentedJSON(http.StatusCreated, &actor)
}

func searchUser(c *gin.Context) {
	acct := c.Query("acct")
	actor, status := app.SearchUser(acct)
	c.IndentedJSON(status, actor)
}

func getUser(c *gin.Context) {
	username := c.Param("username")
	actor, status := app.GetUser(username)
	c.IndentedJSON(status, actor)
}

func postToOutbox(c *gin.Context) {
	username := c.Param("username")
	var activity model.Activity

	if err := c.BindJSON(&activity); err != nil {
		return
	}
	id, status := app.AddToOutbox(username, activity)
	c.Header("Location", id)
	c.IndentedJSON(status, nil)
}

func postToInbox(c *gin.Context) {
	username := c.Param("username")
	var activity model.Activity

	if err := c.BindJSON(&activity); err != nil {
		return
	}
	status := app.AddToInbox(username, activity)
	c.IndentedJSON(status, nil)
}
