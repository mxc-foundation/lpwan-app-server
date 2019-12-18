import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Button, Card, CardBody,
  CardSubtitle, CardTitle, Col, Container, Row, Spinner } from 'reactstrap';

import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import TableCellLink from "../../components/TableCellLink";
import DataTable from "../../components/DataTable";
import AdvancedTable from "../../components/AdvancedTable";

import NetworkServerStore from "../../stores/NetworkServerStore";

const NetworkServerColumn = (cell, row, index, extraData) => {
  return <Link to={`/network-servers/${row.id}`}>{row.name}</Link>;
}

const NetworkServerAddressColumn = (cell, row, index, extraData) => {
  return <div>{row.server}</div>;
}

const columns = [{
  dataField: 'networkServerName',
  text: i18n.t(`${packageNS}:tr000042`),
  sort: false,
  formatter: NetworkServerColumn
}, {
  dataField: 'networkServerAddress',
  text: i18n.t(`${packageNS}:tr000043`),
  sort: false,
  formatter: NetworkServerAddressColumn
}];

class ListNetworkServers extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      loading: true
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage + 1;

    let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    }

    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    const defaultOrgId = 0;
    this.setState({ loading: true });
    NetworkServerStore.list(defaultOrgId, limit, offset, (res) => {
      this.setState({
        data: res.result,
        loading: false
      });
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    return(
      <React.Fragment>
        <TitleBar
          buttons={[
            <Link
              aria-label={i18n.t(`${packageNS}:tr000277`)}
              key={'b-1'}
              to={`/network-servers/create`}
              className="btn btn-primary">{i18n.t(`${packageNS}:tr000277`)}
            </Link>,
          ]}
        >
          <Breadcrumb style={{ fontSize: "1.25rem", margin: "0rem" }}>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000040`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card className="shadow-sm">
              <CardBody className="text-center">
                {
                  this.state.data.length && this.state.data[0].name
                  ? <AdvancedTable
                      data={this.state.data}
                      columns={columns}
                      keyField="id"
                      onTableChange={this.handleTableChange}
                      rowsPerPage={10}
                      searchEnabled={false}
                    />
                  : <Spinner color="primary" style={{ align: 'center', width: '3rem', height: '3rem' }} />
                }
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default ListNetworkServers;
