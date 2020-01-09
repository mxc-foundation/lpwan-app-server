import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import AdvancedTable from "../../components/AdvancedTable";

import GatewayProfileStore from "../../stores/GatewayProfileStore";


const GatewayColumn = (cell, row, index, extraData) => {
  return <Link to={`/gateway-profiles/${row.id}`}>{row.name}</Link>;
}

const NetworkColumn = (cell, row, index, extraData) => {
  return <Link to={`/network-servers/${row.networkServerID}`}>{row.networkServerName}</Link>;
}

const getColumns = () => (
  [{
    dataField: 'name',
    text: i18n.t(`${packageNS}:tr000042`),
    sort: false,
    formatter: GatewayColumn
  }, {
    dataField: 'networkServerName',
    text: i18n.t(`${packageNS}:tr000047`),
    sort: false,
    formatter: NetworkColumn,
  }]
);

class ListGatewayProfiles extends Component {

  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage + 1;
    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    GatewayProfileStore.list(0, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    return (<React.Fragment>
      <TitleBar
        buttons={[
          <Link
            key={'b-1'}
            to={`/gateway-profiles/create`}
            className="btn btn-primary">{i18n.t(`${packageNS}:tr000277`)}
          </Link>,
        ]}
      >
        <Breadcrumb>
          <BreadcrumbItem><Link to={`/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)}</Link></BreadcrumbItem>
        </Breadcrumb>
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <AdvancedTable data={this.state.data} columns={getColumns()} keyField="id" totalSize={this.state.totalSize} onTableChange={this.handleTableChange}></AdvancedTable>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(ListGatewayProfiles);
