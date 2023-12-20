import { GasPrice } from "@titan-cosmjs/stargate";
import { Field, Form, Formik } from "formik";
import { Button, FormGroup } from "react-bootstrap";
import * as Yup from "yup";
import { TitanSigningStargateClient } from "../titan_signingstargateclient";
import { parseCoins, validateCoins, validateGasPrice } from "../utils/helper";

interface SendProps {
  client: TitanSigningStargateClient;
}

const Send = ({ client }: SendProps) => {
  interface SendInputs {
    sender: string;
    receiver: string;
    amount: string;
    gas: string;
    gasPrice: string;
    memo?: string;
  }

  const initialValues: SendInputs = {
    sender: "titan1xnpqgn3hz2xvg054w4cc5d5r9dh00nvj9cmkte",
    receiver: "titan1zhxfglgt5ch2sls6gx346vvy4kn7w26jxy3xju",
    amount: "1tkx",
    gas: "auto",
    gasPrice: "10utkx",
  };

  const sendSchema = Yup.object().shape({
    sender: Yup.string().required(),
    receiver: Yup.string().required(),
    amount: Yup.string()
      .required()
      .test("validate-amount", "Invalid amount", validateCoins),
    gas: Yup.string()
      .required()
      .matches(/^(auto|\d+)$/, "Gas must be auto or number"),
    gasPrice: Yup.string()
      .required()
      .test("validate-price", "Invalid gas price", validateGasPrice),
    memo: Yup.string(),
  });

  const send = async ({
    sender,
    receiver,
    amount,
    gas,
    gasPrice,
    memo,
  }: SendInputs) => {
    try {
      const resp = await client.sendTokens(
        sender,
        receiver,
        parseCoins(amount),
        {
          gas: gas === "auto" ? "auto" : Number(gas),
          gasPrice: GasPrice.fromString(gasPrice),
        },
        memo
      );
      if (resp.code === 0) window.alert("Sent successfully");
      else window.alert(JSON.stringify(resp));
    } catch (e) {
      window.alert(e);
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={sendSchema}
      onSubmit={send}
    >
      {({ errors, touched, isValid }) => (
        <Form>
          <FormGroup>
            <Field name="sender" placeholder="Sender" />
            {errors.sender && touched.sender ? (
              <div>{errors.sender}</div>
            ) : null}
          </FormGroup>
          <FormGroup>
            <Field name="receiver" placeholder="Receiver" />
            {errors.receiver && touched.receiver ? (
              <div>{errors.receiver}</div>
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
            Send
          </Button>
        </Form>
      )}
    </Formik>
  );
};

export default Send;
