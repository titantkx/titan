import { Ed25519Pubkey } from "@titan-cosmjs/amino";
import { Decimal } from "@titan-cosmjs/math";
import { encodePubkey } from "@titan-cosmjs/proto-signing";
import {
  GasPrice,
  MsgCreateValidatorEncodeObject,
} from "@titan-cosmjs/stargate";
import { Field, Form, Formik } from "formik";
import { Button, FormGroup } from "react-bootstrap";
import { MsgCreateValidator } from "titan-cosmjs-types/cosmos/staking/v1beta1/tx";
import * as Yup from "yup";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import { parseCoin, validateCoin, validateGasPrice } from "../utils/helper";

interface CreateValidatorForOtherProps {
  client: TitanSigningStargateClient;
}

const CreateValidatorForOther = ({ client }: CreateValidatorForOtherProps) => {
  interface StakeInputs {
    moniker: string;
    identity: string;
    website: string;
    securityContact: string;
    details: string;
    commissionRate: string;
    commissionMaxRate: string;
    commissionMaxChangeRate: string;
    minSelfDelegation: string;
    delegator: string;
    validator: string;
    amount: string;
    pubkey: string;
    gas: string;
    gasPrice: string;
    memo?: string;
  }

  const initialValues: StakeInputs = {
    moniker: "mokier",
    identity: "identity",
    website: "https://exmaple.com",
    securityContact: "email@example.com",
    details: "details",
    commissionRate: "0.1",
    commissionMaxRate: "0.2",
    commissionMaxChangeRate: "0.01",
    minSelfDelegation: "5000000000000000000",
    delegator: "titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr",
    validator: "titanvaloper16e6pnctgxcnv8y9n27p285gdnmgyl6ndwge0tj",
    pubkey: "dKL3iZHRviSB2FGVVjMDJHAEMODeJN9lsjY8sKR9guo=",
    amount: "10tkx",
    gas: "auto",
    gasPrice: `${10 * 1e10}atkx`,
  };

  const stakeSchema = Yup.object().shape({
    moniker: Yup.string().required(),
    identity: Yup.string().required(),
    website: Yup.string().required().url(),
    securityContact: Yup.string().required().email(),
    details: Yup.string().required(),
    commissionRate: Yup.string()
      .required()
      .matches(/^\d+(\.\d{0,18})?$/),
    commissionMaxRate: Yup.string()
      .required()
      .matches(/^\d+(\.\d{0,18})?$/),
    commissionMaxChangeRate: Yup.string()
      .required()
      .matches(/^\d+(\.\d{0,18})?$/),
    minSelfDelegation: Yup.string().matches(/^\d+$/),
    delegator: Yup.string().required(),
    validator: Yup.string().required(),
    pubkey: Yup.string().required(),
    amount: Yup.string()
      .required()
      .test("validate-amount", "Invalid amount", validateCoin),
    gas: Yup.string()
      .required()
      .matches(/^(auto|\d+)$/, "Gas must be auto or number"),
    gasPrice: Yup.string()
      .required()
      .test("validate-price", "Invalid gas price", validateGasPrice),
    memo: Yup.string(),
  });

  const stake = async ({
    moniker,
    identity,
    website,
    securityContact,
    details,
    commissionRate,
    commissionMaxRate,
    commissionMaxChangeRate,
    minSelfDelegation,
    delegator,
    validator,
    pubkey,
    amount,
    gas,
    gasPrice,
    memo,
  }: StakeInputs) => {
    try {
      const pk: Ed25519Pubkey = {
        type: "tendermint/PubKeyEd25519",
        value: pubkey,
      };
      const msg: MsgCreateValidatorEncodeObject = {
        typeUrl: "/cosmos.staking.v1beta1.MsgCreateValidator",
        value: MsgCreateValidator.fromPartial({
          description: {
            moniker: moniker,
            identity: identity,
            website: website,
            securityContact: securityContact,
            details: details,
          },
          commission: {
            rate: Decimal.fromUserInput(commissionRate, 18).atomics,
            maxRate: Decimal.fromUserInput(commissionMaxRate, 18).atomics,
            maxChangeRate: Decimal.fromUserInput(commissionMaxChangeRate, 18)
              .atomics,
          },
          minSelfDelegation: minSelfDelegation,
          delegatorAddress: delegator,
          validatorAddress: validator,
          pubkey: encodePubkey(pk),
          value: parseCoin(amount),
        }),
      };
      const resp = await client.signAndBroadcast(
        delegator,
        [msg],
        {
          gas: gas === "auto" ? "auto" : Number(gas),
          gasPrice: GasPrice.fromString(gasPrice),
        },
        memo
      );
      if (resp.code === 0) {
        window.alert("Created validator successfully");
        console.log(resp.transactionHash);
      } else window.alert(JSON.stringify(resp));
    } catch (e) {
      window.alert(e);
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={stakeSchema}
      onSubmit={stake}
    >
      {({ errors, touched, isValid }) => (
        <Form>
          <FormGroup>
            <Field name="moniker" placeholder="Moniker" />
            {errors.moniker && touched.moniker ? (
              <div>{errors.moniker}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="identity" placeholder="Identity" />
            {errors.identity && touched.identity ? (
              <div>{errors.identity}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="website" placeholder="Website" />
            {errors.website && touched.website ? (
              <div>{errors.website}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="securityContact" placeholder="Security Contact" />
            {errors.securityContact && touched.securityContact ? (
              <div>{errors.securityContact}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="details" placeholder="Details" />
            {errors.details && touched.details ? (
              <div>{errors.details}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="commissionRate" placeholder="Commission Rate" />
            {errors.commissionRate && touched.commissionRate ? (
              <div>{errors.commissionRate}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="commissionMaxRate" placeholder="Commission Max Rate" />
            {errors.commissionMaxRate && touched.commissionMaxRate ? (
              <div>{errors.commissionMaxRate}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field
              name="commissionMaxChangeRate"
              placeholder="Commission Max Change Rate"
            />
            {errors.commissionMaxChangeRate &&
            touched.commissionMaxChangeRate ? (
              <div>{errors.commissionMaxChangeRate}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="minSelfDelegation" placeholder="Min Self Delegation" />
            {errors.minSelfDelegation && touched.minSelfDelegation ? (
              <div>{errors.minSelfDelegation}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="delegator" placeholder="Delegator" />
            {errors.delegator && touched.delegator ? (
              <div>{errors.delegator}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="validator" placeholder="Validator" />
            {errors.validator && touched.validator ? (
              <div>{errors.validator}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="pubkey" placeholder="Validator Public Key" />
            {errors.pubkey && touched.pubkey ? (
              <div>{errors.pubkey}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="amount" placeholder="Amount" />
            {errors.amount && touched.amount ? (
              <div>{errors.amount}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="gas" placeholder="Gas" />
            {errors.gas && touched.gas ? <div>{errors.gas}</div> : null}
          </FormGroup>
          <FormGroup>
            <Field name="gasPrice" placeholder="Gas price" />
            {errors.gasPrice && touched.gasPrice ? (
              <div>{errors.gasPrice}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="memo" placeholder="Memo" />
            {errors.memo && touched.memo ? <div>{errors.memo}</div> : null}
          </FormGroup>
          <Button type="submit" disabled={!isValid}>
            Create Validator
          </Button>
        </Form>
      )}
    </Formik>
  );
};

export default CreateValidatorForOther;
