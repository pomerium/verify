import React, { FC } from 'react';

import { VerifyInfo } from 'src/api';

type Props = {
  info: VerifyInfo;
};
const VerifyStatus: FC<Props> = ({ info }) => {
  return (
    <div className="white box">
      <div className="largestatus">
        {info?.error ? (
          <>
            <span className="status-bubble status-down"> </span>
            <div className="title-wrapper">
              <span className="title">Identity verification failed</span>
              <label className="status-time">
                <span>
                  We tried to verify the incoming user, but failed with the
                  following error:{' '}
                </span>
                <code>{info?.error}</code>
              </label>
            </div>
          </>
        ) : !info?.request?.tlsValid ? (
          <>
            <span className="status-bubble status-warn"> </span>
            <div className="title-wrapper">
              <span className="title">TLS Certificate verification failed</span>
              <label className="status-time">
                <span>
                  TLS certificate verification failed when verifying the JWT.
                </span>
              </label>
            </div>
          </>
        ) : (
          <>
            <span className="status-bubble status-up"> </span>
            <div className="title-wrapper">
              <span className="title">Identity found and verified 🚀</span>
              <label className="status-time">
                <span>
                  This app is properly configured and is being secured by
                  Pomerium.
                </span>
              </label>
            </div>
          </>
        )}
      </div>
    </div>
  );
};
export default VerifyStatus;
