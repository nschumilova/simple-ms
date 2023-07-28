FROM migrate/migrate:v4.16.2

COPY database/migration/postgres /migrations

ENV DB_HOST "localhost"
ENV DB_PORT 5432
ENV DB_NAME ""
ENV DB_USER ""
ENV DB_PASSWORD ""
ENV COMMAND "up"

ENTRYPOINT migrate -path=/migrations -database=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable ${COMMAND}