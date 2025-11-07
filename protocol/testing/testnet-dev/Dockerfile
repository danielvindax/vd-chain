FROM vindax-base

COPY ./testing/testnet-dev/dev.sh /vindax/dev.sh
COPY ./testing/genesis.sh /vindax/genesis.sh
COPY ./testing/start.sh /vindax/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /vindax/exchange_config
COPY ./testing/delaymsg_config /vindax/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /vindax/dev.sh

ENV DAEMON_NAME=vindaxd
ENTRYPOINT ["cosmovisor", "run"]
