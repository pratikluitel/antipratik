# scripts

## create_user.sh

Creates or refreshes a user in the antipratik SQLite database. On success, prints the user's JWT token — copy it to use with protected API routes.

### Usage

```bash
# Create a new user
./scripts/create_user.sh --db ./data/antipratik.db --username admin --password secret

# Refresh the JWT token for an existing user (no password change)
./scripts/create_user.sh --db ./data/antipratik.db --username admin --refresh
```

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--db` | yes | Path to the SQLite database file |
| `--username` | yes | Username |
| `--password` | unless `--refresh` | Password (bcrypt-hashed before storage) |
| `--refresh` | no | Skip password change; only issue a new token for an existing user |

### Notes

- The JWT secret is generated automatically on first use and stored in the database — no manual configuration needed.
- Tokens expire after **7 days**. Run with `--refresh` to issue a new one.
- Creating a user that already exists is an error — use `--refresh` instead.
