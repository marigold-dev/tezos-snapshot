import React, { CSSProperties, useEffect, useState } from 'react'
import axios from 'axios'
import { useTheme, ThemeProvider } from '@mui/material/styles'
import { Snapshot } from './models/Snapshot'
import SnapshotItem from './components/SnapshotItem'
import HeaderBar from './components/HeaderBar'
import DateSeparator from './components/DateSeparator'

export function App() {
  const theme = useTheme()

  const [snapshots, setSnapshots] = useState<Array<Snapshot>>([])
  useEffect(() => {
    if (snapshots.length < 1) {
      axios.get(process.env.REACT_APP_BACKEND_URL!).then((request: any) => {
        setSnapshots(request.data)
      })
    }
  }, [snapshots, setSnapshots])

  const styles = {
    page: {
      backgroundColor: theme.palette.primary.main,
      overflow: 'auto',
      height: 'calc(100vh - 66px)',
    } as CSSProperties,
    content: {
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      fontSize: 'calc(2vmin)',
    } as CSSProperties,
  }

  return (
    <ThemeProvider theme={theme}>
      <HeaderBar></HeaderBar>

      <div style={styles.page}>
        <div style={styles.content}>
          {snapshots.map((snapshot, index, array) => (
            <div
              key={snapshot.PublicURL}
              style={{ paddingBottom: '30px', textAlign: 'left' }}
            >
              {hasNotPreviousDate(index, array, snapshot) && (
                <DateSeparator snapshot={snapshot}></DateSeparator>
              )}

              <SnapshotItem snapshot={snapshot}></SnapshotItem>
            </div>
          ))}
        </div>
      </div>
    </ThemeProvider>
  )
}

const hasNotPreviousDate = (
  index: number,
  array: Snapshot[],
  snapshot: Snapshot
) => (index > 0 && !(array[index - 1].Date === snapshot.Date)) || index === 0
