CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,
  email VARCHAR(320) UNIQUE NOT NULL,
  master_password VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL,
  created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_vault (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  vault_id VARCHAR(255),
  user_id VARCHAR(255),
  CONSTRAINT fk_vault FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE vaults (
  id VARCHAR(255) PRIMARY KEY,
  key VARCHAR(255) NOT NULL
  owner_id VARCHAR(255) NOT NULL,
  CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE items (
  id VARCHAR(255) PRIMARY KEY,
  vault_id VARCHAR(255) NOT NULL,
  name VARCHAR(1000) NOT NULL,
  note VARCHAR(10000),
  type SMALLINT NOT NULL,
  created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_vault FOREIGN KEY(vault_id) REFERENCES vaults(id)
);

CREATE TABLE logins (
  item_id VARCHAR(255) PRIMARY KEY,
  username VARCHAR(1000),
  password VARCHAR(1000),
  password_revision_date TIMESTAMPTZ
  CONSTRAINT fk_item FOREIGN KEY(item_id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE TABLE uris (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  login_id VARCHAR(255) NOT NULL,
  uri VARCHAR(10000),
  CONSTRAINT fk_login FOREIGN KEY(login_id) REFERENCES logins(item_id) ON DELETE CASCADE
);

CREATE TABLE password_history (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  login_id VARCHAR(255) NOT NULL,
  password VARCHAR(10000),
  last_used TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_login FOREIGN KEY(login_id) REFERENCES logins(item_id) ON DELETE CASCADE
);