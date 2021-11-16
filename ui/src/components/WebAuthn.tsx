import { decode, encodeUrl } from '@borderless/base64';
import Alert, { AlertColor } from '@mui/material/Alert';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import FormControl from '@mui/material/FormControl';
import Grid from '@mui/material/Grid';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import React, { FC, useState } from 'react';

import {
  toWebAuthnAuthenticateOptions,
  toWebAuthnRegisterOptions,
  webAuthnAuthenticate,
  webAuthnRegister
} from 'src/api';

function getKnownCredentials(): ArrayBuffer[] {
  let raw = [];
  try {
    raw =
      (JSON.parse(localStorage.getItem('known-credentials')) as string[]) || [];
  } catch (e) {}
  return raw.map((v) => decode(v));
}

function addKnownCredential(rawID: ArrayBuffer) {
  let raw = [];
  try {
    raw =
      (JSON.parse(localStorage.getItem('known-credentials')) as string[]) || [];
  } catch (e) {}
  const set = new Set<string>(raw);
  set.add(encodeUrl(rawID));
  localStorage.setItem(
    'known-credentials',
    JSON.stringify(Array.from(set.values()))
  );
}

async function authenticate(username: string) {
  const challenge = crypto.getRandomValues(new Uint8Array(32));
  const options: PublicKeyCredentialRequestOptions = {
    allowCredentials: getKnownCredentials().map((rawID) => ({
      type: 'public-key',
      id: rawID
    })),
    challenge: challenge,
    rpId: location.hostname,
    userVerification: 'preferred'
  };
  const credential = (await navigator.credentials.get({
    publicKey: options
  })) as PublicKeyCredential;
  const credentialResponse =
    credential.response as AuthenticatorAssertionResponse;
  await webAuthnAuthenticate({
    options: toWebAuthnAuthenticateOptions(options),
    credential: {
      id: credential.id,
      type: credential.type,
      rawId: encodeUrl(credential.rawId),
      response: {
        authenticatorData: encodeUrl(credentialResponse.authenticatorData),
        clientDataJSON: encodeUrl(credentialResponse.clientDataJSON),
        signature: encodeUrl(credentialResponse.signature),
        userHandle: encodeUrl(
          credentialResponse.userHandle ||
            Uint8Array.from(username, (c) => c.charCodeAt(0))
        )
      }
    }
  });
}

async function register(
  username: string,
  attestationType: AttestationConveyancePreference,
  authenticatorAttachment: AuthenticatorAttachment
) {
  const challenge = crypto.getRandomValues(new Uint8Array(32));
  const options: PublicKeyCredentialCreationOptions = {
    attestation: attestationType || undefined,
    authenticatorSelection: {
      authenticatorAttachment: authenticatorAttachment || undefined,
      residentKey: 'preferred',
      userVerification: 'preferred'
    },
    challenge: challenge,
    pubKeyCredParams: [
      { alg: -65535, type: 'public-key' },
      { alg: -257, type: 'public-key' },
      { alg: -7, type: 'public-key' }
    ],
    rp: {
      id: location.hostname,
      name: 'Pomerium'
    },
    user: {
      id: Uint8Array.from(username, (c) => c.charCodeAt(0)),
      name: username,
      displayName: username
    }
  };
  const credential = (await navigator.credentials.create({
    publicKey: options
  })) as PublicKeyCredential;
  const credentialResponse =
    credential.response as AuthenticatorAttestationResponse;
  addKnownCredential(credential.rawId);
  await webAuthnRegister({
    options: toWebAuthnRegisterOptions(options),
    credential: {
      id: credential.id,
      type: credential.type,
      rawId: encodeUrl(credential.rawId),
      response: {
        attestationObject: encodeUrl(credentialResponse.attestationObject),
        clientDataJSON: encodeUrl(credentialResponse.clientDataJSON)
      }
    }
  });
}

type ActionResult = {
  severity: AlertColor;
  message: string;
};

