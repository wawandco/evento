DELETE FROM events;
DELETE FROM hotel;
DELETE FROM event_hotel_rooms;
DELETE FROM reservations;

INSERT INTO
    events (id, name)
VALUES
    ('7f3535c6-d5cb-44f0-b89b-4b349f01e49d', 'Evento')
ON CONFLICT (id) DO NOTHING;


INSERT INTO hotel (id, name)
VALUES
    ('1c3f3e2a-3f4e-4b62-8f4e-2b3e4f6c7d8e', 'Hotel A'),
    ('2d4e5f6a-4b5e-6c7d-8e9f-0a1b2c3d4e5f', 'Hotel B'),
    ('3e5f6a7b-5c6d-7e8f-9a0b-1c2d3e45f6aa', 'Hotel C')
ON CONFLICT (id) DO NOTHING;


INSERT INTO
    event_hotel_rooms (id, event_id, hotel_id, contracted, locked, reserved)
VALUES
    ('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', '7f3535c6-d5cb-44f0-b89b-4b349f01e49d', '1c3f3e2a-3f4e-4b62-8f4e-2b3e4f6c7d8e', 100, 0, 0),
    ('b2c3d4e5-f6a7-8b9c-0d1e-2f34a5b6c7da', '7f3535c6-d5cb-44f0-b89b-4b349f01e49d', '2d4e5f6a-4b5e-6c7d-8e9f-0a1b2c3d4e5f', 150, 0, 0),
    ('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', '7f3535c6-d5cb-44f0-b89b-4b349f01e49d', '3e5f6a7b-5c6d-7e8f-9a0b-1c2d3e45f6aa', 200, 0, 0)
ON CONFLICT (id) DO NOTHING;
