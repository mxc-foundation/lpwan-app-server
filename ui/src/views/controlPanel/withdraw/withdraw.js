import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Grid from "@material-ui/core/Grid";
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "../../withdraw/WithdrawStyle"
import theme from "../../../theme";
import TableCell from "@material-ui/core/TableCell";
import i18n, {packageNS} from "../../../i18n";

import breadcrumbStyles from "../../common/BreadcrumbStyles";

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class SuperAdminWithdraw extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        const { classes } = this.props;

        return (
            <Grid container spacing={24} className={this.props.classes.backgroundColor}>
                <Grid item xs={12} className={this.props.classes.divider}>
                    <div className={this.props.classes.TitleBar}>
                        <Breadcrumb className={classes.breadcrumb}>
                            <BreadcrumbItem>
                                <Link
                                    className={classes.breadcrumbItemLink}
                                    to={`/organizations`}
                                    onClick={() => {
                                        // Change the sidebar content
                                        this.props.switchToSidebarId('DEFAULT');
                                    }}
                                >
                                    Control Panel
                                </Link>
                            </BreadcrumbItem>
                            <BreadcrumbItem className={classes.breadcrumbItem}>Wallet</BreadcrumbItem>
                            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.withdraw.withdraw`)}</BreadcrumbItem>
                        </Breadcrumb>    
                    </div>
                </Grid>

                <Grid item xs={6}>
                    <TableCell align={this.props.align}>
                        <span style={
                            {
                                textDecoration: "none",
                                color: theme.palette.primary.main,
                                cursor: "pointer",
                                padding: 0,
                                fontWeight: "bold",
                                fontSize: 20,
                                opacity: 0.7,
                                "&:hover": {
                                    opacity: 1,
                                },
                                margin: 16
                            }
                        } className={this.props.classes.link} >
                            {i18n.t(`${packageNS}:menu.messages.coming_soon`)}
                        </span>
                    </TableCell>
                </Grid>
            </Grid>
        );
    }
}

export default withStyles(styles)(withRouter(SuperAdminWithdraw));