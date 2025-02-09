-- Create "subscription_configs" table
CREATE TABLE "public"."subscription_configs" (
  "id" bigserial NOT NULL,
  "provider" text NOT NULL,
  "description" text NULL,
  "logo" text NULL,
  "website" text NULL,
  "status" text NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_subscription_configs_provider" to table: "subscription_configs"
CREATE UNIQUE INDEX "idx_subscription_configs_provider" ON "public"."subscription_configs" ("provider");
-- Create "subscriptions" table
CREATE TABLE "public"."subscriptions" (
  "id" bigserial NOT NULL,
  "uuid" text NOT NULL,
  "name" text NOT NULL,
  "price" numeric NOT NULL,
  "billing_cycle" bigint NOT NULL,
  "start_date" timestamptz NOT NULL,
  "logo" text NULL,
  "user_id" text NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "subscription_config_plans" table
CREATE TABLE "public"."subscription_config_plans" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "price" numeric NOT NULL,
  "currency" text NOT NULL,
  "billing_cycle" integer NOT NULL,
  "status" text NULL DEFAULT 'active',
  "subscription_config_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subscription_configs_plans" FOREIGN KEY ("subscription_config_id") REFERENCES "public"."subscription_configs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
