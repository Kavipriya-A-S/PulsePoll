package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"pulsepoll/backend/config"
	"pulsepoll/backend/models"
)

// CreatePoll creates a new poll
func CreatePoll(c *gin.Context) {
	var poll models.Poll

	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate poll question
	if poll.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Poll question is required",
		})
		return
	}

	// Validate options
	if len(poll.Options) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least 2 options are required",
		})
		return
	}

	// Generate poll ID
	poll.ID = bson.NewObjectID()

	// Store creation time
	poll.CreatedAt = time.Now().Format(time.RFC3339)

	// Generate IDs for options and initialize votes
	for i := range poll.Options {
		poll.Options[i].ID = bson.NewObjectID().Hex()
		poll.Options[i].Votes = 0
	}

	// Insert poll into MongoDB
	_, err := config.DB.Collection("polls").InsertOne(
		context.Background(),
		poll,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create poll",
		})
		return
	}

	c.JSON(http.StatusCreated, poll)
}

// GetPolls returns all polls
func GetPolls(c *gin.Context) {
	cursor, err := config.DB.Collection("polls").Find(
		context.Background(),
		bson.D{},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch polls",
		})
		return
	}

	defer cursor.Close(context.Background())

	var polls []models.Poll

	if err := cursor.All(context.Background(), &polls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to decode polls",
		})
		return
	}

	// Return empty array instead of null when there are no polls
	if polls == nil {
		polls = []models.Poll{}
	}

	c.JSON(http.StatusOK, polls)
}

// VotePoll records one vote for a poll option
func VotePoll(c *gin.Context) {
	pollID := c.Param("id")
	optionID := c.Param("optionId")

	// Convert poll ID to MongoDB ObjectID
	objectID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid poll ID",
		})
		return
	}

	// Find the poll and selected option
	filter := bson.M{
		"_id":        objectID,
		"options.id": optionID,
	}

	// Increase selected option's vote count
	update := bson.M{
		"$inc": bson.M{
			"options.$.votes": 1,
		},
	}

	result, err := config.DB.Collection("polls").UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record vote",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll or option not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Vote recorded successfully",
		"pollId":   pollID,
		"optionId": optionID,
	})
}

// GetPollResults returns vote results for a poll
func GetPollResults(c *gin.Context) {
	pollID := c.Param("id")

	// Convert poll ID to MongoDB ObjectID
	objectID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid poll ID",
		})
		return
	}

	var poll models.Poll

	// Find poll
	err = config.DB.Collection("polls").FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&poll)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found",
		})
		return
	}

	// Calculate total votes
	totalVotes := 0

	for _, option := range poll.Options {
		totalVotes += option.Votes
	}

	// Result structure
	type Result struct {
		Option     string  `json:"option"`
		Votes      int     `json:"votes"`
		Percentage float64 `json:"percentage"`
	}

	results := make([]Result, 0, len(poll.Options))

	// Calculate percentage for each option
	for _, option := range poll.Options {
		percentage := float64(0)

		if totalVotes > 0 {
			percentage = float64(option.Votes) / float64(totalVotes) * 100
		}

		results = append(results, Result{
			Option:     option.Text,
			Votes:      option.Votes,
			Percentage: percentage,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pollId":     poll.ID.Hex(),
		"question":   poll.Question,
		"totalVotes": totalVotes,
		"results":    results,
	})
}