import React, { Component } from "react";

import { Collapse, Button } from 'reactstrap';
import Grid from '@material-ui/core/Grid';
import Card from "@material-ui/core/Card";
import CardHeader from "@material-ui/core/CardHeader";
import CardContent from "@material-ui/core/CardContent";
import LinearProgress from '@material-ui/core/LinearProgress';
import Typography from "@material-ui/core/Typography";

import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';

import moment from "moment";

import i18n, { packageNS } from '../../i18n';
import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";


class FUOTADeploymentDetails extends Component {
  constructor() {
    super();

    this.state = {
      collapseCard1: true,
      collapseCard2: true,
      progress: 0,
      stepProgress: 0,
      lastReload: 0,
    };
  }

  componentDidMount() {
    this.timer = setInterval(this.progress, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }

    this.progress();
  }

  setCollapseCard1 = () => {
    const { collapseCard1 } = this.state;
    this.setState({
      collapseCard1: !collapseCard1
    })
  }

  setCollapseCard2 = () => {
    const { collapseCard2 } = this.state;
    this.setState({
      collapseCard2: !collapseCard2
    })
  }

  progress = () => {
    const { fuotaDeployment } = this.props;
    const start = moment(fuotaDeployment.updatedAt).unix();
    const nextStepAfter = moment(fuotaDeployment.fuotaDeployment.nextStepAfter).unix();
    const now = moment().unix();

    const span = nextStepAfter - start;
    const stepProgress = now - start;
    const progressScaled = stepProgress / span * 100;

    if (progressScaled > 100) {
      if (moment().unix() - this.state.lastReload > 5) {
        this.setState({
          lastReload: moment().unix(),
        }, FUOTADeploymentStore.emitReload());
      }
    }

    const states = 8;
    let state = 0;

    switch(fuotaDeployment.fuotaDeployment.state) {
      case "MC_CREATE":
        state = 0;
        break;
      case "MC_SETUP":
        state = 1;
        break;
      case "FRAG_SESS_SETUP":
        state = 2;
        break;
      case "MC_SESS_C_SETUP":
        state = 3;
        break;
      case "ENQUEUE":
        state = 4;
        break;
      case "STATUS_REQUEST":
        state = 5;
        break;
      case "SET_DEVICE_STATUS":
        state = 6;
        break;
      case "CLEANUP":
        state = 7;
        break;
      case "DONE":
        state = 8;
        break;
      default:
        state = 0;
        break;
    }

    this.setState({
      stepProgress: progressScaled,
      progress: state / states * 100,
    });
  }

  render() {
    const { collapseCard1, collapseCard2 } = this.state;
    const { fuotaDeployment } = this.props;
    let multicastTimeout = 0;
    if (fuotaDeployment.fuotaDeployment.groupType === "CLASS_C") {
      multicastTimeout = (1 << fuotaDeployment.fuotaDeployment.multicastTimeout);
    }

    const createdAt = moment(fuotaDeployment.createdAt).format('lll');
    const updatedAt = moment(fuotaDeployment.updatedAt).format('lll');
    const nextStepAfter = moment(fuotaDeployment.fuotaDeployment.nextStepAfter).format('lll');

    return(
      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <Card>
            <Button color="secondary" onClick={this.setCollapseCard1}>
              <i className={`mdi mdi-arrow-${collapseCard2 ? 'up' : 'down'}`}></i>
              &nbsp;&nbsp;
              <h5 style={{ color: "#fff", display: "inline" }}>
                {i18n.t(`${packageNS}:tr000280`)}
              </h5>
            </Button>
            <Collapse isOpen={collapseCard1}>
              <CardContent>
                <Table>
                  <TableBody>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000320`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.name}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000344`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.redundancy}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000345`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.unicastTimeout}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000346`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.dr}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000347`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.frequency}Hz</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000348`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.groupType}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000349`)}</TableCell>
                      <TableCell>{multicastTimeout} {i18n.t(`${packageNS}:tr000357`)}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </CardContent>
            </Collapse>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <Button color="secondary" onClick={this.setCollapseCard2}>
              <i className={`mdi mdi-arrow-${collapseCard2 ? 'up' : 'down'}`}></i>
              &nbsp;&nbsp;
              <h5 style={{ color: "#fff", display: "inline" }}>
                {i18n.t(`${packageNS}:tr000282`)}
              </h5>
            </Button>
            <Collapse isOpen={collapseCard2}>
              <CardContent>
                <Table>
                  <TableBody>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000321`)}</TableCell>
                      <TableCell>{createdAt}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000322`)}</TableCell>
                      <TableCell>{updatedAt}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000350`)}</TableCell>
                      <TableCell>{fuotaDeployment.fuotaDeployment.state}</TableCell>
                    </TableRow>
                    {fuotaDeployment.fuotaDeployment.state !== "DONE" && <TableRow>
                      <TableCell>{i18n.t(`${packageNS}:tr000351`)}</TableCell>
                      <TableCell>{nextStepAfter}</TableCell>
                    </TableRow>}
                  </TableBody>
                </Table>
              </CardContent>
              {fuotaDeployment.fuotaDeployment.state !== "DONE" && <CardContent>
                <Typography variant="subtitle2" gutterBottom>
                  {i18n.t(`${packageNS}:tr000352`)}
                </Typography>
                <LinearProgress variant="determinate" value={this.state.progress} />
              </CardContent>}
              {fuotaDeployment.fuotaDeployment.state !== "DONE" && <CardContent>
                <Typography variant="subtitle2" gutterBottom>
                  {i18n.t(`${packageNS}:tr000353`)}
                </Typography>
                  <LinearProgress variant="determinate" value={this.state.stepProgress} />
              </CardContent>}
            </Collapse>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default FUOTADeploymentDetails;

