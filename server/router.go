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
	router.POST("/users/:username/post", postToObjectOutbox)
	router.POST("/users/:username/outbox", postToActivityOutbox)
	router.POST("/users/:username/inbox", postToInbox)
	router.GET("/search", searchUser)
	router.GET("/users/:username/object/:id", getObject)
	router.GET("/users/:username", getUser)
	router.GET("/users/:username/followers/requests", getFollowRequests)
	router.GET("/users/:username/followers", getFollowers)
	router.GET("/users/:username/following", getFollowing)
	router.GET("/users/:username/posts", getPostsByUser)
	router.GET("/users/:username/timeline", getUserTimeline)

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
		c.IndentedJSON(http.StatusBadRequest, model.Actor{})
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

func getFollowRequests(c *gin.Context) {
	username := c.Param("username")
	followRequests, status := app.GetFollowRequests(username)
	c.IndentedJSON(status, followRequests)
}

func getFollowers(c *gin.Context) {
	username := c.Param("username")
	followers, status := app.GetFollowers(username)
	c.IndentedJSON(status, followers)
}

func getFollowing(c *gin.Context) {
	username := c.Param("username")
	following, status := app.GetFollowing(username)
	c.IndentedJSON(status, following)
}

func getPostsByUser(c *gin.Context) {
	username := c.Param("username")
	posts, status := app.GetPosts(username)
	c.IndentedJSON(status, posts)
}

func getUserTimeline(c *gin.Context) {
	username := c.Param("username")
	timeline, status := app.GetTimeline(username)
	c.IndentedJSON(status, timeline)
}

func getObject(c *gin.Context) {
	username := c.Param("username")
	id := c.Param("id")
	objects, status := app.GetObjectById(username, id)
	c.IndentedJSON(status, objects)
}

func postToActivityOutbox(c *gin.Context) {
	username := c.Param("username")
	var activity model.Activity

	if err := c.BindJSON(&activity); err != nil {
		return
	}
	id, status := app.AddActivityToOutbox(username, activity)
	c.Header("Location", id)
	c.IndentedJSON(status, nil)
}

func postToObjectOutbox(c *gin.Context) {
	username := c.Param("username")
	var object model.Object

	if err := c.BindJSON(&object); err != nil {
		return
	}
	id, status := app.AddObjectToOutbox(username, object)
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
