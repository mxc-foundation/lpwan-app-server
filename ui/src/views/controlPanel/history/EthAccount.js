import React, { Component } from "react";
import AdvancedTable from "../../../components/AdvancedTable";
import Loader from "../../../components/Loader";

import i18n, { packageNS } from '../../../i18n';
import HistoryStore from "../../../stores/HistoryStore";

import { ETHER } from "../../../util/CoinType"
import { SUPER_ADMIN } from "../../../util/M2mUtil";

const CreatedAtColumn = (cell, row, index, extraData) => {
  return row.createdAt.substring(0, 10);
}

const getColumns = () => (
  [{
    dataField: 'addr',
    text: i18n.t(`${packageNS}:menu.staking.account`),
    sort: false,
  }, {
    dataField: 'status',
    text: i18n.t(`${packageNS}:menu.staking.status`),
    sort: false,
  }, {
    dataField: 'createdAt',
    text: i18n.t(`${packageNS}:menu.staking.date`),
    formatter: CreatedAtColumn,
    sort: false,
  }]
);

class SuperNodeEthAccount extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
    this.state = {
      data: [],
      stats: {},
      totalSize: 0
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
    HistoryStore.getChangeMoneyAccountHistory(ETHER, SUPER_ADMIN, limit, offset, res => {
      const object = this.state;
      object.totalSize = res.count;
      object.data = res.changeHistory;
      object.loading = false;
      this.setState({object});
    }); 
  }

  componentDidMount() {
    this.getPage(10);
  }

  render() {
    return(
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable data={this.state.data} columns={getColumns()}
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} totalSize={this.state.totalSize} rowsPerPage={10}></AdvancedTable>
      </div>
    );
  }
}

export default SuperNodeEthAccount;
