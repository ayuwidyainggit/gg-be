CREATE TABLE IF NOT EXISTS sys.password_reset_requests (
    id serial4 NOT NULL,
    user_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,

    otp_code VARCHAR(255) NOT NULL,
    otp_expired_at TIMESTAMP NOT NULL,
    otp_attempt_count INT NOT NULL DEFAULT 0,
    otp_max_attempt INT NOT NULL DEFAULT 3,

    resend_count INT NOT NULL DEFAULT 0,
    resend_max INT NOT NULL DEFAULT 3,
    resend_cooldown_until TIMESTAMP,

    request_id VARCHAR(100) NOT NULL,

    reset_token VARCHAR(255),
    reset_token_expired_at TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

	CONSTRAINT password_reset_requests_pkey PRIMARY KEY (id),
    CONSTRAINT password_reset_requests_request_id_key UNIQUE (request_id)
);
