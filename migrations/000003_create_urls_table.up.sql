CREATE TABLE "urls"(
    "id" UUID NOT NULL,
    "short_url" VARCHAR(100) NOT NULL,
    "long_url" TEXT NOT NULL,
    "user_id" UUID NOT NULL,
    "expired_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
ALTER TABLE
    "urls" ADD PRIMARY KEY("id");
ALTER TABLE
    "urls" ADD CONSTRAINT "urls_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;