```txt
  _____        __ _             _       
 |_   _|      / _(_)           | |      
   | |  _ __ | |_ _ _ __   ___ | |_ ___ 
   | | | '_ \|  _| | '_ \ / _ \| __/ _ \
  _| |_| | | | | | | | | | (_) | ||  __/
 |_____|_| |_|_| |_|_| |_|\___/ \__\___|

```

_Dependencies_

-   [go](https://golang.org/)
-   [node](https://nodejs.org/en/)
-   [postgres](https://www.postgresql.org/)
-   [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
-   [docker-compose](https://docs.docker.com/compose/install/)

_Included dependent binaries_

-   [caddy](https://github.com/caddyserver/caddy)
-   [sqlboiler](https://github.com/volatiletech/sqlboiler)
-   [migrate](https://github.com/golang-migrate/migrate)
-   [mockery](https://github.com/vektra/mockery)
-   [realize](https://github.com/oxequa/realize)

## Development

**Database**

```bash
docker run -d -p 5438:5432 \
--name infinote-db \
-e POSTGRES_USER=infinote \
-e POSTGRES_PASSWORD=dev \
-e POSTGRES_DB=infinote \
postgres:11-alpine
```

```bash
docker exec -it infinote-db psql -U infinote
```

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
```

**Update Tools (if required)**

```bash
cd server
go generate -tags tools ./tools/...
```

**Web**

```bash
cd web
npm install
npm start
```

**Server**

```bash
cd server
../scripts/db-prepare.sh
go generate

cd graphql
go generate

cd ..

go run cmd/platform/main.go db-drop
go run cmd/platform/main.go db-migrate
go run cmd/platform/main.go db-seed
../bin/realize start
```

## Packaging

```bash
./scripts/build-docker.sh
```
