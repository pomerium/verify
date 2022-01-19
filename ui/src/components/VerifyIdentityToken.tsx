import React, { FC } from 'react';
import uniq from 'lodash/uniq';
import isString from 'lodash/isString';

import { VerifyInfo } from '../api';

const makeUnique = (entries: string|string[]): string[] => {
  if (!entries) {
    return [];
  }
  if (isString(entries)) {
    return [entries];
  }
  return uniq(entries);
};

type Props = {
  info: VerifyInfo;
};
const VerifyIdentityToken: FC<Props> = ({ info }) => {
  return (
    <div className="category white box">
      <div className="messages">
        <div className="box-inner">
          <div className="category-header clearfix">
            <span className="category-title">Signed Identity Token</span>
            {!info?.error &&
            info?.identity?.raw_jwt &&
            info?.identity?.public_key ? (
              <>
                <a
                  href={`https://jwt.io/#debugger-io?token=${info?.identity?.raw_jwt}&publicKey=${info?.identity?.public_key}`}
                >
                  <span className="category-icon"> </span>
                </a>
              </>
            ) : (
              <></>
            )}
          </div>
          <ul className="category-contents plain">
            <table>
              <thead>
                <tr>
                  <th>Claim</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Email</td>
                  <td>{info?.identity?.email}</td>
                </tr>
                <tr>
                  <td>Issuer</td>
                  <td>{info?.identity?.iss}</td>
                </tr>
                <tr>
                  <td>User</td>
                  <td>{info?.identity?.user}</td>
                </tr>
                <tr>
                  <td>Subject</td>
                  <td>{info?.identity?.sub}</td>
                </tr>
                <tr>
                  <td>Audience</td>
                  <td>
                    {makeUnique(info?.identity?.aud).map((v) => (
                      <p key={v}>{v}</p>
                    ))}
                  </td>
                </tr>
                <tr>
                  <td>Expiry</td>
                  <td>{info?.identity?.exp}</td>
                </tr>
                <tr>
                  <td>IssuedAt</td>
                  <td>{info?.identity?.iat}</td>
                </tr>
                <tr>
                  <td>ID</td>
                  <td>{info?.identity?.jti}</td>
                </tr>
                <tr>
                  <td>Groups</td>
                  <td>
                    {makeUnique(info?.identity?.groups).map((v) => (
                      <p key={v}>{v}</p>
                    ))}
                  </td>
                </tr>
              </tbody>
            </table>
          </ul>
        </div>
        <div className="category-link">
          Pomerium adds a signed JWT token to the incoming request headers (
          <code>X-Pomerium-Jwt-Assertion</code>) which can then be used to
          assert a {" "}
          <a href="https://www.pomerium.com/docs/topics/getting-users-identity.html#verification">
            user's identity details
          </a>
          .
        </div>
      </div>
    </div>
  );
};
export default VerifyIdentityToken;
