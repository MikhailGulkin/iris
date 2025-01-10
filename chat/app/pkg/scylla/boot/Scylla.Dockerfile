FROM quay.io/debezium/connect:2.4

# Устанавливаем ScyllaDB CDC Source Connector
RUN curl -L https://github.com/scylladb/scylla-cdc-source-connector/releases/download/scylla-cdc-source-connector-1.0.1/ScyllaDB-scylla-cdc-source-connector-1.0.1.zip \
    -o /tmp/ScyllaDB-scylla-cdc-source-connector-1.0.1.zip && \
    unzip /tmp/ScyllaDB-scylla-cdc-source-connector-1.0.1.zip -d /kafka/connect && \
    rm /tmp/ScyllaDB-scylla-cdc-source-connector-1.0.1.zip

# Устанавливаем дополнительные плагины, если нужно
