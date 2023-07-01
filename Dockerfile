FROM golang:1.19
ARG version=435.0.0

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-$version-linux-x86_64.tar.gz
RUN mv google-cloud-cli-$version-linux-x86_64.tar.gz /root
RUN tar -xf /root/google-cloud-cli-$version-linux-x86_64.tar.gz -C /root
ENV PATH $PATH:/root/google-cloud-sdk/bin