import { Coin } from "@cosmjs/stargate";
import { Decimal } from "decimal.js/decimal";

export function validateCoin(amount: string): boolean {
  return /^\d+\.?\d*[a-zA-Z]+$/.test(amount);
}

export function validateCoins(amount: string): boolean {
  return /^\d+\.?\d*[a-zA-Z]+(,\d+\.?\d*[a-zA-Z]+)*$/.test(amount);
}

export function validateGasPrice(price: string): boolean {
  return /^[0-9.]+[a-z][a-z0-9]*$/i.test(price);
}

export function normalizeCoin(coin: Coin): Coin {
  if (coin.denom === "tkx") {
    const amount = new Decimal(coin.amount).mul(new Decimal("1e+18"));
    return { denom: "utkx", amount: amount.toString() };
  }
  return coin;
}

export function parseCoin(amount: string): Coin {
  const match = amount.match(/^(\d+\.?\d*)([a-zA-Z]+)$/);
  if (!match) throw new Error("Got an invalid coin string");
  return normalizeCoin({
    amount: match[1].replace(/^0+(?![1-9.])/, "") || "0",
    denom: match[2],
  });
}

export function parseCoins(amount: string): Coin[] {
  return amount.split(",").map((part) => parseCoin(part));
}
