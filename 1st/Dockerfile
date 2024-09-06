ARG GO_IMAGE

FROM ${GO_IMAGE}

ARG EXECUTABLE
ARG EXECUTABLE_PATH

ENV EXECUTABLE=${EXECUTABLE}
ENV EXECUTABLE_PATH=${EXECUTABLE_PATH}

# COPY ./scripts/healthcheck.sh /usr/local/bin/healthcheck.sh
# RUN chmod +x /usr/local/bin/healthcheck.sh

COPY ./${EXECUTABLE} ${EXECUTABLE_PATH}/${EXECUTABLE}
RUN chmod +x ${EXECUTABLE_PATH}/${EXECUTABLE}

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]