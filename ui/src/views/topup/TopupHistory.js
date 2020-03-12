import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Row, Col, Card, CardBody } from 'reactstrap';

import Loader from "../../components/Loader";
import AdvancedTable from "../../components/AdvancedTable";
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TopupStore from "../../stores/TopupStore";
import i18n, { packageNS } from '../../i18n';


const tableCols = [
  {
    dataField: 'amount',
    text: i18n.t(`${packageNS}:menu.topup.history.amount`),
    sort: false,
    formatter: (cell, row, rowIndex, formatExtraData) => {
      return <React.Fragment>{cell} MXC</React.Fragment>
    },
  },
  {
    dataField: 'lastUpdateTime',
    text: i18n.t(`${packageNS}:menu.topup.history.date`),
    sort: false,
  },
  {
    dataField: 'txHash',
    text: i18n.t(`${packageNS}:menu.topup.history.tx_hash`),
    sort: false,
  }
]



class TopupHistory extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
      data: [],
      totalSize: 0
    };
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage;
    
    /* let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    } */

    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  };


  getPage(offset) {
    this.setState({ loading: true });
    TopupStore.getTopUpHistory(this.props.organizationID, offset, MAX_DATA_LIMIT, res => {
      const object = this.state;

      object.totalSize = Number(res.count);
      object.data = res.topupHistory;
      object.loading = false;
      this.setState({ object });

      console.log(this.state);
    }, error => {
      this.setState({ loading: false });
    });
  }

  render() {


    return (<React.Fragment>
      <Row>
        <Col>
          <Card>
            <CardBody className="pb-0">
              <div className="position-relative">
                {this.state.loading ? <Loader /> : null}
                <AdvancedTable data={this.state.data} columns={tableCols} keyField="id" totalSize={this.state.totalSize}
                  onTableChange={this.handleTableChange} searchEnabled={true}></AdvancedTable>
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(TopupHistory);
