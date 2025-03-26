-- name: CreateBanner :one
INSERT INTO banners (
    title, image_url, target_url, is_active, position
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBanner :one
SELECT * FROM banners
WHERE id = $1 LIMIT 1;

-- name: ListBanners :many
SELECT * FROM banners
ORDER BY created_at DESC;

-- name: ListActiveBanners :many
SELECT * FROM banners
WHERE is_active = true
ORDER BY created_at DESC;

-- name: UpdateBanner :one
UPDATE banners
SET title = $2, 
    image_url = $3, 
    target_url = $4, 
    is_active = $5, 
    position = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBanner :exec
DELETE FROM banners
WHERE id = $1;
