-- name: CreatePost :one
INSERT INTO posts (post_id,
                  slug,
                  title,
                  introduction,
                  post_category_id,
                  main_image_alt, 
                  main_image_path, 
                  thumbnail_image_alt, 
                  thumbnail_image_path, 
                  content, 
                  author, 
                  author_image_path, 
                  author_image_alt,
                  published,
                  published_at,
                  published_by, 
                  created_at, 
                  created_by) 
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) 
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
  set post_category_id=$2,
  slug=$3,
  title=$4, 
  introduction=$5,
  main_image_alt=$6, 
  thumbnail_image_alt=$7, 
  content=$8,
  author=$9, 
  author_image_path=$10, 
  author_image_alt=$11,
  published=$12,
  published_at=$13,
  published_by=$14,
  updated_at=$15,
  updated_by=$16
WHERE post_id = $1
RETURNING *;

-- name: UpdateMainImagePost :one
UPDATE posts
  set main_image_path=$2, 
  updated_at=$3,
  updated_by=$4
WHERE post_id = $1
RETURNING *;

-- name: UpdateThumbnailImagePost :one
UPDATE posts
  set thumbnail_image_path=$2, 
  updated_at=$3,
  updated_by=$4
WHERE post_id = $1
RETURNING *;


-- name: ListPosts :many
SELECT * FROM posts 
ORDER BY title;

-- name: GetPost :one
SELECT * FROM posts 
WHERE slug=$1 LIMIT 1;

-- name: GetPublicPosts :many
SELECT * FROM posts
WHERE published=$1 order by published_at;
