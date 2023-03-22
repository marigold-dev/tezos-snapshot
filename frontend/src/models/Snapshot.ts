export type SnapshotType = 'ROLLING' | 'FULL'
export type NetworkType = 'MAINNET' | 'TESTNET'
export type NetworkProtocolType = 'MAINNET' | 'HANGZHOUNET' | 'ITHACANET'

export type Snapshot = {
  file_name: string
  network_type: NetworkType
  chain_name: NetworkProtocolType
  date: Date
  snapshot_type: SnapshotType
  block_hash: string
  sha256: string
  block_height: number
  url: string
  filesize_bytes: number
}
