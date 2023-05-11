CREATE TABLE IF NOT EXISTS public.account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.secret (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stype TEXT NOT NULL,
    title TEXT NOT NULL,
    data_encrypted BYTEA NOT NULL,
    account_id UUID REFERENCES public.account (id),
    UNIQUE(stype, title)
);
