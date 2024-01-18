# FROM registry.access.redhat.com/ubi7/ubi-minimal:latest
# FROM nvidia/cuda:11.0.3-base-ubuntu18.04
FROM ubuntu:18.04

ENV OPERATOR=/usr/local/bin/instance-metric-collector \
    USER_UID=root \
    USER_NAME=discovery
COPY _output/bin/instance-metric-collector ${OPERATOR}

RUN mkdir -p /mnt
RUN cd /mnt
RUN mkdir -p power
RUN mkdir -p cpu/stat

ENTRYPOINT ["/usr/local/bin/instance-metric-collector"]

USER ${USER_UID}
