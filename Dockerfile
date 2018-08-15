FROM golang:1.10.0-stretch
MAINTAINER JimSung "hu0620711@126.com"

ENV HOST localhost
ENV PORT 50002
ENV DIRTY_WORDS ./dirty.txt

ENTRYPOINT ["/go/bin/Dirtyfilter"]
