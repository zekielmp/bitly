CREATE TABLE product_reviews (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT      NOT NULL,
    user_id     BIGINT      NOT NULL,
    data        TEXT        NOT NULL,
    rating      INTEGER     NOT NULL CHECK (rating BETWEEN 1 AND 5),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_product_reviews_deleted_at ON product_reviews (deleted_at);
CREATE INDEX idx_product_reviews_product_id ON product_reviews (product_id);
CREATE INDEX idx_product_reviews_user_id    ON product_reviews (user_id);