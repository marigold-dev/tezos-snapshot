export type SnapshotType = 'ROLLING' | 'FULL'
export type NetworkType = 'MAINNET' | 'TESTNET'
export type NetworkProtocolType = 'MAINNET' | 'HANGZHOUNET' | 'ITHACANET'

export type Snapshot = {
  file_name: string
  network_type: NetworkType
  chain: NetworkProtocolType
  date: Date
  snapshot_type: SnapshotType
  block_hash: string
  sha256: string
  block_height: string
  url: string
  filesize_bytes: number
}
