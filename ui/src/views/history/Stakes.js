import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";

import i18n, { packageNS } from '../../i18n';
import StakeStore from "../../stores/StakeStore";

const StartColumn = (cell, row, index, extraData) => {
  return row.start.substring(0, 10);
}

const EndColumn = (cell, row, index, extraData) => {
  return row.end.substring(0, 10);
}

const getColumns = () => (
  [{
    dataField: 'stakeAmount',
    text: i18n.t(`${packageNS}:menu.staking.stake_amount`),
    sort: false,
  }, {
    dataField: 'start',
    text: i18n.t(`${packageNS}:menu.staking.start`),
    formatter: StartColumn,
    sort: false,
  }, {
    dataField: 'end',
    text: i18n.t(`${packageNS}:menu.staking.end`),
    formatter: EndColumn,
    sort: false,
  }, {
    dataField: 'revMonth',
    text: i18n.t(`${packageNS}:menu.staking.revenue_month`),
    sort: false,
  }, {
    dataField: 'networkIncome',
    text: i18n.t(`${packageNS}:menu.staking.network_income`),
    sort: false,
  }, {
    dataField: 'monthlyRate',
    text: i18n.t(`${packageNS}:menu.staking.monthly_rate`),
    sort: false,
  }, {
    dataField: 'revenue',
    text: i18n.t(`${packageNS}:menu.staking.revenue`),
    sort: false,
  }, {
    dataField: 'balance',
    text: i18n.t(`${packageNS}:menu.staking.balance`),
    sort: false,
  }]
);

class Stakes extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      stats: {}
    }
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
    const offset = (page - 1) * sizePerPage + 1;

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
  getPage = (limit, offset) => {
    this.setState({ loading: true });
    StakeStore.getStakingHistory(this.props.organizationID, offset, limit, data => {
      this.setState({ data: data.stakingHist, loading: false });
    });
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    return (
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable data={this.state.data} columns={getColumns()}
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} rowsPerPage={10}></AdvancedTable>
      </div>
    );
  }
}

export default withRouter(Stakes);