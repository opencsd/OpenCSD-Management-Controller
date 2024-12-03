FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/opencsd-controller \
    USER_UID=root \
    USER_NAME=opencsd-controller

COPY bin/opencsd-controller ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup
RUN mkdir /root/.kube
RUN mkdir /mnt/instance-deploy

RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
