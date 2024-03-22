package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgError    error
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			pgInstance = nil
			pgError = fmt.Errorf("unable to create connection pool: %w", err)
		} else {
			pgInstance = &postgres{db}
			pgError = nil
		}
	})
	return pgInstance, pgError
}

const databaseUrl = "postgres://postgres:lockpicks@localhost:5432/postgres"

func getDefaultPG() (*postgres, error) {
	return NewPG(context.Background(), databaseUrl)
}
func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

// vote represents data about a record vote.
type vote struct {
	ID         int  `json:"id"`
	CampaignId int  `json:"campaign_id"`
	UserId     int  `json:"user_id"`
	UnionId    int  `json:"union_id"`
	Approve    bool `json:"approve"`
}

var (
	voteTableOnce sync.Once
)

func beforeVoteAccess() {
	voteTableOnce.Do(func() {
		pg, _ := getDefaultPG()
		_, err := pg.db.Query(context.Background(), "CREATE TABLE votes (CampaignId: int, UserId: int, UnionId: int, Approve: bit)")
		if err != nil {
			fmt.Println("vote table already created")
		}
	})
}

// postVotes adds an vote from JSON received in the request body.
func postVotes(c *gin.Context) {
	// This should probably be middleware but that for later
	beforeVoteAccess()
	var newvote vote

	// Call BindJSON to bind the received JSON to
	// newvote.
	if err := c.BindJSON(&newvote); err != nil {
		return
	}
	// library to consider: jet for SQL?
	// pgx is appropriate for the demo in question.
	// Add the new vote to the slice.
	// votes = append(votes, newvote)
	c.IndentedJSON(http.StatusCreated, newvote)
}

// getVotes responds with the list of all votes as JSON.
func getVotes(c *gin.Context) {
	// This should probably be middleware but let's not for now
	pg, _ := getDefaultPG()
	rows, _ := pg.db.Query(context.Background(), "select * from votes")
	votes := []vote{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal("error while iterating dataset")
		}

		// convert DB types to Go types

		votes = append(votes, vote{
			ID:         values[0].(int),
			CampaignId: values[1].(int),
			UserId:     values[2].(int),
			UnionId:    values[3].(int),
			Approve:    values[4].(bool),
		})

	}
	c.IndentedJSON(http.StatusOK, votes)
}

func main() {
	router := gin.Default()
	router.GET("/votes", getVotes)
	router.POST("/votes", postVotes)
	router.Run("localhost:8080")
	// router.RunTLS(":8080", "server.pem", "server.key")
}
