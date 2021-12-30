import React, { useEffect, useMemo, useState, createContext } from 'react'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import { PaletteMode, PaletteOptions } from '@mui/material'
import useMediaQuery from '@mui/material/useMediaQuery'
import { App } from './App'

export const ColorModeContext = createContext({ toggleColorMode: () => { } })

export default function ThemeContext () {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')
  const [mode, setMode] = useState<PaletteMode>('dark')

  useEffect(() => {
    setMode(prefersDarkMode ? 'dark' : 'light')
  }, [prefersDarkMode])

  const colorMode = useMemo(
    () => ({
      toggleColorMode: () => {
        setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'))
      }
    }),
    []
  )

  const getPalette: (mode: PaletteMode) => PaletteOptions =
    (mode: PaletteMode) => mode === 'light'
      ? {
        // palette values for light mode
          primary: {
            main: '#fcfcfc'
          },
          secondary: {
            main: '#eb3448'
          },
          text: {
            primary: 'rgba(0,0,0,.83)',
            secondary: '#00000'
          },
          action: {
            active: '#1976D2'
          }
        }
      : {
        // palette values for dark mode
          primary: {
            main: '#1c1d22'
          },
          secondary: {
            main: '#eb3448'
          },
          text: {
            primary: '#FFFFFF',
            secondary: '#00000'
          },
          action: {
            active: '#90CAF9'
          }
        }

  const theme = useMemo(() => createTheme({
    palette: getPalette(mode)
  }), [mode])

  return (
    <ColorModeContext.Provider value={colorMode}>
      <ThemeProvider theme={theme}>
        <App />
      </ThemeProvider>
    </ColorModeContext.Provider>
  )
}
