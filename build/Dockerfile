FROM alpine:latest

RUN apk update && apk upgrade && apk add bash curl
#RUN apt-get update \
#    && apt-get -y install curl

# Installing kubectl command-line tool
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

# Installing helm command-line tool
#RUN curl -LO https://get.helm.sh/helm-v3.0.2-linux-amd64.tar.gz
#RUN tar -xzvf helm-v3.0.2-linux-amd64.tar.gz
#RUN rm -f *.tar.gz
#RUN cd linux-amd64 && chmod +x helm && mv ./helm /usr/local/bin/helm

# Copying sources to container
RUN mkdir /jaeger
WORKDIR /jaeger
COPY output/*  /jaeger/
#RUN chmod +x create.sh
#RUN helm repo add jaegertracing https://jaegertracing.github.io/helm-charts

EXPOSE 8080

# Command to start container with
CMD ./deployer