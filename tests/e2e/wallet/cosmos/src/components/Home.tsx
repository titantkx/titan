import { ChainInfo, Keplr, Window as KeplrWindow } from "@keplr-wallet/types";
import React, { useState } from "react";
import { Col, Container, Row } from "react-bootstrap";
import Dropdown from "react-bootstrap/Dropdown";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import CreateValidator from "./CreateValidator";
import CreateValidatorForOther from "./CreateValidatorForOther";
import Send from "./Send";
import Stake from "./Stake";
import StakeForOther from "./StakeForOther";
import Unstake from "./Unstake";
import WithdrawRewards from "./WithdrawRewards";

declare global {
  interface Window extends KeplrWindow {
    leap: any;
  }
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
      coinDenom: "TKX",
      coinMinimalDenom: "atkx",
      coinDecimals: 18,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: "TKX",
      coinMinimalDenom: "atkx",
      coinDecimals: 18,
      gasPriceStep: {
        low: 10 * 1e10,
        average: 11 * 1e10,
        high: 20 * 1e10,
      },
    },
  ],
  stakeCurrency: {
    coinDenom: "TKX",
    coinMinimalDenom: "atkx",
    coinDecimals: 18,
  },
  features: ["cosmwasm", "eth-address-gen", "eth-key-sign"],
};

const Home = () => {
  const [walletName, setWalletName] = useState<string>();
  const [client, setClient] = useState<TitanSigningStargateClient>();

  const addTitanToWallet = async (wallet: Keplr | undefined) => {
    if (!wallet) {
      alert("You need to install " + walletName);
      throw new Error("You need to install " + walletName);
    }
    wallet.defaultOptions = {
      sign: {
        preferNoSetFee: true,
      },
    };
    await wallet.experimentalSuggestChain(chainInfo);
    await wallet.enable(process.env.REACT_APP_CHAIN_ID!);
    const client = await TitanSigningStargateClient.connectWithSigner(
      process.env.REACT_APP_RPC_URL!,
      wallet.getOfflineSigner(process.env.REACT_APP_CHAIN_ID!),
      { isEthermint: true }
    );
    setClient(client);
  };

  const handleSelectWallet = (evtKey: any) => {
    if (evtKey === "keplr") {
      setWalletName("Keplr");
      addTitanToWallet(window.keplr);
    } else if (evtKey === "leap") {
      setWalletName("Leap");
      addTitanToWallet(window.leap);
    }
  };

  return (
    <Container fluid className="p-4">
      <Row>
        <Dropdown onSelect={handleSelectWallet}>
          <Dropdown.Toggle id="dropdown-wallet">
            {walletName || "Select Wallet"}
          </Dropdown.Toggle>
          <Dropdown.Menu>
            <Dropdown.Item eventKey="keplr">Keplr</Dropdown.Item>
            <Dropdown.Item eventKey="leap">Leap</Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown>
      </Row>
      {client && (
        <React.Fragment>
          <Row className="mt-4">
            <Col>
              <Send client={client} />
            </Col>
            <Col>
              <Stake client={client} />
            </Col>
            <Col>
              <StakeForOther client={client} />
            </Col>
            <Col>
              <Unstake client={client} />
            </Col>
            <Col>
              <WithdrawRewards client={client} />
            </Col>
          </Row>
          <Row className="mt-4">
            <Col>
              <CreateValidator client={client} />
            </Col>
            <Col>
              <CreateValidatorForOther client={client} />
            </Col>
          </Row>
        </React.Fragment>
      )}
    </Container>
  );
};

export default Home;
