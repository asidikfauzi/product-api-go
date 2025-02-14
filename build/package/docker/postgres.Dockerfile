FROM postgres:latest

# Copy init script ke dalam direktori inisialisasi Docker PostgreSQL
COPY build/package/postgres/init-db.sh /docker-entrypoint-initdb.d/init-db.sh

# Beri izin eksekusi ke skrip
RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh