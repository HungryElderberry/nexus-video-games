CREATE TYPE payment_status AS ENUM ('PENDING', 'PAID', 'EXPIRED', 'FAILED');

CREATE TABLE xendit_topups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    xendit_invoice_id VARCHAR(255) UNIQUE NOT NULL,
    amount BIGINT NOT NULL,
    status payment_status DEFAULT 'PENDING' NOT NULL,
    external_checkout_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_xendit_topups_user ON xendit_topups(user_id);