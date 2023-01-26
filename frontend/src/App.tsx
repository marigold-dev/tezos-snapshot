import React, { CSSProperties, useEffect, useState } from 'react'
import axios from 'axios'
import { useTheme, ThemeProvider } from '@mui/material/styles'
import { Snapshot } from './models/Snapshot'
import HeaderBar from './components/HeaderBar'
import Content from './components/Content'

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
  }

  return (
    <ThemeProvider theme={theme}>
      <HeaderBar></HeaderBar>

      <div style={styles.page}>
        <Content />
      </div>
    </ThemeProvider>
  )
}

const hasNotPreviousDate = (
  index: number,
  array: Snapshot[],
  snapshot: Snapshot
) => (index > 0 && !(array[index - 1].Date === snapshot.Date)) || index === 0
