CREATE TABLE boards (
    -- Primary & Identity
    internal_id       BIGSERIAL PRIMARY KEY,
    public_id         UUID NOT NULL DEFAULT gen_random_uuid(),

    -- Data Fields
    title             VARCHAR(255) NOT NULL,
    description       TEXT,

    -- Relations & Audit
    owner_internal_id BIGINT NOT NULL,
    owner_public_id   UUID NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Unique Constraints
    CONSTRAINT boards_public_id_unique UNIQUE (public_id),

    -- Foreign Key Constraints
    CONSTRAINT fk_boards_owner_internal_id
        FOREIGN KEY (owner_internal_id)
        REFERENCES users (internal_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_boards_owner_public_id
        FOREIGN KEY (owner_public_id)
        REFERENCES users (public_id)
        ON DELETE CASCADE
);