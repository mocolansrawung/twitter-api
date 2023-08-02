package tweet

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

// Tweet

type Tweet struct {
	ID        uuid.UUID   `db:"id" validate:"required"`
	Content   string      `db:"content" validate:"required"`
	Retweets  int         `db:"retweets" validate:"required,min=0"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

func (t *Tweet) IsDeleted() (deleted bool) {
	return t.DeletedAt.Valid && t.DeletedBy.Valid
}

func (t *Tweet) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ToResponseFormat())
}

// NewFromRequestFormat creates a new Tweet from its request format
func (t Tweet) NewFromRequestFormat(req TweetRequestFormat, userID uuid.UUID) (newTweet Tweet, err error) {
	tweetID, _ := uuid.NewV4()
	newTweet = Tweet{
		ID:        tweetID,
		Content:   req.Content,
		Retweets:  generateRetweets(),
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = newTweet.Validate()

	return
}

func (t *Tweet) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(t)
}

func (t Tweet) ToResponseFormat() TweetResponseFormat {
	resp := TweetResponseFormat{
		ID:        t.ID,
		Content:   t.Content,
		Retweets:  t.Retweets,
		CreatedAt: t.CreatedAt,
		CreatedBy: t.CreatedBy,
		UpdatedAt: t.UpdatedAt,
		UpdatedBy: t.UpdatedBy.Ptr(),
		DeletedAt: t.DeletedAt,
		DeletedBy: t.DeletedBy.Ptr(),
	}

	return resp
}

type TweetRequestFormat struct {
	Content string `json:"content" validate:"required"`
}

type TweetResponseFormat struct {
	ID        uuid.UUID  `json:"id"`
	Content   string     `json:"content"`
	Retweets  int        `json:"retweets"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy uuid.UUID  `json:"createdBy"`
	UpdatedAt null.Time  `json:"updatedAt"`
	UpdatedBy *uuid.UUID `json:"updatedBy"`
	DeletedAt null.Time  `json:"deletedAt,omitempty"`
	DeletedBy *uuid.UUID `json:"deletedBy,omitempty"`
}

// Internal function
func generateRetweets() int {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100)

	return randomNumber
}
