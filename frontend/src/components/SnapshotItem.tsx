
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
        <Typography style={{ fontWeight: 'bold', fontSize: 'x-small', overflowWrap: 'break-word', color: 'rgb(235, 52, 72)' }} component="div">Block Hash:</Typography>
        <Typography style={{ fontWeight: 'bold', fontSize: 'larger', overflowWrap: 'break-word' }} component="div">
           {props.snapshot.Blockhash}
        </Typography>
        <Typography style={{ fontWeight: 'bold', fontSize: 'x-small', overflowWrap: 'break-word', color: 'rgb(235, 52, 72)' }} component="div">SHA256 Checksum:</Typography>
        <Typography style={{ fontWeight: 'bold', fontSize: '13px', overflowWrap: 'break-word' }} component="div">
           {props.snapshot.SHA256Checksum}
        </Typography>
      </CardContent>

      <CardActions sx={{ justifyContent: 'right' }}>
        <Button sx={{ color: theme.palette.text.primary, textDecoration: 'underline' }} size="small"
         href={'https://' + props.snapshot.NetworkProtocol + '.tzkt.io/' + props.snapshot.Blockhash}>
            TzKT
        </Button>
        <Button sx={{ color: theme.palette.text.primary, textDecoration: 'underline' }} size="small"
         href={getTzStatsLink(props.snapshot.NetworkProtocol) + props.snapshot.Blockhash}>
            TzStats
        </Button>
        <Button sx={{ color: theme.palette.text.primary, textDecoration: 'underline' }} size="small"
         href={props.snapshot.PublicURL}>
            Download
        </Button>
      </CardActions>
    </Card>
  )
}

const getTzStatsLink = (networkProtocol: string) => {
  if (networkProtocol === 'MAINNET') {
    return ('https://tzstats.com/')
  }
  const network = networkProtocol.slice(0, -3)
  return ('https://' + network + '.tzstats.com/')
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
