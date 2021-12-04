FROM alpine
RUN mkdir /app
COPY ./main /app/main
ENTRYPOINT ["/app/main"]