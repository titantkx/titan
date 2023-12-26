import { GasPrice } from "@titan-cosmjs/stargate";
import { Field, Form, Formik } from "formik";
import { Button, FormGroup } from "react-bootstrap";
import * as Yup from "yup";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import { parseCoin, validateCoin, validateGasPrice } from "../utils/helper";

interface UnstakeProps {
  client: TitanSigningStargateClient;
}

const Unstake = ({ client }: UnstakeProps) => {
  interface UnstakeInputs {
    delegator: string;
    validator: string;
    amount: string;
    gas: string;
    gasPrice: string;
    memo?: string;
  }

  const initialValues: UnstakeInputs = {
    delegator: "titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr",
    validator: "titanvaloper1rn4kdmg3qm6ga279kvlqa2n8mf4rheex7medn9",
    amount: "1tkx",
    gas: "auto",
    gasPrice: "10atkx",
  };

  const unstakeSchema = Yup.object().shape({
    delegator: Yup.string().required(),
    validator: Yup.string().required(),
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

  const unstake = async ({
    delegator,
    validator,
    amount,
    gas,
    gasPrice,
    memo,
  }: UnstakeInputs) => {
    try {
      console.log("stake");
      const resp = await client.undelegateTokens(
        delegator,
        validator,
        parseCoin(amount),
        {
          gas: gas === "auto" ? "auto" : Number(gas),
          gasPrice: GasPrice.fromString(gasPrice),
        },
        memo
      );
      if (resp.code === 0) window.alert("Unstaked successfully");
      else window.alert(JSON.stringify(resp));
    } catch (e) {
      window.alert(e);
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={unstakeSchema}
      onSubmit={unstake}
    >
      {({ errors, touched, isValid }) => (
        <Form>
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
            Unstake
          </Button>
        </Form>
      )}
    </Formik>
  );
};

export default Unstake;
