package posts_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Romasmi/social-network/internal/domain/post"
	"github.com/Romasmi/social-network/internal/metrics"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const cachedFeedLength = 1000
const postLifetime = time.Hour * 24 * 30
const feedLifetime = time.Hour * 24 * 7

type PostsCache struct {
	client *redis.Client
}

func NewPostsCache(client *redis.Client) *PostsCache {
	return &PostsCache{client: client}
}

func (c *PostsCache) PushToFeeds(ctx context.Context, profileIDs []uuid.UUID, posts []*post.Post) error {
	for _, v := range profileIDs {
		err := c.cacheFeed(ctx, v, posts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PostsCache) DeleteFromFeeds(ctx context.Context, profileIDs []uuid.UUID, postIDs []uuid.UUID) error {
	if len(profileIDs) == 0 || len(postIDs) == 0 {
		return nil
	}
	p := c.client.Pipeline()
	for _, v := range profileIDs {
		p.ZRem(ctx, getFeedKey(v), utils.UUIDsToStrings(postIDs))
	}
	cmds, err := p.Exec(ctx)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			// Log error
			fmt.Printf("Redis command error: %v\n", cmd.Err())
		}
	}
	return nil
}

func (c *PostsCache) cacheFeed(ctx context.Context, profileID uuid.UUID, posts []*post.Post) error {
	if len(posts) == 0 {
		return nil
	}
	if len(posts) > cachedFeedLength {
		posts = posts[:cachedFeedLength]
	}
	postIds := make([]string, len(posts))
	p := c.client.Pipeline()
	for i, v := range posts {
		postIds[i] = v.ID.String()
		p.ZAdd(ctx, getFeedKey(profileID), redis.Z{Score: float64(v.ID.Time()), Member: v.ID.String()})
	}
	p.ZRemRangeByRank(ctx, getFeedKey(profileID), 0, -(cachedFeedLength + 1))
	p.Expire(ctx, getFeedKey(profileID), feedLifetime)
	cmds, err := p.Exec(ctx)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			// Log error
			fmt.Printf("Redis command error: %v\n", cmd.Err())
		}
	}
	return nil
}

func (c *PostsCache) cachePosts(ctx context.Context, posts []*post.Post) error {
	p := c.client.Pipeline()
	for _, v := range posts {
		postToCache, err := json.Marshal(v)
		if err != nil {
			return err
		}
		p.Set(ctx, getPostKey(v.ID), postToCache, postLifetime)
	}
	cmds, err := p.Exec(ctx)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			// Log error
			fmt.Printf("Redis command error: %v\n", cmd.Err())
		}
	}
	return nil
}

func (c *PostsCache) getCachedPosts(ctx context.Context, postIds []uuid.UUID) ([]*post.Post, missedCache, error) {
	keys := make([]string, len(postIds))
	for i, v := range postIds {
		keys[i] = getPostKey(v)
	}
	postsCache, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, nil, err
	}
	mc := make(missedCache, 0, len(postIds))
	posts := make([]*post.Post, len(postIds))
	for i, v := range postsCache {
		if v == nil {
			metrics.CacheMissesTotal.WithLabelValues("posts").Inc()
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, postIds[i]})
			continue
		}
		metrics.CacheHitsTotal.WithLabelValues("posts").Inc()
		data, ok := v.(string)
		if !ok {
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, postIds[i]})
			continue
		}
		var post *post.Post
		if err := json.Unmarshal([]byte(data), &post); err != nil {
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, postIds[i]})
			continue
		}
		posts[i] = post
	}
	return posts, mc, nil
}

func (c *PostsCache) deletePost(ctx context.Context, postID uuid.UUID) {
	c.client.Del(ctx, getPostKey(postID))
}

func (c *PostsCache) getPostIdsFromFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	postIds, err := c.client.ZRevRange(ctx, getFeedKey(profileID), int64(offset), int64(limit+offset-1)).Result()
	if err != nil {
		return nil, err
	}
	if len(postIds) > 0 {
		metrics.CacheHitsTotal.WithLabelValues("feed").Inc()
	} else {
		metrics.CacheMissesTotal.WithLabelValues("feed").Inc()
	}
	output := make([]uuid.UUID, len(postIds))
	for i, v := range postIds {
		output[i] = uuid.MustParse(v)
	}
	return output, nil
}
