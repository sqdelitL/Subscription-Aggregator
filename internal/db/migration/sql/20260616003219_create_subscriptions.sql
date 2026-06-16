-- migrate:up
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    service_name VARCHAR(200),
    price INT,
    start_date DATE,
    end_date DATE
);

-- migrate:down
DROP TABLE subscriptions;