const WebAuthn: FC = ({}) => {
  const [username, setUsername] = useState<string>(null);
  const [attestationType, setAttestationType] =
    useState<AttestationConveyancePreference>(null);
  const [authenticatorAttachment, setAuthenticatorAttachment] =
    useState<AuthenticatorAttachment>(null);
  const [result, setResult] = useState<ActionResult>(null);
  const knownCredentials = getKnownCredentials();

  function onChangeUsername(evt: React.ChangeEvent<HTMLInputElement>) {
    setUsername(evt.target.value);
  }
  function onChangeAttestationType(
    evt: SelectChangeEvent<AttestationConveyancePreference>
  ) {
    setAttestationType(
      evt.target.value === 'none'
        ? null
        : (evt.target.value as AttestationConveyancePreference)
    );
  }
  function onChangeAuthenticatorType(
    evt: SelectChangeEvent<AuthenticatorAttachment>
  ) {
    setAuthenticatorAttachment(
      evt.target.value === 'unspecified'
        ? null
        : (evt.target.value as AuthenticatorAttachment)
    );
  }
  function onClickLogin(evt: React.MouseEvent<HTMLButtonElement>) {
    evt.preventDefault();

    (async () => {
      setResult(null);
      try {
        await authenticate(username);
        setResult({
          severity: 'success',
          message: `Authentication Successful!`
        });
      } catch (e) {
        setResult({ severity: 'error', message: `${e}` });
      }
    })();
  }
  function onClickRegister(evt: React.MouseEvent<HTMLButtonElement>) {
    evt.preventDefault();

    (async () => {
      setResult(null);
      try {
        await register(username, attestationType, authenticatorAttachment);
        setResult({
          severity: 'success',
          message: `Registration Successful! Now try Login.`
        });
      } catch (e) {
        setResult({ severity: 'error', message: `${e}` });
      }
    })();
  }

  return (
    <div className="inner">
      <div className="header clearfix">
        <div className="heading">
          <a href="/" className="logo"></a> <span>WebAuthn</span>
        </div>
      </div>

      <div className="category white box">
        <div className="messages">
          <div className="box-inner">
            <div className="category-header clearfix">
              <span className="category-title"></span>
              <a href="/json">
                <span className="webauthn-icon"></span>
              </a>
            </div>

            <form>
              <Container>
                <Grid container spacing={3}>
                  <Grid item xs={12}>
                    <TextField
                      fullWidth
                      label="Username"
                      onChange={onChangeUsername}
                      value={username || ''}
                      variant="outlined"
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <FormControl fullWidth>
                      <InputLabel>Attestation Type</InputLabel>
                      <Select
                        label="Attestation Type"
                        onChange={onChangeAttestationType}
                        value={attestationType || 'none'}
                      >
                        <MenuItem value="none">None</MenuItem>
                        <MenuItem value="indirect">Indirect</MenuItem>
                        <MenuItem value="direct">Direct</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                  <Grid item xs={12}>
                    <FormControl fullWidth>
                      <InputLabel>Authenticator Type</InputLabel>
                      <Select
                        label="Authenticator Type"
                        onChange={onChangeAuthenticatorType}
                        value={authenticatorAttachment || 'unspecified'}
                      >
                        <MenuItem value="unspecified">Unspecified</MenuItem>
                        <MenuItem value="cross-platform">
                          Cross-Platform
                        </MenuItem>
                        <MenuItem value="platform">Platform (TPM)</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                  <Grid item xs={6}>
                    <Button
                      disabled={!username}
                      onClick={onClickRegister}
                      variant="contained"
                      fullWidth
                    >
                      Register
                    </Button>
                  </Grid>
                  <Grid item xs={6}>
                    <Button
                      disabled={!username || !knownCredentials?.length}
                      onClick={onClickLogin}
                      variant="contained"
                      fullWidth
                    >
                      Login
                    </Button>
                  </Grid>
                  {result ? (
                    <Grid item xs={12}>
                      <Alert severity={result.severity}>{result.message}</Alert>
                    </Grid>
                  ) : (
                    <></>
                  )}
                </Grid>
              </Container>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
export default WebAuthn;
