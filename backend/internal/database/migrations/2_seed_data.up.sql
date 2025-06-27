-- Insert seed data for testing

-- Insert merchant1's store
INSERT INTO stores (id, name, owner_id, created_at) 
VALUES (
    '123e4567-e89b-12d3-a456-426614174000',
    'Pet Paradise Store',
    'merchant1',
    CURRENT_TIMESTAMP
) ON CONFLICT (owner_id) DO NOTHING;

-- Insert sample pets with images
INSERT INTO pets (id, store_id, name, species, age, picture_url, description, breeder_name, breeder_email_encrypted, status, created_at) VALUES
-- Cats
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Luna',
    'Cat',
    2,
    'https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300&h=200&fit=crop&crop=center',
    'Beautiful black and white cat with striking green eyes. Very playful and affectionate.',
    'Whiskers Cattery',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '1 day'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Milo',
    'Cat',
    3,
    'https://images.unsplash.com/photo-1574144611937-0df059b5ef3e?w=300&h=200&fit=crop&crop=center',
    'Orange tabby cat with a calm and gentle personality. Perfect for families.',
    'Feline Friends',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '2 hours'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Shadow',
    'Cat',
    1,
    'https://images.unsplash.com/photo-1596854407944-bf87f6fdd49e?w=300&h=200&fit=crop&crop=center',
    'Young black kitten with bright yellow eyes. Full of energy and curiosity.',
    'Midnight Cattery',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '6 hours'
),

-- Dogs
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Max',
    'Dog',
    4,
    'https://images.unsplash.com/photo-1552053831-71594a27632d?w=300&h=200&fit=crop&crop=center',
    'Friendly Golden Retriever who loves to play fetch and swim. Great with children.',
    'Golden Valley Kennels',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '3 hours'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Bella',
    'Dog',
    2,
    'https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=300&h=200&fit=crop&crop=center',
    'Sweet and loyal Labrador mix. Well-trained and loves long walks.',
    'Happy Tails Rescue',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '5 hours'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Rocky',
    'Dog',
    5,
    'https://images.unsplash.com/photo-1518717758536-85ae29035b6d?w=300&h=200&fit=crop&crop=center',
    'Strong and protective German Shepherd. Excellent guard dog and family companion.',
    'Elite K9 Breeders',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '8 hours'
),

-- Frogs
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Kermit',
    'Frog',
    1,
    'https://images.ctfassets.net/cnu0m8re1exe/4txgybYHtUH0z6Dy9IIFGH/e9f4612ef512ae7ad8a580a39557e4d2/Glass_Frog.jpg?fm=jpg&fl=progressive&w=660&h=433&fit=fill',
    'Adorable green tree frog with a cheerful personality. Easy to care for.',
    'Amphibian Adventures',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '4 hours'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Lily',
    'Frog',
    2,
    'https://www.pbs.org/wnet/nature/files/2021/05/frog-1280x720.png',
    'Beautiful poison dart frog with vibrant colors. Requires specialized care.',
    'Exotic Amphibians Co',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '7 hours'
),

(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Snowball',
    'Cat',
    1,
    NULL,
    'Pure white kitten with blue eyes. Very affectionate and loves to cuddle.',
    'Arctic Cats Breeder',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '1 hour'
),
(
    gen_random_uuid(),
    '123e4567-e89b-12d3-a456-426614174000',
    'Buddy',
    'Dog',
    3,
    NULL,
    'Energetic Border Collie mix. Highly intelligent and loves to learn new tricks.',
    'Smart Paws Training',
    'encrypted_email_data_here',
    'available',
    CURRENT_TIMESTAMP - INTERVAL '30 minutes'
);