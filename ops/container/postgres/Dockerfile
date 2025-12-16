FROM postgres:16-bullseye

# Install dependencies for PostgreSQL Anonymizer
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
       pgxnclient \
       postgresql-server-dev-16 \
       gcc \
       make \
  && rm -rf /var/lib/apt/lists/*

# Install Anonymizer
RUN pgxn install postgresql_anonymizer=1.1.0

RUN mkdir -p /docker-entrypoint-initdb.d
