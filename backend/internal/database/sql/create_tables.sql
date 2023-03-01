CREATE TYPE batches_status AS ENUM ('CREATED', 'PROCESSING', 'ERROR', 'PROCESSED');

CREATE TABLE "batches" (
  "id" uuid NOT NULL PRIMARY KEY,
  "affiliate_id" uuid NOT NULL,
  "batch_raw" varchar NOT NULL,
  "status" batches_status NOT NULL,
  "errors" jsonb NULL DEFAULT '[]',
  "created_at" timestamp DEFAULT NOW()
);

CREATE TABLE "transactions" (
  "batch_id" uuid NOT NULL,
  "affiliate_id" uuid NOT NULL,
  "type" smallint NOT NULL,
  "date" timestamp NOT NULL,
  "product" varchar NOT NULL,
  "value" decimal(8,2) NOT NULL,
  "seller" varchar NOT NULL,
  CONSTRAINT fk_batch FOREIGN KEY(batch_id) REFERENCES batches(id)
);

CREATE INDEX idx_batch ON transactions using btree (batch_id);