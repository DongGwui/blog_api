FROM postgres:16-alpine

# Install build dependencies
RUN apk add --no-cache --virtual .build-deps \
    git \
    make \
    gcc \
    musl-dev \
    clang15 \
    llvm15

# Clone and build pg_bigm
RUN git clone https://github.com/pgbigm/pg_bigm.git /tmp/pg_bigm \
    && cd /tmp/pg_bigm \
    && make USE_PGXS=1 \
    && make USE_PGXS=1 install

# Cleanup build dependencies
RUN apk del .build-deps \
    && rm -rf /tmp/pg_bigm

# Add pg_bigm to shared_preload_libraries
RUN echo "shared_preload_libraries = 'pg_bigm'" >> /usr/local/share/postgresql/postgresql.conf.sample