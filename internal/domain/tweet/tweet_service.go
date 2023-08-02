package tweet

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type TweetService interface {
	Create(requestFormat TweetRequestFormat, userID uuid.UUID) (tweet Tweet, err error)
	ResolveByID(id uuid.UUID) (tweet Tweet, err error)
	ResolveTweets(page int, limit int, sort string, order string) (tweets []Tweet, err error)
	SoftDelete(id uuid.UUID, userID uuid.UUID) (tweet Tweet, err error)
	Update(id uuid.UUID, requestFormat TweetRequestFormat, userID uuid.UUID) (tweet Tweet, err error)
}

type TweetServiceImpl struct {
	TweetRepository TweetRepository
	Config          *configs.Config
}

func ProvideTweetServiceImpl(tweetRepository TweetRepository, config *configs.Config) *TweetServiceImpl {
	s := new(TweetServiceImpl)
	s.TweetRepository = tweetRepository
	s.Config = config

	return s
}

// Create creates a new Tweet
func (s *TweetServiceImpl) Create(requestFormat TweetRequestFormat, userID uuid.UUID) (tweet Tweet, err error) {
	tweet, err = tweet.NewFromRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}

	if err != nil {
		return tweet, failure.BadRequest(err)
	}

	err = s.TweetRepository.Create(tweet)

	if err != nil {
		return
	}

	return
}

// Resolve All Tweets
func (s *TweetServiceImpl) ResolveTweets(page int, limit int, sort string, order string) (tweets []Tweet, err error) {
	tweets, err = s.TweetRepository.ResolveTweets(page, limit, sort, order)
	if err != nil {
		return tweets, failure.BadRequest(err)
	}

	return
}

// ResolveByID
func (s *TweetServiceImpl) ResolveByID(id uuid.UUID) (tweet Tweet, err error) {
	tweet, err = s.TweetRepository.ResolveByID(id)

	if tweet.IsDeleted() {
		return tweet, failure.NotFound("tweet")
	}

	return
}

// SoftDelete
func (s *TweetServiceImpl) SoftDelete(id uuid.UUID, userID uuid.UUID) (tweet Tweet, err error) {
	tweet, err = s.TweetRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = tweet.SoftDelete(userID)
	if err != nil {
		return
	}

	err = s.TweetRepository.Update(tweet)
	return
}

// Update updates a Tweet
func (s *TweetServiceImpl) Update(id uuid.UUID, requestFormat TweetRequestFormat, userID uuid.UUID) (tweet Tweet, err error) {
	tweet, err = s.TweetRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = tweet.Update(requestFormat, userID)
	if err != nil {
		return
	}

	err = s.TweetRepository.Update(tweet)
	return
}
