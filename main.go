package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// votes slice to seed record vote data.
var votes = []vote{
	{CampaignId: "1", UserId: "1", UnionId: "1", Approve: false},
}

// vote represents data about a record vote.
type vote struct {
	ID         string `json:"id"`
	CampaignId string `json:"campaign_id"`
	UserId     string `json:"user_id"`
	UnionId    string `json:"union_id"`
	Approve    bool   `json:"approve"`
}

// postvotes adds an vote from JSON received in the request body.
func postvotes(c *gin.Context) {
	var newvote vote

	// Call BindJSON to bind the received JSON to
	// newvote.
	if err := c.BindJSON(&newvote); err != nil {
		return
	}
	// library to consider: jet for SQL?
	// Add the new vote to the slice.
	votes = append(votes, newvote)
	c.IndentedJSON(http.StatusCreated, newvote)
}

// getvotes responds with the list of all votes as JSON.
func getvotes(c *gin.Context) {
	return pgx.
		c.IndentedJSON(http.StatusOK, votes)
}

func main() {
	router := gin.Default()
	router.GET("/votes", getvotes)
	router.POST("/votes", postvotes)
	router.Run("localhost:8080")
	// router.RunTLS(":8080", "server.pem", "server.key")
}
