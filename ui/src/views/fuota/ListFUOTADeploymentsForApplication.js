import React, { Component } from "react";
import { Link } from "react-router-dom";

import moment from "moment";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import Admin from "../../components/Admin";
import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";
import theme from "../../theme";


const styles = {
  buttons: {
    textAlign: "right",
  },
  button: {
    marginLeft: 2 * theme.spacing(1),
  },
  icon: {
    marginRight: theme.spacing(1),
  },
};

const FUOTADeploymentNameColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const applicationId = extraData['applicationId'];
  return <Link to={
    applicationId
    ? `/organizations/${currentOrgID}/applications/${applicationId}/fuota-deployments/${row.id}`
    : `/organizations/${currentOrgID}/fuota-deployments/${row.id}`
  }>{row.name}</Link>;
}

const FUOTADeploymentCreatedAtColumn = (cell, row, index) => {
  return moment(row.createdAt).format('lll');
}

const FUOTADeploymentUpdatedAtColumn = (cell, row, index) => {
  return moment(row.updatedAt).format('lll');
}

const FUOTADeploymentNextStepAfterColumn = (cell, row, index) => {
  return moment(row.nextStepAfter).format('lll');
}

const getColumns = (currentOrgID, applicationId) => (
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
      formatter: FUOTADeploymentNameColumn,
      formatExtraData: {
        currentOrgID: currentOrgID,
        applicationId: applicationId
      }
    }, {
      dataField: 'createdAt',
      text: i18n.t(`${packageNS}:tr000321`),
      sort: false,
      formatter: FUOTADeploymentCreatedAtColumn,
    }, {
      dataField: 'updatedAt',
      text: i18n.t(`${packageNS}:tr000322`),
      sort: false,
      formatter: FUOTADeploymentUpdatedAtColumn,
    }, {
      dataField: 'state',
      text: "State",
      sort: false,
    }, {
      dataField: 'nextStepAfter',
      text: "Next Step After",
      sort: false,
      formatter: FUOTADeploymentNextStepAfterColumn,
    }
  ]
);


class ListFUOTADeploymentsForApplication extends Component {
  constructor() {
    super();

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
    this.setState({ loading: true });

    FUOTADeploymentStore.list({
      applicationID: this.props.match.params.applicationID,
      limit: limit,
      offset: offset,
    }, (res) => {
      this.setState({
        data: res.result,
        loading: false
      });
    });
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    return(
      <React.Fragment>
        <TitleBar
          // buttons={
          //   <Admin organizationID={currentOrgID}>
          //     <TitleBarButton
          //       key={1}
          //       label={i18n.t(`${packageNS}:tr000277`)}
          //       icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
          //       to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/1/fuota-deployments/create`}
          //       // FIXME - do not hard-code /1/ above, it should be the below...
          //       // :devEUI([\w]{16})
          //     />
          //   </Admin>
          // }
        >
          <TitleBarTitle title="FUOTA Deployments" />
        </TitleBar>
        <div className="position-relative">
          {this.state.loading && <Loader />}
          <AdvancedTable
            data={this.state.data}
            columns={getColumns(currentOrgID, currentApplicationID)}
            keyField="id"
            onTableChange={this.handleTableChange}
            rowsPerPage={10}
            searchEnabled={false}
          />
        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ListFUOTADeploymentsForApplication);

