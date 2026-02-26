FROM golang:1.26 AS build

RUN apt install openssl && mkdir /app && \
    cd /app && openssl genrsa -out private-key.pem 2048 && \
    openssl rsa -in private-key.pem -pubout -out public-key.pem


COPY ./gateway /app/gateway
RUN cd /app/gateway && go build cmd/main.go

COPY ./auth_service /app/auth_service
RUN cd /app/auth_service && go build cmd/main.go

COPY ./websocket_manager /app/websocket_manager
RUN cd /app/websocket_manager && go build cmd/main.go

FROM ubuntu:24.04 AS auth-final
COPY --from=build /app/auth_service/main /app/private-key.pem /app/public-key.pem /app/auth_service/
COPY ./auth_service/config /app/auth_service/config
WORKDIR /app/auth_service
CMD ["./main"]

FROM ubuntu:24.04 AS websocket-final
COPY --from=build /app/websocket_manager/main /app/public-key.pem /app/websocket_manager/
COPY ./websocket_manager/config /app/websocket_manager/config
WORKDIR /app/websocket_manager
CMD ["./main"]

FROM ubuntu:24.04 AS gateway-final
COPY --from=build /app/gateway/main /app/gateway/
WORKDIR /app/gateway
CMD ["./main"]