import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Row, Col, Card, CardBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import AdvancedTable from "../../components/AdvancedTable";

import ServiceProfileStore from "../../stores/ServiceProfileStore";


class ListServiceProfiles extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.serviceProfileColumn = this.serviceProfileColumn.bind(this)
    this.state = {
      data: [],
      columns: [{
        dataField: 'name',
        text: i18n.t(`${packageNS}:tr000042`),
        sort: false,
        formatter: this.serviceProfileColumn,
      }],
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage + 1;
    this.getPage(this.props.match.params.organizationID, sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (organizationID, limit, offset) => {
    ServiceProfileStore.list(organizationID, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  serviceProfileColumn = (cell, row, index, extraData) => {
    return <Link to={`/organizations/${this.props.match.params.organizationID}/service-profiles/${row.id}`}>{row.name}</Link>;
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
              <AdvancedTable data={this.state.data} columns={this.state.columns} keyField="id" totalSize={this.state.totalSize} onTableChange={this.handleTableChange}></AdvancedTable>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListServiceProfiles;
