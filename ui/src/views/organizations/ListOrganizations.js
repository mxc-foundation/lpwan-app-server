import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Row, Col, Card, CardBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import AdvancedTable from "../../components/AdvancedTable";

import OrganizationStore from "../../stores/OrganizationStore";
import Check from "mdi-material-ui/Check";
import Close from "mdi-material-ui/Close";
import TitleBarButton from "../../components/TitleBarButton";

class ListOrganizations extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);

    this.getColumns = this.getColumns.bind(this);
    this.organizationNameColumn = this.organizationNameColumn.bind(this);
    this.canHaveGatewaysColumn = this.canHaveGatewaysColumn.bind(this);

    this.state = {
      data: []
    }
  }

  organizationNameColumn = (cell, row, index, extraData) => {
    return <Link to={`/organizations/${row.id}`}>{row.name}</Link>;
  };

  canHaveGatewaysColumn = (cell, row, index, extraData) => {
    if (row.canHaveGateways) {
        return  <Check />;
    } else {
        return  <Close />;
    }
  };

  serviceProfileColumn = (cell, row, index, extraData) => {
      return <div>
          <div>
              <Link to={`/organizations/${row.id}/service-profiles/create`}>ADD</Link>
          </div>
          <div>
            <Link to={`/organizations/${row.id}/service-profiles`}>CHECK</Link>
          </div>
      </div>;
  };

  getColumns = () => (
      [{
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000042`),
        sort: false,
        formatter: this.organizationNameColumn,
        },
        {
        dataField: 'displayName',
        text: i18n.t(`${packageNS}:tr000126`),
        sort: false,
        },
        {
        dataField: 'canHaveGateways',
        text: i18n.t(`${packageNS}:tr000380`),
        sort: false,
        formatter: this.canHaveGatewaysColumn,
        },
        {
        dataField: 'serviceProfiles',
        text: i18n.t(`${packageNS}:tr000078`),
        sort: false,
        formatter: this.serviceProfileColumn,
        }]
  );

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage + 1;
    this.getPage(sizePerPage, offset);
  };

  getPage(limit, offset) {
    OrganizationStore.list("", limit, offset,  (res) => {
      this.setState({ data: res.result });
    });
  }

  componentDidMount() {
    this.getPage(10, 0);
  }

  render() {
    return (<React.Fragment>
    <TitleBar buttons={
        <TitleBarButton
            key={1}
            label={i18n.t(`${packageNS}:tr000277`)}
            icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
            to={`/organizations/create`}
        />}
    >
        <TitleBarTitle title={i18n.t(`${packageNS}:tr000049`)} />
    </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <AdvancedTable data={this.state.data} columns={this.getColumns()} keyField="id" onTableChange={this.handleTableChange}></AdvancedTable>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListOrganizations;
