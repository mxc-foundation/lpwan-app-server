import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Button, Collapse, NavLink } from 'reactstrap';
import i18n, { packageNS } from '../../../../i18n';




const CURRENT_CARD = "detailsCard";

class DetailsCard extends Component {
  render() {
    const { collapseCard, setCollapseCard, deviceProfile } = this.props;

    return(
      <Card>
        <Button color="secondary" onClick={() => setCollapseCard(CURRENT_CARD)}>
          <i className={`mdi mdi-arrow-${collapseCard[CURRENT_CARD] ? 'up' : 'down'}`}></i>
          &nbsp;&nbsp;
          <h5 style={{ color: "#fff", display: "inline" }}>
            {i18n.t(`${packageNS}:tr000280`)}
          </h5>
        </Button>
        <Collapse isOpen={collapseCard[CURRENT_CARD]}>
          <CardContent>
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                  <TableCell>{this.props.device.device.name}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000079`)}</TableCell>
                  <TableCell>{this.props.device.device.description}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000281`)}</TableCell>
                  <TableCell>
                    {
                      deviceProfile && deviceProfile.deviceProfile.id
                      ? 
                        (
                          <Button color="primary" style={{ margin: "0.5em" }}>
                            <NavLink
                              style={{ color: "#fff", padding: "0" }}
                              tag={Link}
                              to={`/organizations/${this.props.match.params.organizationID}/device-profiles/${deviceProfile.deviceProfile.id}`}
                            >
                              {deviceProfile.deviceProfile.name}
                            </NavLink>
                          </Button>
                        )
                      : (
                        <>
                          <Button color="primary" style={{ margin: "0.5em" }}>
                            <NavLink
                              style={{ color: "#fff", padding: "0" }}
                              tag={Link}
                              to={`/organizations/${this.props.match.params.organizationID}/devices/${this.props.match.params.devEUI}/edit`}
                            >
                              Associate with existing Device Profile
                            </NavLink>
                          </Button>
                          <Button color="primary" style={{ margin: "0.5em" }}>
                            <NavLink
                              style={{ color: "#fff", padding: "0" }}
                              tag={Link}
                              to={`/organizations/${this.props.match.params.organizationID}/devices-profiles/create`}
                            >
                              Create Device Profile
                            </NavLink>
                          </Button>
                        </>
                      )
                    }
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Collapse>
      </Card>
    );
  }
}

export default DetailsCard;
