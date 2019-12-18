import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Row, Col, Card, CardBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import AdvancedTable from "../../components/AdvancedTable";

import ServiceProfileStore from "../../stores/ServiceProfileStore";


const ServiceProfileColumn = (cell, row, index, extraData) => {
  return <Link to={`/organizations/2/service-profiles/${row.id}`}>{row.name}</Link>;
}

const columns = [{
    dataField: 'test_gateway_profile',
    text: i18n.t(`${packageNS}:tr000042`),
    sort: false,
    formatter: ServiceProfileColumn
}];

class ListServiceProfiles extends Component {

  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: []
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
    ServiceProfileStore.list(0, limit, offset, (res) => {
      this.setState({ data: res.result });
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
            to={`/organizations/${this.props.match.params.organizationID}/service-profiles/create`}
            className="btn btn-primary">{i18n.t(`${packageNS}:tr000277`)}
          </Link>,
        ]}
      >
        <TitleBarTitle title={i18n.t(`${packageNS}:tr000069`)} />
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <AdvancedTable data={this.state.data} columns={columns} keyField="id" onTableChange={this.handleTableChange}></AdvancedTable>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ListServiceProfiles;
