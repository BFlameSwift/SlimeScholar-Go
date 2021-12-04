FROM ubuntu
RUN mkdir /app
COPY ./main /app/
RUN chmod +x /app/main
ENTRYPOINT [ "/app/main" ]