import React, {Component} from "react";
import {Link, withRouter} from 'react-router-dom';
import {Breadcrumb, BreadcrumbItem, Button, Card, CardBody, Col, Row} from 'reactstrap';
import TitleBar from "../../../components/TitleBar";
import i18n, {packageNS} from '../../../i18n';
import WithdrawStore from "../../../stores/WithdrawStore";
import * as Yup from "yup";
import {Field, Form, Formik} from "formik";
import {AsyncAutoComplete, ReactstrapInput} from "../../../components/FormInputs";

class UpdateWithdrawFeeForm extends Component {

    constructor(props) {
        super(props);

        this.state = {
            loading: true,
            currency: "ETH_MXC",
        };
    }

    componentDidMount() {
        this.loadSettings();
    }

    loadSettings = async () => {
        try {
            WithdrawStore.getWithdrawFee(this.state.currency, (resp) => {
                this.setState({withdrawFee: resp.withdrawFee});
            }, () => {
            });

        } catch (e) {
            console.log("Error", e)
        }
    };

    saveSettings = async (data) => {
        try {
            WithdrawStore.setWithdrawFee({
                currency: data.currency,
                withdrawFee: data.withdrawFee.toString(),
                password: data.password || "",
            }, (resp) => {
            });
        } catch (e) {
            console.log("Error", e)
        }
    };

    reset = () => {
        this.loadSettings();
    }

    handleChange = (name, event) => {
        this.setState({
            [name]: event.target.value
        });
    };

    getCurrencyOptions = (search, callbackFunc) => {
        this.setState({loading: true});
        const res = {
            "result": [
                {"name": "ETH_MXC"},
            ]
        };
        const options = res.result.map((cur) => {
            return {label: cur.name}
        });
        this.setState({loading: false});
        callbackFunc(options);
    }

    onCurrencySelect = (v) => {
        console.log("2", v.label)
        try {
            WithdrawStore.getWithdrawFee(v.label, (resp) => {
                this.setState({currency: v.label, withdrawFee: resp.withdrawFee});
            }, () => {
                this.setState({currency: v.label, withdrawFee: 0.0});
                console.log("2.2")
            });
        } catch (e) {
            console.log("Error", e)
        }
    }

    render() {
        let fieldsSchema = {
            withdrawFee: Yup.number().test(
                'is-decimal',
                'invalid decimal',
                value => (value + "").match(/^\d*\.?\d+$/),
            ),
            currency: Yup.string().trim(),
        }

        const formSchema = Yup.object().shape(fieldsSchema);

        return (
            <React.Fragment>
                <Formik
                    enableReinitialize
                    initialValues={
                        {
                            currency: this.state.currency,
                            withdrawFee: this.state.withdrawFee,
                        }
                    }
                    validationSchema={formSchema}
                    onSubmit={(values) => {
                        const castValues = formSchema.cast(values);
                        this.saveSettings({...castValues})
                    }}>
                    {props => {
                        const {
                            handleSubmit,
                            handleChange,
                            setFieldValue,
                            values,
                            handleBlur,
                        } = props;
                        // errors && console.error('validation errors', errors);
                        return (
                            <Form onSubmit={handleSubmit}>
                                <Field
                                    type="text"
                                    label={i18n.t(`${packageNS}:menu.settings.currency`)}
                                    name="currency"
                                    id="currency"
                                    value={values.currency || ""}
                                    getOption={(search, callbackFunc) => {
                                        callbackFunc({label: this.state.currency})
                                    }}
                                    getOptions={this.getCurrencyOptions}
                                    onChange={this.onCurrencySelect}
                                    setFieldValue={() => {
                                    }}
                                    inputProps={{
                                        clearable: true,
                                        cache: false,
                                    }}
                                    component={AsyncAutoComplete}
                                />

                                <Field
                                    type="text"
                                    label={i18n.t(`${packageNS}:menu.settings.withdraw_fee`)}
                                    name="withdrawFee"
                                    id="withdrawFee"
                                    value={values.withdrawFee || ""}
                                    component={ReactstrapInput}
                                    validate
                                    inputProps={{
                                        clearable: true,
                                        cache: false,
                                    }}
                                />

                                <Button type="submit" className="btn-block"
                                        color="primary">{this.props.submitLabel || i18n.t(`${packageNS}:tr000066`)}</Button>
                            </Form>
                        );
                    }}
                </Formik>
            </React.Fragment>
        );
    }
}

class UpdateWithdrawFee extends Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {

        return (
            <React.Fragment>
                <TitleBar>
                    <Breadcrumb>
                        <BreadcrumbItem>
                            <Link
                                to={`/organizations`}
                                onClick={() => {
                                    // Change the sidebar content
                                    this.props.switchToSidebarId('DEFAULT');
                                }}
                            >
                                {i18n.t(`${packageNS}:menu.control_panel`)}
                            </Link>
                        </BreadcrumbItem>
                        <BreadcrumbItem>{i18n.t(`${packageNS}:tr000451`)}</BreadcrumbItem>
                        <BreadcrumbItem
                            active>{i18n.t(`${packageNS}:menu.settings.update_withdraw_fee`)}</BreadcrumbItem>
                    </Breadcrumb>
                </TitleBar>
                <Row>
                    <Col>
                        <Card>
                            <CardBody>
                                <UpdateWithdrawFeeForm
                                    submitLabel={i18n.t(`${packageNS}:tr000066`)}
                                />
                            </CardBody>
                        </Card>
                    </Col>
                </Row>
            </React.Fragment>
        );
    }
}

export default withRouter(UpdateWithdrawFee);
