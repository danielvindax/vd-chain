FROM vindax-base

RUN apk add --no-cache bash jq aws-cli
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

COPY ./testing/version/. /vindax/
COPY ./testing/mainnet/. /vindax/

ENV HOME /vindax
WORKDIR $HOME

RUN /vindax/mainnet.sh

ENTRYPOINT ["/vindax/start.sh"]
