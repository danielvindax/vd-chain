FROM vindax-base

COPY ./testing/containertest/vindaxd_preupgrade* /bin/
COPY ./testing/containertest/containertest.sh /vindax/containertest.sh
COPY ./testing/containertest/preupgrade_entrypoint.sh /vindax/preupgrade_entrypoint.sh
COPY ./testing/containertest/preupgrade_genesis.json /vindax/preupgrade_genesis.json
COPY ./testing/genesis.sh /vindax/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /vindax/exchange_config
COPY ./testing/delaymsg_config /vindax/delaymsg_config
COPY ./testing/version/. /vindax/

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /vindax/containertest.sh

COPY ./testing/containertest/config/pricefeed_exchange_config.toml /vindax/chain/.alice/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /vindax/chain/.bob/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /vindax/chain/.carl/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /vindax/chain/.dave/config/pricefeed_exchange_config.toml

ENTRYPOINT ["vindaxd"]
