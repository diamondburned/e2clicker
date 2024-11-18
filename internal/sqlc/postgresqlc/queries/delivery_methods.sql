-- Language: postgresql
--
/*
 * Delivery Method
 */
-- name: DeliveryMethods :many
SELECT *
FROM delivery_methods;

-- name: DeliveryMethod :one
SELECT *
FROM delivery_methods
WHERE name = $1;
