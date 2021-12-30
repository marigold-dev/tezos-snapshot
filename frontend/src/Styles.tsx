
import { makeStyles } from '@material-ui/styles'
import { useTheme } from '@mui/material/styles'

const theme = useTheme()
const useStyles = makeStyles({
  page: {
    backgroundColor: theme.palette.primary.main,
    overflow: 'auto',
    height: 'calc(100vh - 66px)'
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    fontSize: 'calc(2vmin)'
  },
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
})

export default useStyles
