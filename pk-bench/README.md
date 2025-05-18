# pk-bench

## setup environment
```bash
docker compose up -d
```

## run application
```bash
❯ go run main.go

=== MySQL Insert Benchmark ===
Inserted 500000 records in 12.689810541s

=== MySQL Query Benchmark ===
UUIDv4 query took: 708.283666ms
UUIDv7 query took: 234.022916ms

=== PostgreSQL Insert Benchmark ===
Inserted 500000 records in 4.386692542s

=== PostgreSQL Query Benchmark ===
UUIDv4 query took: 75.465458ms
UUIDv7 query took: 32.69025ms
```

## check data
```bash
docker compose exec postgres /bin/bash
```

```bash
psql -U postgres -d benchmark_db
```

## check statistics

### select
```sql
EXPLAIN ANALYZE SELECT * from users;
```

```
benchmark_db=# EXPLAIN ANALYZE SELECT * from users ;
-[ RECORD 1 ]-------------------------------------------------------------------------------------------------------------
QUERY PLAN | Seq Scan on users  (cost=0.00..10155.00 rows=500000 width=51) (actual time=0.074..84.684 rows=500000 loops=1)
-[ RECORD 2 ]-------------------------------------------------------------------------------------------------------------
QUERY PLAN | Planning Time: 0.824 ms
-[ RECORD 3 ]-------------------------------------------------------------------------------------------------------------
QUERY PLAN | Execution Time: 104.725 ms
```

### insert
```sql
BEGIN;
EXPLAIN ANALYZE INSERT INTO users (id_v4, id_v7, name, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', '018e3c2e-7b6a-7cc2-bb2c-2b7b7b7b7b7b', 'Alice', '2024-07-01 12:34:56');
ROLLBACK;
```

```
benchmark_db=*# EXPLAIN ANALYZE INSERT INTO users (id_v4, id_v7, name, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', '018e3c2e-7b6a-7cc2-bb2c-2b7b7b7b7b7b', 'Alice', '2024-07-01 12:34:56');
-[ RECORD 1 ]-------------------------------------------------------------------------------------------
QUERY PLAN | Insert on users  (cost=0.00..0.01 rows=0 width=0) (actual time=1.141..1.141 rows=0 loops=1)
-[ RECORD 2 ]-------------------------------------------------------------------------------------------
QUERY PLAN |   ->  Result  (cost=0.00..0.01 rows=1 width=556) (actual time=0.023..0.023 rows=1 loops=1)
-[ RECORD 3 ]-------------------------------------------------------------------------------------------
QUERY PLAN | Planning Time: 0.534 ms
-[ RECORD 4 ]-------------------------------------------------------------------------------------------
QUERY PLAN | Execution Time: 1.381 ms
```

## check VACUUM info
```sql
SELECT * FROM pg_stat_all_tables WHERE relname = 'users';
```

```
benchmark_db=# SELECT * FROM pg_stat_all_tables WHERE relname = 'users';
-[ RECORD 1 ]-------+------------------------------
relid               | 16389
schemaname          | public
relname             | users
seq_scan            | 4
seq_tup_read        | 1500000
idx_scan            | 0
idx_tup_fetch       | 0
n_tup_ins           | 500001
n_tup_upd           | 0
n_tup_del           | 0
n_tup_hot_upd       | 0
n_live_tup          | 500000
n_dead_tup          | 1
n_mod_since_analyze | 0
n_ins_since_vacuum  | 1
last_vacuum         | 
last_autovacuum     | 2025-05-18 04:28:45.809697+00
last_analyze        | 
last_autoanalyze    | 2025-05-18 04:28:46.080775+00
vacuum_count        | 0
autovacuum_count    | 1
analyze_count       | 0
autoanalyze_count   | 1
```

```sql
SELECT
  relname,
  n_live_tup,
  n_dead_tup,
  CASE n_dead_tup WHEN 0 THEN 0 ELSE round(n_dead_tup*100/(n_live_tup+n_dead_tup) ,2) END AS ratio
FROM
  pg_stat_user_tables;
```

```
-[ RECORD 1 ]------
relname    | users
n_live_tup | 500000
n_dead_tup | 1
ratio      | 0.00
```

## cleanup
```bash
docker compose down -v
```