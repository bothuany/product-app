# Start the PostgreSQL container
docker run --name postgres-test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 6432:5432 -d postgres:latest

Write-Host "Postgresql starting..."
Start-Sleep -Seconds 3

# Create the database
docker exec -it postgres-test psql -U postgres -d postgres -c "CREATE DATABASE productapp;"
Start-Sleep -Seconds 3
Write-Host "Database productapp created"

# Create the table
docker exec -it postgres-test psql -U postgres -d productapp -c "
create table if not exists products
(
  id bigserial not null primary key,
  name varchar(255) not null,
  price double precision not null,
  discount double precision,
  store varchar(255) not null
);
"

Start-Sleep -Seconds 3
Write-Host "Table products created"