import React, { Component } from "react";
import { Col, Row, Table } from "reactstrap";
import EditableTextInput from "../../components/EditableTextInput";
import i18n, { packageNS } from "../../i18n";


class GatewayFormClassB extends Component {
  constructor(props) {
    super(props);

    this.state = {
      records: props.records
    };

    this.onDataChanged = this.onDataChanged.bind(this);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.records !== this.props.records) {
      this.setState({ records: this.props.records });
    }
  }

  /**
   * On data edited
   * @param {*} object
   * @param {*} field
   * @param {*} changedValue
   */
  onDataChanged(e) {
    let records = this.state.records;
    records[e.id] = e.value;

    if (this.props.onDataChanged) {
      this.props.onDataChanged(records);
    } else {
      this.setState({ records });
    }
  }

  render() {
    if (this.state.records.beacon_period === undefined) {
      return <div></div>;
    }

    return (
      <React.Fragment>
        <Row>
          <Col>
            <h5>{i18n.t(`${packageNS}:tr000601`)}</h5>
            <Table>
              <thead>
                <tr>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000604`)}</th>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000605`)}</th>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000606`)}</th>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000607`)}</th>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000608`)}</th>
                  <th width="16.33%">{i18n.t(`${packageNS}:tr000609`)}</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_period}
                      id={'beacon_period'}
                      field={"beacon_period"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_freq_hz}
                      id={'beacon_freq_hz'}
                      field={"beacon_freq_hz"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_datarate}
                      id={'beacon_datarate'}
                      field={"beacon_datarate"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_bw_hz}
                      id={'beacon_bw_hz'}
                      field={"beacon_bw_hz"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_power}
                      id={'beacon_power'}
                      field={"beacon_power"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                  <td>
                    <EditableTextInput
                      value={this.state.records.beacon_infodesc}
                      id={'beacon_infodesc'}
                      field={"beacon_infodesc"}
                      onDataChanged={this.onDataChanged}
                    />
                  </td>
                </tr>

                {/* {this.state.records.length
                  ? this.state.records.map((record, idx) => {
                    console.log('record', record);
                      return ();
                    })
                  : null} */}
              </tbody>
            </Table>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayFormClassB;
