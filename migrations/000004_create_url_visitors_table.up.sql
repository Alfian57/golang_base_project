CREATE TABLE "url_visitors"(
    "id" UUID NOT NULL,
    "url_id" UUID NOT NULL,
    "ip_address" VARCHAR(15) NOT NULL,
    "user_agent" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "url_visitors" ADD PRIMARY KEY("id");
ALTER TABLE
    "url_visitors" ADD CONSTRAINT "url_visitors_url_id_foreign" FOREIGN KEY("url_id") REFERENCES "urls"("id") ON DELETE CASCADE;