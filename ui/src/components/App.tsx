import ThemeProvider from '@mui/material/styles/ThemeProvider';
import createTheme from '@mui/material/styles/createTheme';
import React, { FC } from 'react';

import Verify from './Verify';
import WebAuthn from './WebAuthn';

const theme = createTheme({
  palette: {
    action: {
      active: '#39256C'
    },
    background: {
      default: '#FFFFFF',
      paper: '#FFFFFF'
    },
    primary: {
      main: '#6F43E7'
    },
    secondary: {
      main: '#49AAA1'
    }
  }
});

const App: FC = () => {
  return (
    <ThemeProvider theme={theme}>
      {location.href.indexOf('webauthn') >= 0 ? <WebAuthn /> : <Verify />}
    </ThemeProvider>
  );
};
export default App;
