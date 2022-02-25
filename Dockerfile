FROM golang:1.17 as gobuilder
WORKDIR /bootstrap
# install required commands
RUN go get github.com/GeertJohan/go.rice/rice
# cache
COPY bootstrap/go.mod bootstrap/go.sum ./
RUN go mod download
# build
COPY bootstrap ./
RUN go generate ./...
RUN CGO_ENABLED=0 go build

FROM alpine/git:latest as git
ARG VERSION=3.7.0
WORKDIR /checkout
RUN git clone https://github.com/hakimel/reveal.js.git . && \
    git checkout $VERSION && \
    rm -rf .git

FROM node:10-alpine 
COPY --from=git /checkout /reveal.js
WORKDIR /reveal.js
RUN npm install
RUN sed -i -e "s/open: true/open: false/" Gruntfile.js
COPY --from=gobuilder /bootstrap/bootstrap /bin/bootstrap
ENTRYPOINT [ "/bin/bootstrap" ]
CMD [ "start" ]