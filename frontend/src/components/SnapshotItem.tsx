
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
      minWidth: '80vh',
      backgroundColor: theme.palette.primary.main,
      color: theme.palette.text.primary,
      boxShadow: 'none',
      border: 'solid',
      borderWidth: '1px',
      borderRadius: '0px'
    }}>
      <CardContent>
        <Typography sx={{ fontSize: 14, display: 'flex' }} gutterBottom>
          <span style={{ flex: '1' }}>{props.snapshot.SnapshotType} - {props.snapshot.Network}</span>
        </Typography>
        <Typography style={{ fontWeight: 'bold', fontSize: '18px' }} component="div">
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

export default SnapshotItem
