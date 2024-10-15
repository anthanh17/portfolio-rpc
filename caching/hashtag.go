package cache

import (
	"context"
	"fmt"
	"portfolio-profile-rpc/util"

	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

type SessionType struct {
	SessionID string `json:"sessionId"`
	Username  string `json:"username"`
}

type HashtagCache struct {
	cachier Cachier
	logger  *zap.Logger
}

func NewHashtagCache(cachier Cachier, logger *zap.Logger) *HashtagCache {
	return &HashtagCache{
		cachier: cachier,
		logger:  logger,
	}
}

// ---------------- 1. Hashtag Leaderboard -------------------
// Hashtag - Leaderboard using "sorted set" data structure

// getHashtagLeaderboardCacheKey returns the cache key for the hashtag leaderboard.
// The key is formatted as "hashtag:leaderboard".
func (h HashtagCache) getHashtagLeaderboardCacheKey() string {
	return "hashtag:leaderboard"
}

type HashtagLeaderboard struct {
	Score  float64
	Member string
}

func convertToRedisZ(leaderboard []HashtagLeaderboard) []redis.Z {
	redisZList := make([]redis.Z, len(leaderboard))
	for i, item := range leaderboard {
		redisZList[i] = redis.Z{
			Score:  item.Score,
			Member: item.Member,
		}
	}
	return redisZList
}

func (h HashtagCache) AddHashtagLeaderboardCacheElements(ctx context.Context, leaderboard []HashtagLeaderboard) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	redisZList := convertToRedisZ(leaderboard)

	key := h.getHashtagLeaderboardCacheKey()
	err := cachier.redisClient.ZAdd(ctx, key, redisZList...).Err()

	if err != nil {
		h.logger.With(zap.Error(err)).Error("error AddElementHashtagLeaderboardCache")
		return false, fmt.Errorf("error adding elements to hashtag leaderboard cache: %w", err)
	}

	return true, nil
}

func (h HashtagCache) DeleteHashtagLeaderboardCacheElements(ctx context.Context, hashtags []string) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	key := h.getHashtagLeaderboardCacheKey()
	err := cachier.redisClient.ZRem(ctx, key, util.ConvertSliceStringToSliceInterface(hashtags)...).Err()

	if err != nil {
		h.logger.With(zap.Error(err)).Error("error DeleteHashtagLeaderboardCacheElements")
		return false, fmt.Errorf("error deleting elements from hashtag leaderboard cache: %w", err)
	}

	return true, nil
}

func (h HashtagCache) GetAllHashtagLeaderboardCache(ctx context.Context) ([]string, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return nil, fmt.Errorf("cachier is not redis")
	}

	key := h.getHashtagLeaderboardCacheKey()
	hashtags, err := cachier.redisClient.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		h.logger.With(zap.Error(err)).Error("error DeleteHashtagLeaderboardCacheElements")
		return nil, fmt.Errorf("error deleting elements from hashtag leaderboard cache: %w", err)
	}

	// // If miss cache
	// if len(hashtags) == 0 {
	// 	h.logger.Info("miss cache")
	// 	return nil, ErrCacheMiss
	// }

	return hashtags, nil
}

// ---------------- 2. Hashtag Profile -----------------------
// Hashtag - Profile using "set" data structure
// Key: "hashtag:portfolio_profile:{hashtag}"
// Value - "set": proifileIds

const (
	HashtagProfilesCacheKey string = "hashtag_profiles"
	ProfileHashtagsCacheKey string = "profile_hashtags"
)

func (h HashtagCache) getHashtagProfilesCacheKey(hashtag string) string {
	return fmt.Sprintf("hashtag:%s:%s", HashtagProfilesCacheKey, hashtag)
}

func (h HashtagCache) getProfileHashtagsCacheKey(profileId string) string {
	return fmt.Sprintf("hashtag:%s:%s", ProfileHashtagsCacheKey, profileId)
}

type HashtagProfileIds struct {
	Hashtag    string
	ProfileIds []string
}

func (h HashtagCache) AddHashtagProfileIdsCacheElements(ctx context.Context, hashtagProfileIds HashtagProfileIds) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	key := h.getHashtagProfilesCacheKey(hashtagProfileIds.Hashtag)
	err := cachier.redisClient.SAdd(ctx, key, util.ConvertSliceStringToSliceInterface(hashtagProfileIds.ProfileIds)...).Err()

	if err != nil {
		h.logger.With(zap.Error(err)).Error("error AddHashtagProfileIdsCacheElements")
		return false, fmt.Errorf("error adding elements to hashtag profile cache: %w", err)
	}

	return true, nil
}

type ProfileIdHashtags struct {
	ProfileId string
	Hashtags  []string
}

