FROM nginx:latest

# Preparing container
RUN apt-get update \
    && apt-get -y install curl git \
    && apt-get -y upgrade openssl

# Installing kubectl command-line tool
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

# Installing helm command-line tool
RUN curl -LO https://get.helm.sh/helm-v3.0.2-linux-amd64.tar.gz
RUN tar -xzvf helm-v3.0.2-linux-amd64.tar.gz
RUN cd linux-amd64 && chmod +x helm && mv ./helm /usr/local/bin/helm

# Copying sources to container
RUN mkdir /jaeger
WORKDIR /jaeger
COPY src/* .
RUN chmod +x create.sh
RUN git -c http.sslVerify=false clone https://github.wdf.sap.corp/EngSrv-Monitoring/helm-chart-jaeger.git
RUN cd helm-chart-jaeger && helm dependency update

# Command to start container with
CMD ./create.sh