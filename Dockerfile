FROM alpine

RUN apk --update add ca-certificates

COPY project-files/* /
COPY ./project-creator /project-creator

EXPOSE 8000
ENTRYPOINT ["/project-creator"]
