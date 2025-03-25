FROM ubuntu:18.04

ENV OPERATOR=/usr/local/bin/opencsd-api-server \
    USER_UID=root \
    USER_NAME=opencsd-api-server

COPY bin/opencsd-api-server ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup
RUN mkdir /root/.kube
RUN mkdir /mnt/instance-deploy

RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
