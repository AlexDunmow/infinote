ROLLBACK;
BEGIN;

CREATE TABLE companies (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE groups (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL,
    company_id UUID NOT NULL REFERENCES companies(id),
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT UNIQUE NOT NULL,
    permissions TEXT[] NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    role_id UUID NOT NULL REFERENCES roles(id),
    company_id UUID NOT NULL REFERENCES companies(id),
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verify_token TEXT NOT NULL DEFAULT gen_random_uuid (),
    require_old_password BOOLEAN NOT NULL DEFAULT FALSE,
    reset_token TEXT NOT NULL DEFAULT gen_random_uuid (),
    password_hash TEXT NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE issued_tokens (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    company_id uuid NOT NULL REFERENCES companies (id),
    device TEXT NOT NULL,
    token_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    token_expires TIMESTAMPTZ NOT NULL,
    blacklisted BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE notes (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    owner_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    done BOOLEAN NOT NULL DEFAULT false,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    archived_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_groups (
    user_id UUID NOT NULL REFERENCES users(id),
    group_id UUID NOT NULL REFERENCES groups(id),
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, group_id)
);


CREATE TABLE shares (
    user_id UUID NOT NULL REFERENCES users(id),
    note_id UUID NOT NULL REFERENCES notes(id),
    shared_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, note_id)
);

COMMIT;