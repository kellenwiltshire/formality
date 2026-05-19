# Part 2 - Database

So I've decided to go with Postgres for my DB, but I need something to help manage it within the Go backend. I've decided to go with a tool called Goose. Seems like a lightweight solution to get me going. First though, designing the database schema.

## Schema

I've decided to start with 4 tables: `users, forms, form_submissions, and smtp_settings`. They'll be relational, so a `user` can have many `forms`, `forms` can obviously have many `form_submission`, and a `user` can have 1 `smtp_setting`.

Here's how the tables look in code, for now:

```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Store securely hashed passwords
    role user_role DEFAULT 'user' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE smtp_settings (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    username VARCHAR(255),
    password_encrypted TEXT, -- Encrypt this data before saving!
    encryption_type VARCHAR(50) DEFAULT 'tls', -- tls, ssl, or none
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE forms (
    id VARCHAR(6) PRIMARY KEY DEFAULT generate_short_id(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    target_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE form_submissions (
    id SERIAL PRIMARY KEY,
    form_id VARCHAR(6) NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    status form_status DEFAULT 'received' NOT NULL
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

I also added these roles:

```
CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TYPE form_status AS ENUM ('received', 'processed', 'error')
```

I had some help from `AI` to make sure my syntax was correct, and it also drew up this helper function for generating the `form.id`:

```
CREATE OR REPLACE FUNCTION generate_short_id()
RETURNS TEXT AS $$
DECLARE
    -- Your character set: 26 upper + 26 lower + 10 numbers = 62 characters
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    result TEXT := '';
    i INT;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql VOLATILE;
```

Now I have to actually create a postgres instance for all this to live in. Since I am using Docker anyways, I can start of my docker files for this.

This should be a good start

```
services:
  database:
    image: postgres:17
    networks:
      - internal
      - default
    ports:
      - 15432:5432
    environment:
      - PUID=${PUID}
      - PGID=${PGID}
      - TZ=${TZ}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_USER=${DB_USERNAME}
      - POSTGRES_DB=${DB_DATABASE_NAME}
    volumes:
      - ${PWD}/db-data/:/var/lib/postgresql/data/

networks:
  default:
  internal:
    internal: true
```

Success! The container is created and running. Now we'll need to create the tables with our migrations, might as well have that as the first step in our backend process.

```
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func goDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func main() {
	// Load DB Env Variables
	dbUser := goDotEnvVariable("DB_USERNAME")
	dbPass := goDotEnvVariable("DB_PASSWORD")
	dbHost := goDotEnvVariable("DB_HOST")
	dbPort := goDotEnvVariable("DB_PORT")
	dbName := goDotEnvVariable("DB_DATABASE_NAME")

	// Format: postgres://username:password@host:port/dbname?options
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

    // Connect to database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Run migrations automatically
    if err := goose.Up(db, "migrations"); err != nil {
        log.Fatal(err)
    }

    log.Println("Migrations completed successfully!")
}
```

Little help from Medium here which was great. Loading env's is a bit more of a process than expected, but interesting. Migrations completed, and the tables have been created! I think this is a good time for a first commit. Next up, starting to build the API!
