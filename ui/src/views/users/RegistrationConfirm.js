import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Button } from 'reactstrap';
import { Formik, Form, Field, FieldArray } from 'formik';
import { ReactstrapInput } from '../../components/FormInputs';
import * as Yup from 'yup';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";
import classNames from "classnames";

import DropdownMenuLanguage from "../../components/DropdownMenuLanguage";
//import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import i18n, { packageNS } from '../../i18n';


const styles = {
    languageWrapper: {
        marginLeft: '15px'
    },
    textField: {
        width: "100%",
    },
    link: {
        "& a": {
            color: theme.palette.primary.main,
            textDecoration: "none",
        },
    },
};


class RegistrationConfirmForm extends FormComponent {
    constructor(props) {
        super(props);

        this.state = {}
    }
    componentDidMount() {
        SessionStore.confirmRegistration(this.props.securityToken, (resp) => {
            if (resp) {
                const object = this.state;
                object.object = resp;
                object.isTokenValid = true;
                this.setState({
                    object
                })
                SessionStore.setToken(resp.jwt)
            } else {
                const object = this.state;
                object.isTokenValid = false;
                this.setState({
                    object
                })
            }
        })
    }

    render() {
        const object = this.state;

        if (object === undefined) {
            return (<div></div>);
        }

        let fieldsSchema = {
            username: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
            password: Yup.string().trim().matches(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/g, i18n.t(`${packageNS}:menu.registration.password_hint`)).required(i18n.t(`${packageNS}:tr000431`)),
            passwordConfirmation: Yup.string().trim().matches(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/g, i18n.t(`${packageNS}:menu.registration.password_hint`)).required(i18n.t(`${packageNS}:tr000431`)),
            organizationName: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
            organizationDisplayName: Yup.string().required(i18n.t(`${packageNS}:tr000431`))
        }

        const formSchema = Yup.object().shape(fieldsSchema);
        
        let username = "";
        if(this.state.object !== undefined){
            username = this.state.object.object.username;
        }
        return (
            <Formik
                enableReinitialize
                initialValues={{
                    username: username || '',
                    password: object.password || '',
                    passwordConfirmation: object.passwordConfirmation || '',
                    organizationName: object.organizationName || '',
                    organizationDisplayName: object.organizationDisplayName || ''
                }}
                validationSchema={formSchema}
                onSubmit={this.props.onSubmit}>
                {({
                    handleSubmit,
                    setFieldValue,
                    handleChange,
                    handleBlur,
                    values
                }) => (
                        <Form onSubmit={handleSubmit} noValidate>
                            <Field
                                type="email"
                                label={i18n.t(`${packageNS}:tr000003`)+"*"}
                                name="username"
                                id="username"
                                value={values.username}
                                onChange={handleChange}
                                helpText={i18n.t(`${packageNS}:tr000150`)}
                                component={ReactstrapInput}
                                onBlur={handleBlur}
                                required
                            />

                            <Field
                                type="password"
                                label={i18n.t(`${packageNS}:tr000004`)+"*"}
                                name="password"
                                id="password"
                                value={values.password}
                                onChange={handleChange}
                                helpText={i18n.t(`${packageNS}:menu.registration.password_hint`)}
                                component={ReactstrapInput}
                                onBlur={handleBlur}
                                required
                            />

                            <Field
                                type="password"
                                label={i18n.t(`${packageNS}:tr000023`)+"*"}
                                name="passwordConfirmation"
                                id="passwordConfirmation"
                                value={values.passwordConfirmation}
                                onChange={handleChange}
                                helpText={i18n.t(`${packageNS}:menu.registration.password_hint`)}
                                component={ReactstrapInput}
                                onBlur={handleBlur}
                                required
                            />

                            <Field
                                type="text"
                                label={i18n.t(`${packageNS}:tr000030`)+"*"}
                                name="organizationName"
                                id="organizationName"
                                value={values.organizationName}
                                onChange={handleChange}
                                component={ReactstrapInput}
                                onBlur={handleBlur}
                                required
                            />

                            <Field
                                type="text"
                                label={i18n.t(`${packageNS}:tr000031`)+"*"}
                                name="organizationDisplayName"
                                id="organizationDisplayName"
                                value={values.organizationDisplayName}
                                onChange={handleChange}
                                component={ReactstrapInput}
                                onBlur={handleBlur}
                                required
                            />
                            <Button type="submit" color="primary">{this.props.submitLabel}</Button>
                        </Form>
                    )}
            </Formik>

        );
    }
}


class RegistrationConfirm extends Component {
    constructor() {
        super();

        this.state = {
            isTokenValid: null,
            isPwdMatch: null
        }

        localStorage.setItem('jwt', '')

        this.onSubmit = this.onSubmit.bind(this);
    }

    onChangeLanguage = e => {
        const newLanguage = {
            id: e.id,
            label: e.label,
            value: e.value,
            code: e.code
        }

        this.props.onChangeLanguage(newLanguage);
    }

    onSubmit(data) {
        console.log('data', data);
        if (data.password === data.passwordConfirmation) {
            this.setState({
                isPwdMatch: true
            })

            let request = {
                userId: data.id,
                password: data.password,
                organizationName: data.organizationName,
                organizationDisplayName: data.organizationDisplayName,
            }

            SessionStore.finishRegistration(request, (responseData) => {
                SessionStore.logout(() => {
                    this.props.history.push("/logout");
                });
            })
        } else {
            this.setState({
                isPwdMatch: false
            })
        }

    }

    render() {
        console.log(this.props);
        return (
            <Grid container justify="center">
                <Grid item xs={6} lg={4}>
                    <Card>
                        <div className={classNames(this.props.classes.languageWrapper)}>
                            <DropdownMenuLanguage onChangeLanguage={this.onChangeLanguage} />
                        </div>
                        <CardHeader
                            title={i18n.t(`${packageNS}:tr000019`)}
                        />
                        <CardContent>
                            {this.state.isPwdMatch !== null && this.state.isPwdMatch === false &&
                                <p style={{ color: 'Red', textAlign: 'center' }}>{i18n.t(`${packageNS}:tr000025`)}</p>
                            }
                            <RegistrationConfirmForm
                                submitLabel={i18n.t(`${packageNS}:tr000022`)}
                                onSubmit={this.onSubmit}
                                securityToken={this.props.match.params.securityToken}
                            />
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        );
    }
}

export default withStyles(styles)(withRouter(RegistrationConfirm));
