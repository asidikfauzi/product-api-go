# Project Documentation

## Prerequisites

This project requires PostgreSQL with the `uuid-ossp` extension enabled and Redis. Follow the steps below to ensure your environment is set up correctly:

### Install PostgreSQL
Make sure PostgreSQL is installed on your system. If not, you can download and install it from [PostgreSQL official site](https://www.postgresql.org/download/).

### Enable `uuid-ossp` Extension
After installing PostgreSQL, enable the `uuid-ossp` extension by running the following SQL command in your database:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

Ensure you have sufficient privileges (e.g., superuser) to create extensions.

### Verify Extension
You can verify that the extension is enabled by running:

```sql
SELECT uuid_generate_v4();
```

This should return a randomly generated UUID.

### Install Redis
Redis is required for caching and other performance optimizations. Install Redis using the following commands:

- **MacOS (Homebrew):**
  ```sh
  brew install redis
  ```
  Start Redis with:
  ```sh
  brew services start redis
  ```

- **Ubuntu/Debian:**
  ```sh
  sudo apt update
  sudo apt install redis-server
  ```
  Start Redis with:
  ```sh
  sudo systemctl start redis
  ```

- **Windows (using WSL or Redis for Windows):**
  Download and install Redis from [Redis official site](https://redis.io/download/).

Verify Redis is running by executing:
```sh
redis-cli ping
```
If Redis is working correctly, it should return `PONG`.

## Project Structure

The project uses a `Makefile` to simplify various build, migration, and execution tasks. Below is an explanation of the main commands:

### Makefile Commands

#### General Commands

- **`all`**: Installs dependencies, builds the main application, and runs it.
  ```sh
  make all
  ```
- **`build`**: Builds the main application binary.
  ```sh
  make build
  ```
- **`reload`**: Rebuilds the application binary and runs it.
  ```sh
  make reload
  ```

#### Database Migration Commands

- **`migrate`**: Builds and runs the migration script to apply new migrations.
  ```sh
  make migrate
  ```
- **`rollback`**: Rolls back migrations by a specified number of steps.
  ```sh
  make rollback step=<number_of_steps>
  ```

### Additional Details

- **Go Modules**
  Before running any command, ensure Go dependencies are installed using:
  ```sh
  make mod
  ```

- **Custom Parameters**
  Some commands, like `rollback` and `new-migration`, require additional parameters:
    - `step`: Number of migration steps to roll back.
    - `table`: Name of the table for the new migration.

## Running the Application

1. Build the application binary:
   ```sh
   make build
   ```
2. Run the application:
   ```sh
   make run
   ```

## Running the Application with Docker
This project includes Docker support for simplified deployment and development. Follow the steps below to run the application using Docker.
Ensure you have Docker installed on your system. If not, download and install it from [Docker official site](https://www.docker.com/).

### Running with Docker Compose

1. Build and Start Containers
   ```sh
   make docker-build
   ```
2. Stopping the Containers
   ```sh
   make docker-down
   ```
3. Checking Running Containers
   ```sh
   docker ps
   ```

### Running Database Migrations in Docker
1. Run Migrations Inside the Container
   ```sh
   docker compose exec <container_id> make migrate
   ```
2. Roll Back Migrations
   ```sh
   docker compose exec <container_id> make rollback step=1
   ```

## Notes

- Ensure your `.env` file is correctly set up before running Docker.
- The `POSTGRES_HOST` in `.env` should be set to `postgres` (not `localhost`) when using Docker.
- Use `docker logs <container_name>` to check logs if something isn't working.