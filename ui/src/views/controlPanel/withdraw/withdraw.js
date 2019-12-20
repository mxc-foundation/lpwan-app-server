import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Grid from "@material-ui/core/Grid";
import { withStyles } from "@material-ui/core/styles";
import styles from "../../withdraw/WithdrawStyle"
import theme from "../../../theme";
import TableCell from "@material-ui/core/TableCell";
import TitleBarTitle from "../../../components/TitleBarTitle";
import i18n, {packageNS} from "../../../i18n";


class SuperAdminWithdraw extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <Grid container spacing={24} className={this.props.classes.backgroundColor}>
            <Grid item xs={12} className={this.props.classes.divider}>
                <div className={this.props.classes.TitleBar}>
                    <TitleBarTitle title={i18n.t(`${packageNS}:menu.withdraw.withdraw`)} />
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