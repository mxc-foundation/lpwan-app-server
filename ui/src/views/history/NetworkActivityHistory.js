import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import i18n, { packageNS } from '../../i18n';
import HistoryStore from "../../stores/HistoryStore";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import LinkVariant from "mdi-material-ui/LinkVariant";

const PckRcvColumn = (cell, row, index, extraData) => {
  return parseInt(row.DlCntGw - row.DlCntGwFree);
}

const getColumns = () => (
  [{
    dataField: 'StartAt',
    text: i18n.t(`${packageNS}:menu.staking.time`),
    sort: false,
  }, {
    dataField: 'DlCntDv',
    text: i18n.t(`${packageNS}:menu.staking.packets_sent`),
    sort: false,
  }, {
    dataField: 'DlCntDvFree',
    text: i18n.t(`${packageNS}:menu.staking.free_packets`),
    sort: false,
  }, {
    dataField: 'Packets Received',
    text: i18n.t(`${packageNS}:menu.staking.packets_received`),
    formatter: PckRcvColumn,
    sort: false,
  }, {
    dataField: 'Income',
    text: i18n.t(`${packageNS}:menu.staking.earned`),
    sort: false,
  }, {
    dataField: 'Spend',
    text: i18n.t(`${packageNS}:menu.staking.spent`),
    sort: false,
  },{
    dataField: 'UpdatedBalance',
    text: i18n.t(`${packageNS}:menu.staking.balance`),
    sort: false,
  }]
);

class NetworkActivityHistory extends Component {
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
    HistoryStore.getWalletUsageHist(this.props.organizationID, offset, limit, data => {
      this.setState({ data: data.walletUsageHis, loading: false });
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
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} rowsPerPage={10}></AdvancedTable>
      </div>
    );
  }
}

export default withRouter(NetworkActivityHistory);