import React from 'react';
import { Link } from "react-router-dom";
import { Card, CardBody, CardFooter, CardHeader, CardText } from 'reactstrap';
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';



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
