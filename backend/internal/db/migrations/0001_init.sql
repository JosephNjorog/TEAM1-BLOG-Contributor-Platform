-- Team1 Blog Contributor Platform — initial schema

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM (
    'super_admin',
    'moderator',
    'graphic_designer',
    'publisher',
    'contributor'
);

CREATE TYPE user_status AS ENUM ('active', 'inactive');

CREATE TYPE article_status AS ENUM (
    'draft',
    'submitted',
    'changes_requested',
    'resubmitted',
    'editorial_approved',
    'banner_uploaded',
    'published',
    'payment_initiated',
    'payment_confirmed'
);

CREATE TYPE review_decision AS ENUM ('approved', 'changes_requested');

CREATE TYPE suggestion_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TYPE payment_status AS ENUM ('pending', 'initiated', 'simulated', 'confirmed', 'failed');

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    role            user_role NOT NULL,
    wallet_address  TEXT,
    bio             TEXT,
    status          user_status NOT NULL DEFAULT 'active',
    invited_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE invitations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL,
    role        user_role NOT NULL,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    invited_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE articles (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contributor_id          UUID NOT NULL REFERENCES users(id),
    reviewer_id             UUID REFERENCES users(id),
    designer_id             UUID REFERENCES users(id),
    publisher_id            UUID REFERENCES users(id),
    title                   TEXT NOT NULL,
    content                 TEXT NOT NULL DEFAULT '',
    source_citation         TEXT,
    status                  article_status NOT NULL DEFAULT 'draft',
    word_count              INTEGER NOT NULL DEFAULT 0,
    review_cycle_count      INTEGER NOT NULL DEFAULT 0,
    substack_url            TEXT,
    cloudinary_banner_url   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at            TIMESTAMPTZ,
    published_at            TIMESTAMPTZ
);

CREATE INDEX idx_articles_contributor ON articles(contributor_id);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_reviewer ON articles(reviewer_id);
CREATE INDEX idx_articles_designer ON articles(designer_id);
CREATE INDEX idx_articles_publisher ON articles(publisher_id);

CREATE TABLE review_cycles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id  UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id),
    decision    review_decision NOT NULL,
    summary     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_review_cycles_article ON review_cycles(article_id);

CREATE TABLE suggestions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id      UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    review_cycle_id UUID NOT NULL REFERENCES review_cycles(id) ON DELETE CASCADE,
    reviewer_id     UUID NOT NULL REFERENCES users(id),
    range_start     INTEGER NOT NULL,
    range_end       INTEGER NOT NULL,
    suggestion_text TEXT NOT NULL,
    status          suggestion_status NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_suggestions_article ON suggestions(article_id);

CREATE TABLE banners (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id      UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    designer_id     UUID NOT NULL REFERENCES users(id),
    cloudinary_url  TEXT NOT NULL,
    uploaded_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    marked_ready_at TIMESTAMPTZ
);

CREATE INDEX idx_banners_article ON banners(article_id);

CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id      UUID NOT NULL REFERENCES articles(id),
    contributor_id  UUID NOT NULL REFERENCES users(id),
    wallet_address  TEXT NOT NULL,
    amount_usd      NUMERIC(10,2) NOT NULL DEFAULT 100.00,
    tx_hash         TEXT,
    status          payment_status NOT NULL DEFAULT 'pending',
    initiated_by    UUID REFERENCES users(id),
    initiated_at    TIMESTAMPTZ,
    confirmed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_payments_article ON payments(article_id);

CREATE TABLE notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    article_id  UUID REFERENCES articles(id),
    message     TEXT NOT NULL,
    read        BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user ON notifications(user_id, read);

CREATE TABLE substack_articles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contributor_id      UUID REFERENCES users(id),
    substack_post_id    TEXT NOT NULL UNIQUE,
    title               TEXT NOT NULL,
    url                 TEXT NOT NULL,
    published_at        TIMESTAMPTZ NOT NULL,
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_substack_articles_contributor ON substack_articles(contributor_id);

CREATE TABLE audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id    UUID REFERENCES users(id),
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id   UUID,
    metadata    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);
