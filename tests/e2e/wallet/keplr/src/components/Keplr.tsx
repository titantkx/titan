import { ChainInfo, Window as KeplrWindow } from "@keplr-wallet/types";
import { useEffect, useState } from "react";
import { Col, Container, Row } from "react-bootstrap";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import Send from "./Send";
import Stake from "./Stake";
import Unstake from "./Unstake";
import WithdrawRewards from "./WithdrawRewards";

declare global {
  interface Window extends KeplrWindow {}
}

const chainInfo: ChainInfo = {
  chainId: process.env.REACT_APP_CHAIN_ID!,
  chainName: process.env.REACT_APP_CHAIN_NAME!,
  rpc: process.env.REACT_APP_RPC_URL!,
  rest: process.env.REACT_APP_REST_URL!,
  bip44: {
    coinType: 60,
  },
  bech32Config: {
    bech32PrefixAccAddr: "titan",
    bech32PrefixAccPub: "titanpub",
    bech32PrefixValAddr: "titanvaloper",
    bech32PrefixValPub: "titanvaloperpub",
    bech32PrefixConsAddr: "titanvalcons",
    bech32PrefixConsPub: "titanvalconspub",
  },
  currencies: [
    {
      coinDenom: "tkx",
      coinMinimalDenom: "atkx",
      coinDecimals: 18,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: "tkx",
      coinMinimalDenom: "atkx",
      coinDecimals: 18,
      gasPriceStep: {
        low: 0.001 * 1e10,
        average: 0.025 * 1e10,
        high: 0.04 * 1e10,
      },
    },
  ],
  stakeCurrency: {
    coinDenom: "tkx",
    coinMinimalDenom: "atkx",
    coinDecimals: 18,
  },
  features: ["eth-address-gen", "eth-key-sign"],
};

const KeplrView = () => {
  const [client, setClient] = useState<TitanSigningStargateClient>();

  useEffect(() => {
    const addTitanToKeplr = async () => {
      const { keplr } = window;
      if (!keplr) {
        alert("You need to install Keplr");
        throw new Error("You need to install Keplr");
      }
      keplr.defaultOptions = {
        sign: {
          preferNoSetFee: true,
        },
      };
      await keplr.experimentalSuggestChain(chainInfo);
      await keplr.enable(process.env.REACT_APP_CHAIN_ID!);
      const client = await TitanSigningStargateClient.connectWithSigner(
        process.env.REACT_APP_RPC_URL!,
        keplr.getOfflineSigner(process.env.REACT_APP_CHAIN_ID!),
        { isEthermint: true }
      );
      setClient(client);
    };
    addTitanToKeplr();
  }, []);

  return (
    <Container fluid className="p-5">
      {client && (
        <Row>
          <Col>
            <Send client={client} />
          </Col>
          <Col>
            <Stake client={client} />
          </Col>
          <Col>
            <Unstake client={client} />
          </Col>
          <Col>
            <WithdrawRewards client={client} />
          </Col>
        </Row>
      )}
    </Container>
  );
};

export default KeplrView;
