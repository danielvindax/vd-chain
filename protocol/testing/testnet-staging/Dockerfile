FROM vindax-base

COPY ./testing/testnet-staging/staging.sh /vindax/staging.sh
COPY ./testing/genesis.sh /vindax/genesis.sh
COPY ./testing/start.sh /vindax/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /vindax/exchange_config
COPY ./testing/delaymsg_config /vindax/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /vindax/staging.sh

ENV DAEMON_NAME=vindaxd
ENV DAEMON_HOME=/vindax/chain/.alice
ENV UNSAFE_SKIP_BACKUP=true
ENTRYPOINT ["cosmovisor", "run"]
