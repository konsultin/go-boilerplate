CREATE TABLE IF NOT EXISTS auth_session (
    id BIGSERIAL PRIMARY KEY,
    xid VARCHAR(255) NOT NULL UNIQUE,
    subject_id VARCHAR(255) NOT NULL,
    subject_type_id INT NOT NULL,
    auth_provider_id INT NOT NULL DEFAULT 1,
    device_platform_id INT NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    device JSONB NOT NULL DEFAULT '{}',
    notification_channel_id INT NOT NULL DEFAULT 0,
    notification_token VARCHAR(512),
    expired_at TIMESTAMP NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by JSONB,
    version BIGINT NOT NULL DEFAULT 1,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_auth_session_xid ON auth_session(xid);
CREATE INDEX IF NOT EXISTS idx_auth_session_subject_id ON auth_session(subject_id);
CREATE INDEX IF NOT EXISTS idx_auth_session_status ON auth_session(status_id);
CREATE INDEX IF NOT EXISTS idx_auth_session_expired_at ON auth_session(expired_at);
