CREATE TABLE subscriptions
(
    id           BIGSERIAL PRIMARY KEY,
    service_name TEXT    NOT NULL,
    price        INTEGER NOT NULL CHECK (price >= 0),
    user_id      UUID    NOT NULL,
    start_date   DATE    NOT NULL CHECK (EXTRACT(DAY FROM start_date) = 1),
    end_date     DATE    CHECK (end_date IS NULL OR EXTRACT(DAY FROM end_date) = 1)
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions (service_name);
CREATE INDEX idx_subscriptions_start_date ON subscriptions (start_date);
CREATE INDEX idx_subscriptions_end_date ON subscriptions (end_date);
