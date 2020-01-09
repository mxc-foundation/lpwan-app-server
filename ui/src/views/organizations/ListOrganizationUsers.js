import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Check from "mdi-material-ui/Check";
import Close from "mdi-material-ui/Close";
import AdvancedTable from "../../components/AdvancedTable";
import { Breadcrumb, BreadcrumbItem, Button, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TitleBar from "../../components/TitleBar";

import OrganizationStore from "../../stores/OrganizationStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

const UserNameColumn = (cell, row, index, extraData) => {
  const organizationId = extraData['organizationId'];
  return <Link to={`/organizations/${organizationId}/users/${row.userID}`}>{row.username}</Link>;
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
      totalSize: 0
    }
  }

  createUser = () => {
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users/create`);
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
    OrganizationStore.listUsers(this.props.match.params.organizationID, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = res.totalCount;
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Button color="primary"
              key={1}
              onClick={this.createUser}
              className=""><i className="mdi mdi-account-multiple-plus"></i>{' '}{i18n.t(`${packageNS}:tr000041`)}
            </Button>,
          ]}
        >
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
              >
                  Organizations
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}`}
              >
                {currentOrgID}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000068`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col>
            <Card>
              <CardBody>
                <AdvancedTable data={this.state.data} columns={getColumns(this.props.match.params.organizationID)}
                  keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} totalSize={this.state.totalSize} rowsPerPage={10}></AdvancedTable>

              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(ListOrganizationUsers));
