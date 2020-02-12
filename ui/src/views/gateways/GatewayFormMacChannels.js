import React, { Component } from "react";
import { Row, Col, Table, CustomInput } from "reactstrap";

import i18n, { packageNS } from "../../i18n";

class GatewayFormMacChannels extends Component {
  constructor(props) {
    super(props);

    this.state = {
      records: props.records
    };

    this.onToggle = this.onToggle.bind(this);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.records !== this.props.records) {
      this.setState({ records: this.props.records });
    }
  }

  /**
   * On switch toggle
   * @param {*} idx
   * @param {*} e
   */
  onToggle(idx, e) {
    let records = [...this.state.records];
    records[idx]["enable"] = e.target.checked;
    this.setState({ records });

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
            <h5>{i18n.t(`${packageNS}:tr000599`)}</h5>
            <Table>
              <thead>
                <tr>
                  <th width="33%">{i18n.t(`${packageNS}:tr000610`)}</th>
                  <th width="33%">{i18n.t(`${packageNS}:tr000611`)}</th>
                  <th width="33%">{i18n.t(`${packageNS}:tr000612`)}</th>
                </tr>
              </thead>
              <tbody>
                {this.state.records.map((record, idx) => {
                  return (
                    <tr key={idx}>
                      <td>{record.channel}</td>
                      <td>{record.freq_hz ? record.freq_hz / 1000000 : 0}</td>
                      <td>
                        <CustomInput
                          type="switch"
                          id={`switch-${idx}`}
                          name="channel-enable"
                          label=""
                          checked={record.enable}
                          onChange={e => this.onToggle(idx, e)}
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </Table>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayFormMacChannels;
