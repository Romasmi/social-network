package posts_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const cachedFeedLength = 1000
const postLifetime = time.Hour * 24 * 30
const feedLifetime = time.Hour * 24 * 7

type missedCache []struct {
	i  int
	id uuid.UUID
}

type cachedPostsRepositoryImpl struct {
	postsRepo PostsRepository
	redis     *redis.Client
}

func (c *cachedPostsRepositoryImpl) GetPostsByIds(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error) {
	return c.postsRepo.GetPostsByIds(ctx, postIDs)
}

func CreateCachedPostsRepository(postsRepo PostsRepository, rds *redis.Client) PostsRepository {
	return &cachedPostsRepositoryImpl{postsRepo: postsRepo, redis: rds}
}

func (c *cachedPostsRepositoryImpl) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	created, err := c.postsRepo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	// TODO push post to friends feeds
	return created, err
}

func (c *cachedPostsRepositoryImpl) UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*models.Post, error) {
	updated, err := c.postsRepo.UpdatePost(ctx, postID, profileID, newText)
	if err != nil {
		return nil, err
	}
	go func() {
		err := c.cachePosts(context.WithoutCancel(ctx), []*models.Post{updated})
		if err != nil {
			// TODO Add logger
			fmt.Println(err)
		}
	}()
	return updated, err
}

func (c *cachedPostsRepositoryImpl) DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error {
	c.redis.Del(ctx, getPostKey(postID))
	// TODO remove from friends feeds
	return c.postsRepo.DeletePost(ctx, postID, profileID)
}

func (c *cachedPostsRepositoryImpl) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post, _, err := c.getCachedPosts(ctx, []string{postID.String()})
	if err != nil {
		return nil, err
	}
	if len(post) > 0 {
		return post[0], nil
	}
	return c.postsRepo.GetPost(ctx, postID)
}

func (c *cachedPostsRepositoryImpl) GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	if offset+limit > cachedFeedLength {
		return c.postsRepo.GetFeed(ctx, profileID, limit, offset)
	}

	postIds, err := c.redis.ZRevRange(ctx, getFeedKey(profileID), int64(offset), int64(limit+offset-1)).Result()
	if err != nil {
		return nil, err
	}
	if len(postIds) > 0 {
		posts, missedCache, err := c.getCachedPosts(ctx, postIds)
		if err != nil {
			return nil, err
		}
		missedCacheIds := make([]uuid.UUID, len(missedCache))
		for i, v := range missedCache {
			missedCacheIds[i] = v.id
		}
		missedCachePosts, err := c.GetPostsByIds(ctx, missedCacheIds)
		if err != nil {
			return nil, err
		}
		results := make(map[uuid.UUID]*models.Post, len(missedCachePosts))
		for _, v := range missedCachePosts {
			results[v.ID] = v
		}
		for _, v := range missedCache {
			posts[v.i] = results[v.id]
		}
		if len(missedCachePosts) > 0 {
			go func() {
				err := c.cachePosts(context.WithoutCancel(ctx), missedCachePosts)
				if err != nil {
					// TODO Add logger
					fmt.Println(err)
				}
			}()
		}
		return posts, nil
	}
	feed, err := c.postsRepo.GetFeed(ctx, profileID, limit, offset)
	go func() {
		err := c.cacheFeed(context.WithoutCancel(ctx), profileID, feed)
		if err != nil {
			// TODO Add logger
			fmt.Println(err)
		}
	}()
	return feed, err
}

func (c *cachedPostsRepositoryImpl) cachePosts(ctx context.Context, posts []*models.Post) error {
	p := c.redis.Pipeline()
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

func (c *cachedPostsRepositoryImpl) cacheFeed(ctx context.Context, profileID uuid.UUID, posts []*models.Post) error {
	if len(posts) == 0 {
		return nil
	}
	if len(posts) > cachedFeedLength {
		posts = posts[:cachedFeedLength]
	}
	postIds := make([]string, len(posts))
	p := c.redis.Pipeline()
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

func (c *cachedPostsRepositoryImpl) getCachedPosts(ctx context.Context, postIds []string) ([]*models.Post, missedCache, error) {
	keys := make([]string, len(postIds))
	for i, v := range postIds {
		keys[i] = getPostKey(uuid.MustParse(v))
	}
	postsCache, err := c.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, nil, err
	}
	mc := make(missedCache, 0, len(postIds))
	posts := make([]*models.Post, len(postIds))
	for i, v := range postsCache {
		if v == nil {
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, uuid.MustParse(postIds[i])})
			continue
		}
		data, ok := v.(string)
		if !ok {
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, uuid.MustParse(postIds[i])})
			continue
		}
		var post *models.Post
		if err := json.Unmarshal([]byte(data), &post); err != nil {
			mc = append(mc, struct {
				i  int
				id uuid.UUID
			}{i, uuid.MustParse(postIds[i])})
			continue
		}
		posts[i] = post
	}
	return posts, mc, nil
}
