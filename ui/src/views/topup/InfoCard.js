import React from 'react';
import { Row, Col, Card, CardHeader, CardFooter, CardBody, CardText } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBarTitle from "../../components/TitleBarTitle";
import { Link } from "react-router-dom";


export default function MediaCard(props) {

  return (
    <React.Fragment>
      <Card style={{ height: '100%' }}>
        <CardHeader>{i18n.t(`${packageNS}:menu.topup.synchronize_your_eth_account`)}</CardHeader>
        <CardBody>
          <CardText>{i18n.t(`${packageNS}:menu.topup.note`)}</CardText>
        </CardBody>
        <CardFooter>
          <TitleBarTitle component={Link} to={`${props.path}`} title={i18n.t(`${packageNS}:menu.topup.change_eth_account`)} />
        </CardFooter>
      </Card>
    </React.Fragment>
  );
}
