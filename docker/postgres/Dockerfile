FROM postgres:16-bookworm

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    make \
    gcc \
    postgresql-server-dev-16 \
    && rm -rf /var/lib/apt/lists/*

# Clone and build pg_bigm
RUN git clone https://github.com/pgbigm/pg_bigm.git /tmp/pg_bigm \
    && cd /tmp/pg_bigm \
    && make USE_PGXS=1 \
    && make USE_PGXS=1 install \
    && rm -rf /tmp/pg_bigm

# Cleanup build dependencies (optional, keeps image smaller)
RUN apt-get purge -y --auto-remove git make gcc postgresql-server-dev-16

# Add pg_bigm to shared_preload_libraries
RUN echo "shared_preload_libraries = 'pg_bigm'" >> /usr/share/postgresql/postgresql.conf.sample
