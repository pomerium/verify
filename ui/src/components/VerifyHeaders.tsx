import React, { FC } from 'react';

import { VerifyInfo } from '../api';

type Props = {
  info: VerifyInfo;
};
const VerifyHeaders: FC<Props> = ({ info }) => {
  const headers = Object.entries(info?.headers || {}) || [];

  return (
    <div className="category white box">
      <div className="messages">
        <div className="box-inner">
          <div className="category-header clearfix">
            <span className="category-title">Unsigned Identity Headers</span>
            <a href="/headers">
              <span className="json-icon"></span>
            </a>
          </div>
          {headers.length ? (
            <table>
              <thead>
                <tr>
                  <th>Header</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {headers?.map(([k, vs]) => {
                  <tr key={k}>
                    <td>{k}</td>
                    <td>
                      {vs?.map((v) => (
                        <p key={v}>{v}</p>
                      ))}
                    </td>
                  </tr>;
                })}
              </tbody>
            </table>
          ) : (
            <>No headers found!</>
          )}
        </div>
        <div className="category-link">
          Pomerium allows{' '}
          <a href="https://docs.pomerium.io/reference/#jwt-claim-headers">
            passing identity{' '}
          </a>{' '}
          to upstream applications as HTTP request headers. Note, unlike{' '}
          <code>x-pomerium-jwt-assertion</code> these headers are
          <strong>not signed</strong>.
        </div>
      </div>
    </div>
  );
};
export default VerifyHeaders;
