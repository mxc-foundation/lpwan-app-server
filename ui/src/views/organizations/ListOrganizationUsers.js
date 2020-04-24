import classNames from 'classnames';
import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Button, Card, Col, Row } from 'reactstrap';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import OrganizationStore from "../../stores/OrganizationStore";
import { MAX_DATA_LIMIT } from '../../util/pagination';





const UserNameColumn = (cell, row, index, extraData) => {
  const organizationId = extraData['organizationId'];
  return <Link to={`/organizations/${organizationId}/users/${row.userID}`}>{row.username}</Link>;
}

const AdminColumn = (cell, row, index, extraData) => {
  return <i className={classNames("mdi", {"mdi-check": row.isAdmin, "mdi-close": !row.isAdmin}, "font-20")}></i>;
}

const getColumns = (organizationId) => (
  [{
    dataField: 'userID',
    text: i18n.t(`${packageNS}:tr000077`),
    sort: false,
  }, {
    dataField: 'username',
    text: i18n.t(`${packageNS}:tr000056`),
    sort: false,
    formatter: UserNameColumn,
    formatExtraData: { organizationId: organizationId }
  }, {
    dataField: 'isAdmin',
    text: i18n.t(`${packageNS}:tr000058`),
    sort: false,
    formatter: AdminColumn,
  }]
);

class ListOrganizationUsers extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data: [],
      totalSize: 0,
      loading: false,
    }
  }

  createUser = () => {
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users/create`);
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage;

    /* let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    } */
    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  }

  /**
   * Fetches data from server
   */
  getPage = async (limit, offset) => {
    limit = MAX_DATA_LIMIT;
    this.setState({ loading: true });
    
    const res = await OrganizationStore.listUsers(this.props.match.params.organizationID, limit=10, offset=0);
console.log(res);    
    const object = this.state;
    object.totalSize = Number(res.totalCount);
    object.data = res.result;
    object.loading = false;
    this.setState({ object });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Button color="primary"
              key={1}
              onClick={this.createUser}
              className=""><i className="mdi mdi-plus"></i>{' '}{i18n.t(`${packageNS}:tr000041`)}
            </Button>,
          ]}
        >
          <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000068`), active: true }]}></OrgBreadCumb>
        </TitleBar>
        <Row>
          <Col>
            <Card className="card-box shadow-sm">
            {this.state.loading && <Loader />}
              <AdvancedTable data={this.state.data} columns={getColumns(this.props.match.params.organizationID)}
                keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} totalSize={this.state.totalSize} rowsPerPage={10}></AdvancedTable>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(ListOrganizationUsers);
