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

  getColumns = (canHaveGateways) => (
      [{
            dataField: 'name',
            text: i18n.t(`${packageNS}:tr000042`),
            sort: false,
            formatter: this.organizationNameColumn,
        },
        {
            dataField: 'display_name',
            text: i18n.t(`${packageNS}:tr000126`),
            sort: false,
        },
        {
            dataField: 'can_have_gateways',
            text: i18n.t(`${packageNS}:tr000380`),
            sort: false,
            formatter: this.canHaveGatewaysColumn,
            formatExtraData: {canHaveGateways: canHaveGateways},
        },
        {
            dataField: 'service_profiles',
            text: i18n.t(`${packageNS}:tr000078`),
            sort: false,
           /* formatter: this.serviceProfileColumn,*/
        }]
  )

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage + 1;
    this.getPage(this.props.match.params.organizationID, sizePerPage, offset);
  }
  getPage(limit, offset, callbackFunc) {
    OrganizationStore.list("", limit, offset,  (res) => {
      this.setState({ data: res.result });
    });
  }

  organizationNameColumn = (cell, row, index, extraData) => {
    return <Link to={`/organizations/${this.props.match.params.organizationID}`}>{row.name}</Link>;
  }

  canHaveGatewaysColumn = (cell, row, index, extraData) => {
      if (extraData.canHaveGateways) {
          return  <Check />;
      } else {
          return  <Close />;
      }
  }

  componentDidMount() {
    this.getPage(this.props.match.params.organizationID, 10);
  }

  render() {

    return (<React.Fragment>
      <TitleBar>
        <TitleBarTitle title={i18n.t(`${packageNS}:tr000069`)} />
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <AdvancedTable data={this.state.data} columns={this.getColumns(this.state.data.canHaveGateways)} keyField="id" onTableChange={this.handleTableChange}></AdvancedTable>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListOrganizations;
