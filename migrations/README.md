# Database migrations

This directory is reserved for versioned PostgreSQL schema migrations.

Add migrations in ordered pairs when a domain requires persistence:

```text
000001_create_example.up.sql
000001_create_example.down.sql
```

Do not add schema changes until the corresponding domain model and repository contract are ready. No migration runner or schema is configured yet.
