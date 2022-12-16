import { encodeUrl } from '@borderless/base64';

export type VerifyInfoRequest = {
  host: string;
  hostname: string;
  method: string;
  origin: string;
  url: string;
  tlsValid: boolean;
};
export type VerifyInfoIdentity = {
  iss?: string;
  sub?: string;
  aud?: string | string[];
  exp?: number;
  nbf?: number;
  iat?: number;
  jti?: string;
  groups?: string[];
  user?: string;
  email?: string;
  raw_jwt?: string;
  public_key?: string;
};
export type VerifyInfo = {
  error?: string;
  identity?: VerifyInfoIdentity;
  headers: { [name: string]: string[] };
  request: VerifyInfoRequest;
};

export async function fetchVerifyInfo(): Promise<VerifyInfo> {
  const response = await fetch('/api/verify-info', {
    headers: {
      'Content-Type': 'application/json'
    }
  });
  const result = await response.json();
  return result as VerifyInfo;
}

export type WebAuthnCredentialDescriptor = {
  type: string;
  id: string;
};

export type WebAuthnAuthenticateOptions = {
  allowCredentials?: WebAuthnCredentialDescriptor[];
  challenge: string;
  extensions?: Record<string, unknown>;
  rpId?: string;
  timeout?: number;
  userVerification?: string;
};

export function toWebAuthnAuthenticateOptions(
  options: PublicKeyCredentialRequestOptions
): WebAuthnAuthenticateOptions {
  let obj: WebAuthnAuthenticateOptions = {
    challenge: encodeUrl(
      ArrayBuffer.isView(options.challenge)
        ? options.challenge.buffer
        : options.challenge
    )
  };
  if ('allowCredentials' in options) {
    obj.allowCredentials = options.allowCredentials.map((c) => ({
      id: encodeUrl(ArrayBuffer.isView(c.id) ? c.id.buffer : c.id),
      type: c.type
    }));
  }
  if ('extensions' in options) {
    obj.extensions = Object.fromEntries(Object.entries(options.extensions));
  }
  if ('rpId' in options) {
    obj.rpId = options.rpId;
  }
  if ('timeout' in options) {
    obj.timeout = options.timeout;
  }
  if ('userVerification' in options) {
    obj.userVerification = options.userVerification;
  }
  return obj;
}

export type WebAuthnAuthenticateRequest = {
  options: WebAuthnAuthenticateOptions;
  credential: {
    id: string;
    type: string;
    rawId: string;
    response: {
      authenticatorData: string;
      clientDataJSON: string;
      signature: string;
      userHandle: string;
    };
  };
};

export async function webAuthnAuthenticate(
  request: WebAuthnAuthenticateRequest
) {
  const response = await fetch('/api/webauthn-authenticate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(request)
  });
  const text = await response.text();
  if (!response.ok) {
    throw `${response.status}: ${text}`;
  }
}

export type WebAuthnRegisterOptions = {
  attestation?: string;
  authenticatorSelection?: {
    authenticatorAttachment?: string;
    requireResidentKey?: boolean;
    residentKey?: string;
    userVerification?: string;
  };
  challenge: string;
  excludeCredentials?: WebAuthnCredentialDescriptor[];
  extensions?: Record<string, unknown>;
  pubKeyCredParams: PublicKeyCredentialParameters[];
  rp: {
    id: string;
    name: string;
  };
  timeout?: number;
  user: {
    id: string;
    displayName: string;
    name: string;
  };
};

export function toWebAuthnRegisterOptions(
  options: PublicKeyCredentialCreationOptions
): WebAuthnRegisterOptions {
  let obj: WebAuthnRegisterOptions = {
    challenge: encodeUrl(
      ArrayBuffer.isView(options.challenge)
        ? options.challenge.buffer
        : options.challenge
    ),
    pubKeyCredParams: options.pubKeyCredParams,
    rp: {
      id: options.rp.id,
      name: options.rp.name
    },
    user: {
      id: encodeUrl(
        ArrayBuffer.isView(options.user.id)
          ? options.user.id.buffer
          : options.user.id
      ),
      displayName: options.user.displayName,
      name: options.user.name
    }
  };
  if ('attestation' in options) {
    obj.attestation = options.attestation;
  }
  if ('authenticatorSelection' in options) {
    obj.authenticatorSelection = {};
    [
      'authenticatorAttachment',
      'requireResidentKey',
      'residentKey',
      'userVerification'
    ].forEach((k) => {
      if (k in options.authenticatorSelection) {
        obj.authenticatorSelection[k] = options.authenticatorSelection[k];
      }
    });
  }
  if ('excludeCredentials' in options) {
    obj.excludeCredentials = options.excludeCredentials.map((c) => ({
      id: encodeUrl(ArrayBuffer.isView(c.id) ? c.id.buffer : c.id),
      type: c.type
    }));
  }
  if ('extensions' in options) {
    obj.extensions = Object.fromEntries(Object.entries(options.extensions));
  }
  if ('timeout' in options) {
    obj.timeout = options.timeout;
  }
  return obj;
}

export type WebAuthnRegisterRequest = {
  options: WebAuthnRegisterOptions;
  credential: {
    id: string;
    type: string;
    rawId: string;
    response: {
      attestationObject: string;
      clientDataJSON: string;
    };
  };
};

export async function webAuthnRegister(request: WebAuthnRegisterRequest) {
  const response = await fetch('/api/webauthn-register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(request)
  });
  const text = await response.text();
  if (!response.ok) {
    throw `${response.status}: ${text}`;
  }
}
