import React, { FC } from 'react';

import { VerifyInfo } from '../api';

type Props = {
  info: VerifyInfo;
};
const VerifyRequestDetails: FC<Props> = ({ info }) => {
  const request = Object.entries(info?.request || {});

  return (
    <div className="category white box">
      <div className="messages">
        <div className="box-inner">
          <div className="category-header clearfix">
            <span className="category-title">Request Details</span>
            <a href="/json">
              <span className="json-icon"></span>
            </a>
          </div>

          <table>
            <tbody>
              {request.map(([k, v]) => (
                <tr key={k}>
                  <td>{k}</td>
                  <td>{v}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="category-link">
          A complete dump of the values on this page can be found at the{' '}
          <a href="/json">/json</a> endpoint.
        </div>
      </div>
    </div>
  );
};
export default VerifyRequestDetails;
