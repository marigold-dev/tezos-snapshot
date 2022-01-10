export type SnapshotType = 'ROLLING' | 'FULL'
export type NetworkType = 'MAINNET' | 'TESTNET'

export type Snapshot = {
  FileName: string;
  Network: NetworkType;
  NetworkProtocol: NetworkType;
  Date: Date;
  SnapshotType: SnapshotType;
  Blockhash: string;
  Blocklevel: string;
  PublicURL: string;
  Size: number;
}
