FROM postgres:13.0

ARG POSTGRES_DB
COPY init.sql /docker-entrypoint-initdb.d/
RUN sed -i "s|POSTGRES_DB|${POSTGRES_DB}|g" /docker-entrypoint-initdb.d/init.sql
EXPOSE 5432
CMD ["postgres"]
