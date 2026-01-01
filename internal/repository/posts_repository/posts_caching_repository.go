package posts_repository

import (
	"context"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/post"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type missedCache []struct {
	i  int
	id uuid.UUID
}

type PostsCachingRepository struct {
	postsRepo  PostsRepository
	postsCache *PostsCache
}

func (c *PostsCachingRepository) GetPostsByIds(ctx context.Context, postIDs []uuid.UUID) ([]*post.Post, error) {
	return c.postsRepo.GetPostsByIds(ctx, postIDs)
}

func CreateCachedPostsRepository(postsRepo PostsRepository, rds *redis.Client) PostsRepository {
	return &PostsCachingRepository{postsRepo: postsRepo, postsCache: NewPostsCache(rds)}
}

func (c *PostsCachingRepository) CreatePost(ctx context.Context, p *post.Post) (*post.Post, error) {
	created, err := c.postsRepo.CreatePost(ctx, p)
	if err != nil {
		return nil, err
	}
	return created, err
}

func (c *PostsCachingRepository) UpdatePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID, newText string) (*post.Post, error) {
	updated, err := c.postsRepo.UpdatePost(ctx, postID, profileID, newText)
	if err != nil {
		return nil, err
	}
	go func() {
		err := c.postsCache.cachePosts(context.WithoutCancel(ctx), []*post.Post{updated})
		if err != nil {
			// TODO Add logger
			fmt.Println(err)
		}
	}()
	return updated, err
}

func (c *PostsCachingRepository) DeletePost(ctx context.Context, postID uuid.UUID, profileID uuid.UUID) error {
	c.postsCache.deletePost(ctx, postID)
	return c.postsRepo.DeletePost(ctx, postID, profileID)
}

func (c *PostsCachingRepository) GetPost(ctx context.Context, postID uuid.UUID) (*post.Post, error) {
	cachedPost, _, err := c.postsCache.getCachedPosts(ctx, []uuid.UUID{postID})
	if err != nil {
		return nil, err
	}
	if len(cachedPost) > 0 {
		return cachedPost[0], nil
	}
	return c.postsRepo.GetPost(ctx, postID)
}

func (c *PostsCachingRepository) GetFeed(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*post.Post, error) {
	if offset+limit > cachedFeedLength {
		return c.postsRepo.GetFeed(ctx, profileID, limit, offset)
	}

	postIds, err := c.postsCache.getPostIdsFromFeed(ctx, profileID, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(postIds) > 0 {
		posts, missedCache, err := c.postsCache.getCachedPosts(ctx, postIds)
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
		results := make(map[uuid.UUID]*post.Post, len(missedCachePosts))
		for _, v := range missedCachePosts {
			results[v.ID] = v
		}
		for _, v := range missedCache {
			posts[v.i] = results[v.id]
		}
		if len(missedCachePosts) > 0 {
			go func() {
				err := c.postsCache.cachePosts(context.WithoutCancel(ctx), missedCachePosts)
				if err != nil {
					// TODO Add logger
					fmt.Println(err)
				}
			}()
		}
		return posts, nil
	}
	fullFeed, err := c.postsRepo.GetFeed(ctx, profileID, cachedFeedLength, 0)
	if err != nil {
		return nil, err
	}

	go func() {
		err := c.postsCache.cacheFeed(context.WithoutCancel(ctx), profileID, fullFeed)
		if err != nil {
			fmt.Println(err)
		}
	}()

	end := offset + limit
	if end > len(fullFeed) {
		end = len(fullFeed)
	}
	if offset >= len(fullFeed) {
		return []*post.Post{}, nil
	}
	return fullFeed[offset:end], nil
}
