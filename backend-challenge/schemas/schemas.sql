CREATE SCHEMA kart;

CREATE TABLE IF NOT EXISTS kart.products (
      id          BIGSERIAL PRIMARY KEY,
      name        Varchar(255) NOT NULL,
      category    Varchar(100) NOT NULL,
      price       NUMERIC(10, 2) NOT NULL,
      status      Varchar(20) NOT NULL,
      image       JSONB NOT NULL,
      meta        JSONB,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      UNIQUE (name, category)
);

CREATE INDEX IF NOT EXISTS idx_products_created_at ON kart.products(created_at DESC);

CREATE TABLE IF NOT EXISTS kart.orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_code VARCHAR(20),
    subtotal    NUMERIC(10, 2),
    discount    NUMERIC(10, 2) DEFAULT 0,
    total       NUMERIC(10, 2),
    meta        JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS kart.order_items (
     id          BIGSERIAL PRIMARY KEY,
     order_id    UUID NOT NULL REFERENCES kart.orders(id) ON DELETE CASCADE,
     product_id  BIGINT NOT NULL REFERENCES kart.products(id),
     quantity    INTEGER NOT NULL CHECK (quantity > 0),
     unit_price  NUMERIC(10, 2) NOT NULL,
     discount    NUMERIC(10, 2) DEFAULT 0,
     price       NUMERIC(10, 2) NOT NULL,
     meta        JSONB,
     created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     UNIQUE (order_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON kart.order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON kart.orders(created_at DESC);

CREATE TABLE IF NOT EXISTS kart.coupons (
    code         VARCHAR(10) PRIMARY KEY,
    file_sources varchar(20)[] NOT NULL,
    file_count   INT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_coupons_file_count ON kart.coupons(file_count);