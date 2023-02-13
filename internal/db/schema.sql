
ALTER TABLE IF EXISTS ONLY public.user_order DROP CONSTRAINT IF EXISTS user_order_user_null_fk;
ALTER TABLE IF EXISTS ONLY public.user_balance_withdraw DROP CONSTRAINT IF EXISTS user_balance_withdraw_user_null_fk;
ALTER TABLE IF EXISTS ONLY public.user_balance DROP CONSTRAINT IF EXISTS user_balance_user_null_fk;
ALTER TABLE IF EXISTS ONLY public."user" DROP CONSTRAINT IF EXISTS user_pk;
ALTER TABLE IF EXISTS ONLY public.user_order DROP CONSTRAINT IF EXISTS user_order_pk;
ALTER TABLE IF EXISTS ONLY public.user_balance DROP CONSTRAINT IF EXISTS user_balance_pk;
DROP TABLE IF EXISTS public.user_order;
DROP TABLE IF EXISTS public.user_balance_withdraw;
DROP TABLE IF EXISTS public.user_balance;
DROP TABLE IF EXISTS public."user";

CREATE TABLE IF NOT EXISTS "user" (
                                      id uuid DEFAULT gen_random_uuid() NOT NULL,
                                      login character varying(255) NOT NULL,
                                      password character varying NOT NULL,
                                      token character varying
);

CREATE TABLE IF NOT EXISTS user_balance (
                                            user_id uuid,
                                            current double precision NOT NULL,
                                            withdrawn double precision
);

CREATE TABLE IF NOT EXISTS user_balance_withdraw (
                                                     user_id uuid NOT NULL,
                                                     order_number character varying(255) NOT NULL,
                                                     sum double precision NOT NULL,
                                                     processed_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_order (
                                          user_id uuid NOT NULL,
                                          number character varying(255) NOT NULL,
                                          status character varying(255) NOT NULL,
                                          accrual double precision,
                                          uploaded_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE IF EXISTS ONLY public.user_balance
    ADD CONSTRAINT user_balance_pk UNIQUE (user_id);

ALTER TABLE IF EXISTS ONLY public.user_order
    ADD CONSTRAINT user_order_pk UNIQUE (number);

ALTER TABLE IF EXISTS ONLY public."user"
    ADD CONSTRAINT user_pk PRIMARY KEY (id);

ALTER TABLE IF EXISTS ONLY public.user_balance
    ADD CONSTRAINT user_balance_user_null_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE IF EXISTS ONLY public.user_balance_withdraw
    ADD CONSTRAINT user_balance_withdraw_user_null_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE IF EXISTS ONLY public.user_order
    ADD CONSTRAINT user_order_user_null_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);


