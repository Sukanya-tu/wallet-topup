-- USERS TABLE
CREATE TABLE IF NOT EXISTS public.users (
    user_id bigint NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    balance numeric(12,2),
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;


-- TRANSACTIONS TABLE
CREATE TABLE IF NOT EXISTS public.transactions (
    transaction_id uuid NOT NULL,
    user_id bigint,
    amount numeric(12,2) NOT NULL,
    payment_method text COLLATE pg_catalog."default",
    status text COLLATE pg_catalog."default",
    expires_at timestamp with time zone NOT NULL,
    id uuid,
    CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id),
    CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transactions_status_check CHECK (status = ANY (ARRAY['verified'::character varying::text, 'completed'::character varying::text]))
);

ALTER TABLE IF EXISTS public.transactions
    OWNER to postgres;
