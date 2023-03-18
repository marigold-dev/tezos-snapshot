import React, { CSSProperties, useEffect, useState } from 'react'
import axios from 'axios'
import { useTheme } from '@mui/material/styles'
import { Snapshot } from '../models/Snapshot'
import DateSeparator from './DateSeparator'
import SnapshotItem from './SnapshotItem'
import ReactLoading from 'react-loading'

const hasNotPreviousDate = (
  index: number,
  array: Snapshot[],
  snapshot: Snapshot
) => (index > 0 && !(array[index - 1].date === snapshot.date)) || index === 0

export default function Content() {
  const theme = useTheme()

  const [loading, setLoading] = useState(true)
  const [snapshots, setSnapshots] = useState<Array<Snapshot>>([])
  useEffect(() => {
    if (snapshots.length < 1) {
      axios.get(process.env.REACT_APP_BACKEND_URL!).then((request: any) => {
        setSnapshots(request.data.data)
        setLoading(false)
      })
    }
  }, [snapshots, setSnapshots])

  const styles = {
    content: {
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      fontSize: 'calc(2vmin)',
    } as CSSProperties,
    loading: {
      justifyContent: 'center',
      display: 'flex',
      height: 'calc(100vh - 66px)',
      alignItems: 'center',
    },
  }

  if (loading)
    return (
      <div style={styles.loading}>
        <ReactLoading type="spin" color={theme.palette.secondary.main} />
      </div>
    )

  return (
    <div style={styles.content}>
      {snapshots.map((snapshot, index, array) => (
        <div
          key={snapshot.url}
          style={{ paddingBottom: '30px', textAlign: 'left' }}
        >
          {hasNotPreviousDate(index, array, snapshot) && (
            <DateSeparator snapshot={snapshot}></DateSeparator>
          )}

          <SnapshotItem snapshot={snapshot}></SnapshotItem>
        </div>
      ))}
    </div>
  )
}
