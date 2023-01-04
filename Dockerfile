FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/DanielPickens/Keeper

ENV USER 
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET CtqA-2OvKhcSDuvt

# 
ENV DSN $xxx_MYSQL_DSN

WORKDIR /go/src/github.com/DanielPickens/Keeper

RUN godep go build

EXPOSE 8888
CMD ./Keeper
