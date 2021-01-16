FROM ubuntu:18.04 as base

ENV DEBIAN_FRONTEND noninteractive

RUN apt update -y && \
    apt install -y wget build-essential --no-install-recommends

RUN apt install software-properties-common -y

RUN add-apt-repository ppa:deadsnakes/ppa

RUN add-apt-repository universe

RUN apt update -y

RUN apt install python3.8 -y

RUN apt install python3-pip -y

RUN apt install build-essential python3-dev python3-pip python3-setuptools python3-wheel python3-cffi libcairo2 \
    libpango-1.0-0 libpangocairo-1.0-0 libgdk-pixbuf2.0-0 libffi-dev shared-mime-info -y


RUN pip3 install cairosvg pillow lottie

FROM golang AS builder

COPY . /usr/src

WORKDIR /usr/src

RUN go build -o /usr/app/main ./cmd/main.go

FROM base

WORKDIR /usr/app

COPY --from=builder /usr/app .

COPY convert.py convert.py

RUN mkdir temp-images

EXPOSE 8080

CMD ["./main"]