func (h HashtagCache) AddProfileIdHashtagsCacheElements(ctx context.Context, profileIdHashtags ProfileIdHashtags) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	key := h.getProfileHashtagsCacheKey(profileIdHashtags.ProfileId)
	err := cachier.redisClient.SAdd(ctx, key, util.ConvertSliceStringToSliceInterface(profileIdHashtags.Hashtags)...).Err()

	if err != nil {
		h.logger.With(zap.Error(err)).Error("error AddProfileIdHashtagsCacheElements")
		return false, fmt.Errorf("error adding elements to hashtag profile cache: %w", err)
	}

	return true, nil
}

func (h HashtagCache) GetHashtagProfileIdsCacheElements(ctx context.Context, hashtag string) ([]string, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return nil, fmt.Errorf("cachier is not redis")
	}

	// get key cache
	hashtagProfileIdsKey := h.getHashtagProfilesCacheKey(hashtag)

	// Check key exists
	exists, err := cachier.redisClient.Exists(ctx, hashtagProfileIdsKey).Result()
	if err != nil {
		h.logger.Info("error checking if hashtagProfileIdsKey exists")
		return nil, fmt.Errorf("error checking if hashtagProfileIdsKey exists: %w", err)
	}
	if exists == 0 {
		h.logger.Info("hashtagProfileIdsKey does not exist")
		return nil, fmt.Errorf("hashtagProfileIdsKey does not exist")
	}

	// Get cache list profileIds
	profileIds, err := cachier.redisClient.SMembers(ctx, hashtagProfileIdsKey).Result()
	if err != nil {
		h.logger.Info("error getting profileIds from cache")
		return nil, fmt.Errorf("error getting profileIds from cache: %w", err)
	}

	// // If miss cache
	// if len(profileIds) == 0 {
	// 	h.logger.Info("miss cache")
	// 	return nil, ErrCacheMiss
	// }

	return profileIds, nil
}

func (h HashtagCache) GetProfileIdHashtagsCacheElements(ctx context.Context, profileId string) ([]string, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return nil, fmt.Errorf("cachier is not redis")
	}

	// get key cache
	profileIdhashtagsKey := h.getProfileHashtagsCacheKey(profileId)

	// Check key exists
	exists, err := cachier.redisClient.Exists(ctx, profileIdhashtagsKey).Result()
	if err != nil {
		h.logger.Sugar().Infof("error checking if profileIdhashtagsKey: %s exists", profileIdhashtagsKey)
		return nil, fmt.Errorf("error checking if profileIdhashtagsKey exists: %w", err)
	}
	if exists == 0 {
		h.logger.Sugar().Infof("profileIdhashtagsKey: %s does not exist", profileIdhashtagsKey)
		return []string{}, nil //fmt.Errorf("profileIdhashtagsKey does not exist")
	}

	// Get cache list hashtags
	hashtags, err := cachier.redisClient.SMembers(ctx, profileIdhashtagsKey).Result()
	if err != nil {
		h.logger.Info("error getting hashtags from cache")
		return nil, fmt.Errorf("error getting hashtags from cache: %w", err)
	}

	// // If miss cache
	// if len(hashtags) == 0 {
	// 	h.logger.Info("miss cache")
	// 	return nil, ErrCacheMiss
	// }

	return hashtags, nil
}

func (h HashtagCache) DeleteProfileIdsHashtagProfileCache(ctx context.Context, hashtags []string, profileIds []string) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	for _, hashtag := range hashtags {
		// get key cache
		hashtagProfileKey := h.getHashtagProfilesCacheKey(hashtag)
		err := cachier.redisClient.SRem(ctx, hashtagProfileKey, util.ConvertSliceStringToSliceInterface(profileIds)...).Err()
		if err != nil {
			h.logger.With(zap.Error(err)).Error("error DeleteHashtagLeaderboardCacheElements")
			return false, fmt.Errorf("error deleting elements from hashtag leaderboard cache: %w", err)
		}
	}

	for _, profileId := range profileIds {
		// get key cache
		profileIdHashtagsKey := h.getProfileHashtagsCacheKey(profileId)
		err := cachier.redisClient.SRem(ctx, profileIdHashtagsKey, util.ConvertSliceStringToSliceInterface(hashtags)...).Err()
		if err != nil {
			h.logger.With(zap.Error(err)).Error("error DeleteHashtagLeaderboardCacheElements")
			return false, fmt.Errorf("error deleting elements from hashtag leaderboard cache: %w", err)
		}
	}

	return true, nil
}

func (h HashtagCache) DeleteProfileHashtagsCacheKey(ctx context.Context, profileId string) (bool, error) {
	cachier, ok := h.cachier.(*redisClient)
	if !ok {
		h.logger.Info("cachier is not redis")
		return false, fmt.Errorf("cachier is not redis")
	}

	// get key cache
	hashtagProfileKey := h.getProfileHashtagsCacheKey(profileId)
	err := cachier.redisClient.Del(ctx, hashtagProfileKey, hashtagProfileKey).Err()
	if err != nil {
		h.logger.With(zap.Error(err)).Error("error DeleteHashtagLeaderboardCacheElements")
		return false, fmt.Errorf("error deleting elements from hashtag leaderboard cache: %w", err)
	}

	return true, nil
}
