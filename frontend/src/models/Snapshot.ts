export type ArtifactType = 'tezos-snapshot' | 'tarball'
export type HistoryModeType = 'rolling' | 'full' | 'archive'
export type NetworkProtocolType = 'mainnet' | 'ithacanet' | string

export type Snapshot = {
  file_name: string
  chain_name: NetworkProtocolType
  date: Date
  artifact_type: ArtifactType
  history_mode: HistoryModeType
  block_hash: string
  sha256: string
  block_height: number
  gs_uri: string
  url: string
  filesize: string
  filesize_bytes: number
  tezos_version: {
    version: {
      additional_info: string
      major: number
      minor: number
    }
    implementation: string
    commit_info: {
      commit_date: string
      commit_hash: string
    }
  }
}
