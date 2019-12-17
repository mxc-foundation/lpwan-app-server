import React, { Component } from "react";
import { Link } from "react-router-dom";

import Check from "mdi-material-ui/Check";
import Close from "mdi-material-ui/Close";
import AdvancedTable from "../../components/AdvancedTable";
import { Button, Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import UserStore from "../../stores/UserStore";

const GatewayColumn = (cell, row, index, extraData) => {
  const organizationId = extraData['organizationId'];
  return <Link to={`/users/${row.id}`}>{row.username}</Link>;
}

const ActiveColumn = (cell, row, index, extraData) => {
  if (row.isActive) {
    return <Check/>;
  } else {
    return <Close />;
  }
}

const AdminColumn = (cell, row, index, extraData) => {
  if (row.isAdmin) {
    return <Check />;
  } else {
    return <Close />;
  }
}

const getColumns = (organizationId) => (
  [{
    dataField: 'username',
    text: i18n.t(`${packageNS}:tr000056`),
    sort: false,
    formatter: GatewayColumn,
    formatExtraData: { organizationId: organizationId }
  }, {
    dataField: 'isActive',
    text: i18n.t(`${packageNS}:tr000057`),
    sort: false,
    formatter: ActiveColumn,
  }, {
    dataField: 'isAdmin',
    text: i18n.t(`${packageNS}:tr000058`),
    sort: false,
    formatter: AdminColumn,
  }]
);

class ListUsers extends Component {
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
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage + 1;

    let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    }
    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = (limit, offset) => {
    this.setState({ loading: true });
    UserStore.list("", limit, offset, (res) => {
      this.setState({ data: res.result, loading: false });
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  createUser = () => {
    this.props.history.push(`/users/create`);
  }

  render() {
    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Button color="primary"
              key={1}
              onClick={this.createUser}
              className=""><i class="mdi mdi-account-multiple-plus"></i>{' '}{i18n.t(`${packageNS}:tr000277`)}
            </Button>,
          ]}
        >
          <Breadcrumb>
            <BreadcrumbItem><Link to={`/users`}>{i18n.t(`${packageNS}:tr000036`)}</Link></BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card>
              <CardBody>
                <AdvancedTable data={this.state.data} columns={getColumns(this.props.organizationID)}
                  keyField="id" onTableChange={this.handleTableChange} searchEnabled={true} rowsPerPage={10}></AdvancedTable>

              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default ListUsers;
