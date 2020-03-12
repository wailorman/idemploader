FROM ruby:2.5.1-alpine as builder

RUN mkdir /docs
COPY . /docs/
WORKDIR /docs

RUN bin/assemble.rb

FROM docker.io/swaggerapi/swagger-ui

RUN mkdir /docs
COPY --from=builder /docs/build/openapi.json /docs/openapi.json
ENV SWAGGER_JSON /docs/openapi.json
