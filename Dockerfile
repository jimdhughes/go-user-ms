FROM scratch
WORKDIR /
ADD main /
VOLUME /data
EXPOSE 8080
ENV USERMS_SERVER_TYPE=http
ENV USERMS_GRPC_PORT=:50051
ENV USERMS_HTTP_PORT=:8080
CMD ["/main"]
