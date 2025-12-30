package posts_repository

import "github.com/google/uuid"

func getFeedKey(profileID uuid.UUID) string {
	return "feed:" + profileID.String()
}

func getPostKey(postID uuid.UUID) string {
	return "posts:" + postID.String()
}
