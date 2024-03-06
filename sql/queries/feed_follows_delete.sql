-- name: DeleteSpecificFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;