import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Card, CardBody, Row, Col } from 'reactstrap';
import Grid from "@material-ui/core/Grid";

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import Admin from "../../components/Admin";
import ApplicationStore from "../../stores/ApplicationStore";
import OrganizationDevices from "../devices/OrganizationDevices";

const ApplicationNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  return <Link to={`/organizations/${currentOrgID}/applications/${row.id}`}>{row.name}</Link>;
}

const ApplicationServiceProfileNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  return <Link to={`/organizations/${currentOrgID}/service-profiles/${row.serviceProfileID}`}>{row.serviceProfileName}</Link>;
}

const ApplicationDescriptionColumn = (cell, row, index, extraData) => {
  return <div>{row.description}</div>;
}

const getColumns = (currentOrgID) => (
  [
    {
      dataField: 'id',
      text: i18n.t(`${packageNS}:tr000077`),
      sort: false
    },
    {
      dataField: 'name',
      text: i18n.t(`${packageNS}:tr000042`),
      sort: false,
      formatter: ApplicationNameColumn,
      formatExtraData: { currentOrgID: currentOrgID }
    }, {
      dataField: 'serviceProfileName',
      text: i18n.t(`${packageNS}:tr000078`),
      sort: false,
      formatter: ApplicationServiceProfileNameColumn,
      formatExtraData: { currentOrgID: currentOrgID }
    }, {
      dataField: 'description',
      text: i18n.t(`${packageNS}:tr000079`),
      sort: false,
      formatter: ApplicationDescriptionColumn
    }
  ]
);

class ListApplications extends Component {
  constructor() {
    super();

    this.state = {
      data: [],
      loading: true,
      totalSize: 0
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage ;

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
    limit = MAX_DATA_LIMIT;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    this.setState({ loading: true });

    ApplicationStore.list("", currentOrgID, limit, offset, (res) => {
      const object = this.state;
      object.totalSize = Number(res.totalCount);
      object.data = res.result;
      object.loading = false;
      this.setState({object});
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <Grid container spacing={4}>
        <OrganizationDevices
          mainTabIndex={1}
          organizationID={currentOrgID}
        >
          <TitleBar
            buttons={
              <Admin organizationID={currentOrgID}>
                <TitleBarButton
                  key={1}
                  label={i18n.t(`${packageNS}:tr000277`)}
                  icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
                  to={`/organizations/${currentOrgID}/applications/create`}
                />
              </Admin>
            }
          >
            <TitleBarTitle title={i18n.t(`${packageNS}:tr000076`)} />
          </TitleBar>
          <Row>
            <Col>
              <Card className="card-box shadow-sm">
                <CardBody className="position-relative">
                  {this.state.loading && <Loader />}
                  <AdvancedTable
                    data={this.state.data}
                    columns={getColumns(currentOrgID)}
                    keyField="id"
                    onTableChange={this.handleTableChange}
                    rowsPerPage={10}
                    totalSize={this.state.totalSize}
                    searchEnabled={false}
                  />
                </CardBody>
              </Card>
            </Col>
          </Row>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default ListApplications;
