export type SnapshotType = 'ROLLING' | 'FULL'
export type NetworkType = 'MAINNET' | 'TESTNET'
export type NetworkProtocolType = 'MAINNET' | 'HANGZHOUNET' | 'ITHACANET'

export type Snapshot = {
  FileName: string;
  Network: NetworkType;
  NetworkProtocol: NetworkProtocolType;
  Date: Date;
  SnapshotType: SnapshotType;
  Blockhash: string;
  SHA256Checksum: string;
  Blocklevel: string;
  PublicURL: string;
  Size: number;
}
