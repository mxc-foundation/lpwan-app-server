import React, { Component } from "react";

import {
  Button,
  Card,
  CardSubtitle,
  CardTitle,
  Col,
  Container,
  Row
} from 'reactstrap';

import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import i18n, { packageNS } from '../../i18n';
import TableCellLink from "../../components/TableCellLink";
import DataTable from "../../components/DataTable";

import NetworkServerStore from "../../stores/NetworkServerStore";


class ListNetworkServers extends Component {
  getPage(limit, offset, callbackFunc) {
    NetworkServerStore.list(0, limit, offset, callbackFunc);
  }

  getRow(obj) {
    return(
      <TableRow key={obj.id}>
        <TableCellLink to={`/network-servers/${obj.id}`}>{obj.name}</TableCellLink>
        <TableCell>{obj.server}</TableCell>
      </TableRow>
    );
  }

  render() {
    return(
      <Container>
        <Row>
          <Col md="12" sm="12">
            <Card className="card-box" style={{ minWidth: "25rem" }}>
              <Row>
                <Col md="8" xs="9">
                  <CardTitle className="mt-0 header-title">
                    {i18n.t(`${packageNS}:tr000040`)}
                  </CardTitle>
                  <CardSubtitle className="text-muted font-14 mb-3">
                    List of network servers.
                  </CardSubtitle>
                </Col>
                <Col md="2" xs="3">
                  <Button
                    color="primary"
                    href={`/network-servers/create`}
                    outline
                    size="md"
                  >
                    {i18n.t(`${packageNS}:tr000041`)}
                  </Button>
                </Col>
              </Row>
              <Row>
                <Col md="12" sm="12">
                  <DataTable
                    header={
                      <TableRow>
                        <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                        <TableCell>{i18n.t(`${packageNS}:tr000043`)}</TableCell>
                      </TableRow>
                    }
                    getPage={this.getPage}
                    getRow={this.getRow}
                  />
                </Col>
              </Row>
            </Card>
          </Col>
        </Row>
      </Container>
    );
  }
}

export default ListNetworkServers;
