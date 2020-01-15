import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";

import DeviceStore from "../../../stores/DeviceStore";
import DetailsCard from "./subComponents/DetailsCard";
import EnqueueCard from "./subComponents/EnqueueCard";
import QueueCard from "./subComponents/QueueCard";
import StatusCard from "./subComponents/StatusCard";

class DeviceDetails extends Component {
  constructor() {
    super();
    this.state = {
      collapseCard: {
        detailsCard: true,
        statusCard: true,
        enqueueCard: true,
        queueCard: true
      },
      activated: false,
    };
  }

  componentDidMount() {
    this.setDeviceActivation();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.device !== this.props.device) {
      this.setDeviceActivation();
    }
  }

  setCollapseCard = (collapseCardComponentName) => {
    const { collapseCard } = this.state;
    const existingValue = collapseCard[collapseCardComponentName];
    const newObj = {};
    newObj[collapseCardComponentName] = !existingValue;
    let newCollapseCard = Object.assign(collapseCard, newObj);
    this.setState({
      collapseCard: newCollapseCard
    })
  }

  setDeviceActivation = () => {
    if (this.props.device === undefined) {
      return;
    }

    DeviceStore.getActivation(this.props.device.device.devEUI, resp => {
      if (resp === null) {
        this.setState({
          activated: false,
        });
      } else {
        this.setState({
          activated: true,
        });
      }
    });
  };

  render() {
    const { collapseCard } = this.state;

    return(
      <Grid container spacing={4} style={{ backgroundColor: "#ebeff2", marginTop: "10px" }}>
        <Grid item xs={12} md={6}>
          <DetailsCard
            collapseCard={collapseCard}
            device={this.props.device}
            deviceProfile={this.props.deviceProfile}
            match={this.props.match}
            setCollapseCard={this.setCollapseCard}
          />
        </Grid>
        <Grid item xs={12} md={6}>
          <StatusCard
            collapseCard={collapseCard}
            device={this.props.device}
            setCollapseCard={this.setCollapseCard}
          />
        </Grid>
        {this.state.activated &&
          <Grid item xs={12} lg={6}>
            <EnqueueCard
              collapseCard={collapseCard}
              setCollapseCard={this.setCollapseCard}
            />
          </Grid>
        }
        {this.state.activated &&
          <Grid item xs={12} lg={6}>
            <QueueCard
              collapseCard={collapseCard}
              setCollapseCard={this.setCollapseCard}
            />
          </Grid>
        }
      </Grid>
    );
  }
}

export default DeviceDetails;
