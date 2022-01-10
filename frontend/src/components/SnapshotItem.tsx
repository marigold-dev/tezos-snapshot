
import React from 'react'
import Card from '@mui/material/Card'
import CardActions from '@mui/material/CardActions'
import CardContent from '@mui/material/CardContent'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import { Snapshot } from '../models/Snapshot'
import { useTheme } from '@mui/material/styles'

const SnapshotItem = (props: { snapshot: Snapshot }) => {
  const theme = useTheme()
  return (
    <Card sx={{
      width: '80vw',
      maxWidth: '800px',
      backgroundColor: theme.palette.primary.main,
      color: theme.palette.text.primary,
      boxShadow: 'none',
      border: 'solid',
      borderWidth: '1px',
      borderRadius: '0px'
    }}>
      <CardContent>
        <Typography sx={{ fontSize: '14px', display: 'flex' }} gutterBottom>
          <span style={{ flex: '1' }}> {
          props.snapshot.NetworkProtocol === props.snapshot.Network
            ? props.snapshot.Network
            : `${props.snapshot.Network} - ${props.snapshot.NetworkProtocol}`
          } - {props.snapshot.SnapshotType} - {formatBytes(props.snapshot.Size ?? 0)}</span>
        </Typography>
        <Typography style={{ fontWeight: 'bold', fontSize: '2.3vh', overflowWrap: 'break-word' }} component="div">
          {props.snapshot.Blockhash}
        </Typography>
      </CardContent>

      <CardActions sx={{ justifyContent: 'right' }}>
        <Button sx={{ color: theme.palette.text.primary, textDecoration: 'underline' }} size="small"
          onClick={() => { window.location.href = props.snapshot.PublicURL }}>Download</Button>
      </CardActions>
    </Card>
  )
}

const formatBytes = (bytes: number, decimals: number = 2) => {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

export default SnapshotItem
