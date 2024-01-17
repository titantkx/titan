import { GasPrice } from "@titan-cosmjs/stargate";
import { Field, Form, Formik } from "formik";
import { Button, FormGroup } from "react-bootstrap";
import * as Yup from "yup";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import { validateGasPrice } from "../utils/helper";

interface WithdrawRewardsProps {
  client: TitanSigningStargateClient;
}

const WithdrawRewards = ({ client }: WithdrawRewardsProps) => {
  interface WithdrawRewardsInputs {
    delegator: string;
    validator: string;
    gas: string;
    gasPrice: string;
    memo?: string;
  }

  const initialValues: WithdrawRewardsInputs = {
    delegator: "titan16e6pnctgxcnv8y9n27p285gdnmgyl6ndsuu2nr",
    validator: "titanvaloper1cxjpv02d4cg7jp9qvh2her2lz5ljut0ulc3dua",
    gas: "auto",
    gasPrice:  `${10 * 1e10}atkx`,
  };

  const withdrawRewardsSchema = Yup.object().shape({
    delegator: Yup.string().required(),
    validator: Yup.string().required(),
    gas: Yup.string()
      .required()
      .matches(/^(auto|\d+)$/, "Gas must be auto or number"),
    gasPrice: Yup.string()
      .required()
      .test("validate-price", "Invalid gas price", validateGasPrice),
    memo: Yup.string(),
  });

  const withdrawRewards = async ({
    delegator,
    validator,
    gas,
    gasPrice,
    memo,
  }: WithdrawRewardsInputs) => {
    try {
      console.log("stake");
      const resp = await client.withdrawRewards(
        delegator,
        validator,
        {
          gas: gas === "auto" ? "auto" : Number(gas),
          gasPrice: GasPrice.fromString(gasPrice),
        },
        memo
      );
      if (resp.code === 0) window.alert("Withdrew rewards successfully");
      else window.alert(JSON.stringify(resp));
    } catch (e) {
      window.alert(e);
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={withdrawRewardsSchema}
      onSubmit={withdrawRewards}
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
            Withdraw Rewards
          </Button>
        </Form>
      )}
    </Formik>
  );
};

export default WithdrawRewards;
