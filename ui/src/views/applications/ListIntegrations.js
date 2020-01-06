import React, { Component } from "react";
import { Link } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import ApplicationStore from "../../stores/ApplicationStore";
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

const ApplicationKindColumn = (cell, row, index, extraData) => {
  const currentOrgID = extraData['currentOrgID'];
  const applicationId = extraData['applicationId'];
  const kind = row.kind.toLowerCase();

  return <Link to={`/organizations/${currentOrgID}/applications/${applicationId}/integrations/${kind}`}>{row.kind}</Link>;
}

const getColumns = (currentOrgID, applicationId) => (
  [
    {
      dataField: 'kind',
      text: i18n.t(`${packageNS}:tr000412`),
      sort: false,
      formatter: ApplicationKindColumn,
      formatExtraData: {
        currentOrgID: currentOrgID,
        applicationId: applicationId
      }
    }
  ]
);


class ListIntegrations extends Component {
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
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    this.setState({ loading: true });

    ApplicationStore.listIntegrations("", this.props.match.params.applicationID, currentOrgID, limit, offset, (res) => {
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
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    return(
      <React.Fragment>
        <TitleBar
          buttons={
            <TitleBarButton
              key={1}
              label={i18n.t(`${packageNS}:tr000277`)}
              icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
              to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/integrations/create`}
            />
          }
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000553`)} />
        </TitleBar>

        <div className="position-relative">
          {this.state.loading && <Loader />}
          <AdvancedTable
            data={this.state.data}
            columns={getColumns(currentOrgID, currentApplicationID)}
            keyField="kind"
            onTableChange={this.handleTableChange}
            rowsPerPage={10}
            searchEnabled={true}
          />
        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ListIntegrations);
