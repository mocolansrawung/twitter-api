package tweet

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	tweetQueries = struct {
		selectAllTweets string
		selectTweet     string
		insertTweet     string
	}{
		selectAllTweets: `
			SELECT
				id,
				content,
				retweets,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			FROM tweets
		`,

		selectTweet: `
			SELECT
				id,
				content,
				retweets,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			FROM tweets
		`,

		insertTweet: `
			INSERT INTO tweets (
				id,
				content,
				retweets,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
			) VALUES (
				:id,
				:content,
				:retweets,
				:created_at,
				:created_by,
				:updated_at,
				:updated_by,
				:deleted_at,
				:deleted_by
			)
		`,
	}
)

type TweetRepository interface {
	Create(tweet Tweet) (err error)
	ResolveByID(id uuid.UUID) (tweet Tweet, err error)
	ResolveTweets(page int, limit int, sort string, order string) (tweets []Tweet, err error)
}

type TweetRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideTweetRepositoryMySQL(db *infras.MySQLConn) *TweetRepositoryMySQL {
	s := new(TweetRepositoryMySQL)
	s.DB = db

	return s
}

// Create new Tweet
func (r *TweetRepositoryMySQL) Create(tweet Tweet) (err error) {
	exists, err := r.ExistsByID(tweet.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "tweet", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, tweet); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// Resolve Tweets
func (r *TweetRepositoryMySQL) ResolveTweets(page int, limit int, sort string, order string) (tweets []Tweet, err error) {
	var args []interface{}

	query := tweetQueries.selectAllTweets

	if sort != "" {
		validColumns := map[string]bool{
			"id":         true,
			"content":    false,
			"retweets":   true,
			"created_at": true,
			"created_by": true,
			"updated_at": true,
			"updated_by": true,
			"deleted_at": true,
			"deleted_by": true,
		}
		if !validColumns[sort] {
			return nil, errors.New("Invalid sort parameter")
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}
		if !validOrders[order] {
			return nil, errors.New("Invalid order parameter")
		}

		if order == "" {
			order = "asc"
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	}

	offset := page * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err = r.DB.Read.Select(&tweets, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			err = failure.NotFound("tweets")
			logger.ErrorWithStack(err)
			return
		}

		logger.ErrorWithStack(err)
		return
	}

	return
}

// ResolveByID
func (r *TweetRepositoryMySQL) ResolveByID(id uuid.UUID) (tweet Tweet, err error) {
	err = r.DB.Read.Get(
		&tweet,
		tweetQueries.selectTweet+" WHERE id = ?",
		id.String())

	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("tweet")
		logger.ErrorWithStack(err)
		return
	}

	return
}

// Checking the existence of a Tweet by its ID.
func (r *TweetRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM tweets WHERE id = ?",
		id.String())

	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

// Internal functions
func (r *TweetRepositoryMySQL) txCreate(tx *sqlx.Tx, tweet Tweet) (err error) {
	stmt, err := tx.PrepareNamed(tweetQueries.insertTweet)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tweet)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
