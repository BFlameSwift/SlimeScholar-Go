FROM alpine
RUN mkdir /app
WORKDIR /app
COPY ./main .
CMD ./main