FROM golang:1.22 as gobuilder
WORKDIR /bootstrap
# cache
COPY bootstrap/go.mod bootstrap/go.sum ./
RUN go mod download
# build
COPY bootstrap ./
RUN go generate ./...
RUN CGO_ENABLED=0 go build

FROM alpine/git:latest as git
ARG REVEALJS_VERSION=5.1.0
WORKDIR /checkout
RUN git clone https://github.com/hakimel/reveal.js.git . && \
    git checkout $REVEALJS_VERSION && \
    rm -rf .git

FROM node:22-alpine 
COPY --from=git /checkout /reveal.js
WORKDIR /reveal.js
RUN npm install
COPY --from=gobuilder /bootstrap/bootstrap /bin/bootstrap
ENTRYPOINT [ "/bin/bootstrap" ]
EXPOSE 8000
CMD [ "start" ]