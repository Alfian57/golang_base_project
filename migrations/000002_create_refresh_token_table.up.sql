CREATE TABLE "refresh_tokens" (
    "id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "token_hash" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "expires_at" BIGINT NOT NULL
);

ALTER TABLE
    "refresh_tokens" ADD PRIMARY KEY("id");

ALTER TABLE
    "refresh_tokens" ADD CONSTRAINT "refresh_tokens_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;