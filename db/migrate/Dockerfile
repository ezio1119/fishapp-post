FROM migrate/migrate:latest
COPY ./sql sql/
ENTRYPOINT ["sh", "-c", "/migrate $@", ""]