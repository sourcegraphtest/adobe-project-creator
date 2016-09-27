FROM alpine

RUN apk --update add ca-certificates

COPY project-files/* /
COPY ./adobe-project-creator /adobe-project-creator

EXPOSE 8000
ENTRYPOINT ["/adobe-project-creator"]
