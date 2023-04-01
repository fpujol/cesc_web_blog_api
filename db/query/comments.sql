-- name: GetPostComments :many
SELECT * FROM comments
WHERE post_id=$1;