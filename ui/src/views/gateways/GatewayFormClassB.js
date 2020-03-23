import React, { Component } from "react";
import { Row, Col, Table } from "reactstrap";

import i18n, { packageNS } from "../../i18n";
import EditableTextInput from "../../components/EditableTextInput";

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
  onDataChanged(idx, field, changedValue) {
    let records = [...this.state.records];
    records[idx][field] = changedValue;
    if (this.props.onDataChanged) {
      this.props.onDataChanged(records);
    } else {
      this.setState({ records });
    }
  }

  render() {
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
                {this.state.records.length
                  ? this.state.records.map((record, idx) => {
                      return (
                        <tr key={idx}>
                          <td>
                            <EditableTextInput
                              value={record.beacon_period}
                              id={idx}
                              field={"beacon_period"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                          <td>
                            <EditableTextInput
                              value={record.beacon_freq}
                              id={idx}
                              field={"beacon_freq"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                          <td>
                            <EditableTextInput
                              value={record.beacon_datarate}
                              id={idx}
                              field={"beacon_datarate"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                          <td>
                            <EditableTextInput
                              value={record.beacon_bandwidth}
                              id={idx}
                              field={"beacon_bandwidth"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                          <td>
                            <EditableTextInput
                              value={record.beacon_power}
                              id={idx}
                              field={"beacon_power"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                          <td>
                            <EditableTextInput
                              value={record.beacon_info}
                              id={idx}
                              field={"beacon_info"}
                              onDataChanged={this.onDataChanged}
                            />
                          </td>
                        </tr>
                      );
                    })
                  : null}
              </tbody>
            </Table>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayFormClassB;
