import React from 'react'
import Link from '@mui/material/Link'
import { useTheme } from '@mui/material/styles'

const SnapshotLink = (props: { url: string, children: string }) => {
  const theme = useTheme()
  const styles = {
    button: {
      color: theme.palette.text.primary,
      paddingLeft: '10px',
      height: '100%',
      display: 'flex',
      alignItems: 'center',
      '&:hover': {
        backgroundColor: '#eb3448 !important',
        cursor: 'pointer'
      }
    },
    link: {
      color: theme.palette.text.primary,
      paddingRight: '10px',
      fontFamily: '"Roboto","Helvetica","Arial",sans-serif'
    }
  }

  return (

    <span style={styles.button} onClick={() => { window.location.href = props.url }} >
      <Link sx={styles.link} href={props.url} underline="none">
        {props.children}
      </Link>
    </span>
  )
}

export default SnapshotLink
