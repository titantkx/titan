import {
  GeneratedType,
  OfflineSigner,
  Registry,
} from "@titan-cosmjs/proto-signing";
import {
  SigningStargateClient,
  SigningStargateClientOptions,
  defaultRegistryTypes,
} from "@titan-cosmjs/stargate";
import { CometClient, connectComet } from "@titan-cosmjs/tendermint-rpc";

export const titanDefaultRegistryTypes: ReadonlyArray<[string, GeneratedType]> =
  [...defaultRegistryTypes];

function createDefaultRegistry(): Registry {
  return new Registry(titanDefaultRegistryTypes);
}

export class TitanSigningStargateClient extends SigningStargateClient {
  public static async connectWithSigner(
    endpoint: string,
    signer: OfflineSigner,
    options: SigningStargateClientOptions = {}
  ): Promise<TitanSigningStargateClient> {
    const cmClient = await connectComet(endpoint);
    return new TitanSigningStargateClient(cmClient, signer, {
      registry: createDefaultRegistry(),
      ...options,
    });
  }

  protected constructor(
    cmClient: CometClient | undefined,
    signer: OfflineSigner,
    options: SigningStargateClientOptions
  ) {
    super(cmClient, signer, options);
  }
}
