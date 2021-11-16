import React, { FC, useEffect, useState } from 'react';

import { fetchVerifyInfo, VerifyInfo } from '../api';
import VerifyHeaders from './VerifyHeaders';
import VerifyIdentityToken from './VerifyIdentityToken';
import VerifyRequestDetails from './VerifyRequestDetails';
import VerifyStatus from './VerifyStatus';

const Verify: FC = () => {
  const [info, setInfo] = useState<VerifyInfo>(null);

  useEffect(() => {
    (async () => {
      setInfo(await fetchVerifyInfo());
    })();
  }, []);

  return (
    <div className="inner">
      <div className="header clearfix">
        <div className="heading">
          <a href="/" className="logo"></a> <span>{info?.request?.host}</span>
        </div>
      </div>

      <div className="content">
        <VerifyStatus info={info} />
        <VerifyIdentityToken info={info} />
        <VerifyHeaders info={info} />
        <VerifyRequestDetails info={info} />
      </div>
    </div>
  );
};
export default Verify;
