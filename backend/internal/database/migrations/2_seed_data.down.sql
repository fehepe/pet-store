-- Remove seed data

-- Delete all pets from merchant1's store
DELETE FROM pets WHERE store_id = '123e4567-e89b-12d3-a456-426614174000';

-- Delete merchant1's store (this will cascade delete related data)
DELETE FROM stores WHERE id = '123e4567-e89b-12d3-a456-426614174000';