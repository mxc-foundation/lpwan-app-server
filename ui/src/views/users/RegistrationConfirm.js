import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import classNames from "classnames";
import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Button } from 'reactstrap';
import * as Yup from 'yup';
import FormComponent from "../../classes/FormComponent";
import DropdownMenuLanguage from "../../components/DropdownMenuLanguage";
import { ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import VerificationWith2FA from "../common/VerificationWith2FA";
import theme from "../../theme";




class RegistrationConfirm extends Component {
    constructor() {
        super();
    }

    render() {
        let token = '';
        if(this.props !== undefined){
            token = this.props.match.params.securityToken;
        }
        
        return (
            <VerificationWith2FA restart={`/registration`} next={`/registration-confirm-steptwo/${token}`}/>
        );
    }
}

export default withRouter(RegistrationConfirm);
