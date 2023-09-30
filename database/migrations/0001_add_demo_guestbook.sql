-- +up  <- SQL below runs when applying a migration
CREATE TABLE demo_guestbook (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    message text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

-- Add some demo messages
INSERT INTO demo_guestbook (name, message, created_at) VALUES
    ('John Doe', 'Hello world!', NOW()),
    ('Jane Doe', 'Hi there!', NOW());


-- +down <- SQL below runs when reverting a migration
DROP TABLE demo_guestbook;