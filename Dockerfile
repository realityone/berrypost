FROM node:latest as statics_dist
ADD . /data/berrypost
WORKDIR /data/berrypost/statics
RUN npm install && npm run build

FROM golang:latest as berrypost
ADD . /data/berrypost
COPY --from=statics_dist /data/berrypost/statics/dist/* /data/berrypost/statics/dist/
WORKDIR /data/berrypost
RUN CGO_ENABLED=0 GOOS=linux go build -o berrypost/berrypost berrypost/main.go

FROM alpine:latest
COPY --from=berrypost /data/berrypost/berrypost/berrypost /usr/local/bin/berrypost
CMD /usr/local/bin/berrypost