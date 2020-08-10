import React, { Component, useState, useEffect } from 'react';
import classNames from 'classnames';
import { withRouter } from "react-router-dom";
import AsyncSelect from 'react-select/async';
import confirm from "reactstrap-confirm";
import GatewayStore from "../stores/GatewayStore";
import { CustomInput, FormFeedback, FormGroup, Row, FormText, Input, InputGroup, InputGroupAddon, Label, Alert, Button } from "reactstrap";
import i18n, { packageNS } from '../i18n';

const ReactstrapInput = (
    {
        field: { ...fields },
        form: { touched, errors, ...rest },
        helpText,
        ...props
    }) => (<FormGroup>
        <Label for={props.id}>{props.label}</Label>
        <Input {...props} {...fields} invalid={Boolean(touched[fields.name] && errors[fields.name])} />
        {touched[fields.name] && errors[fields.name] ? <FormFeedback>{errors[fields.name]}</FormFeedback> : null}
        {helpText && <FormText color="muted">{helpText}</FormText>}
    </FormGroup>
    );

const ReactstrapCheckbox = ({
    field,
    form: { isSubmitting, setFieldValue, touched, errors, values },
    helpText,
    ...props
}) => {
    return (
        <FormGroup >
            <CustomInput {...props} type="checkbox" checked={values[field.name]} value={field.value}
                /* onChange={(event, value) => setFieldValue(field.name, event.checked)} */ onChange={props.onChange} />
            {helpText && <FormText color="muted">{helpText}</FormText>}
        </FormGroup>
    )
};

const ReactstrapRadio = ({
    field,
    form: { isSubmitting, setFieldValue, touched, errors, values },
    disabled = false,
    ...props
}) => {
    return (
        <FormGroup check inline>
            <Label for={props.id}>
                <Input {...props} type="radio" name={field.name} checked={values[field.name] === field.value} value={field.value}
                    onChange={(event, value) => setFieldValue(field.name, field.value)} />{props.label}
            </Label>
        </FormGroup>
    )
};


const ReactstrapSelect = ({
    field,
    form: { isSubmitting, touched, errors },
    disabled = false,
    ...props
}) => {
    let error = errors[field.name];
    let touch = touched[field.name];
    return (
        <FormGroup>
            <Label for={props.inputprops.id}>{props.label}</Label>
            <Input id={props.inputprops.id} {...field} {...props} type="select"
                invalid={Boolean(touched[field.name] && errors[field.name])}
                placeholder="Test">
                <option value="">{props.inputprops.defaultOption}</option>
                {props.inputprops.options.map((option, index) => {
                    if (option.name)
                        return (<option value={option.id} key={index}>{option.name}</option>);
                    return (<option value={option} key={index}>{option}</option>)
                })}
            </Input>
            {touch && error && <FormFeedback>{error}</FormFeedback>}
        </FormGroup>
    )
};


const ReactstrapInputGroup = (
    {
        field: { ...fields },
        form: { touched, errors, ...rest },
        ...props
    }) => {
    const InputRenderControl = props.inputComponent || Input;
    return <FormGroup>
        <Label for={props.id}>{props.label}</Label>

        <InputGroup>
            {props.prepend && <InputGroupAddon addonType="prepend">
                {props.prepend}
            </InputGroupAddon>}

            <InputRenderControl {...props} {...fields} onChange={(e) => {
                if (props.onChange) props.onChange(e);
                if (fields.onChange) fields.onChange(e);
            }} invalid={Boolean(touched[fields.name] && errors[fields.name])}
                classes={classNames({ 'is-invalid': Boolean(touched[fields.name] && errors[fields.name]) })} />

            {touched[fields.name] && errors[fields.name] ? <FormFeedback className="order-last">{errors[fields.name]}</FormFeedback> : ''}

            {props.append && <InputGroupAddon addonType="append">
                {props.append}
            </InputGroupAddon>}
        </InputGroup>

        {props.helpText && <FormText color="muted">{props.helpText}</FormText>}

    </FormGroup>
};


class AutocompleteSelect extends Component {
    constructor(props) {
        super(props);

        this.state = {
            options: [],
        };

        this.setInitialOptions = this.setInitialOptions.bind(this);
        this.setSelectedOption = this.setSelectedOption.bind(this);
        this.onAutocomplete = this.onAutocomplete.bind(this);
    }

    componentDidMount() {
        this.setInitialOptions(this.setSelectedOption);
    }

    componentDidUpdate(prevProps) {
        if (prevProps.value === this.props.value && prevProps.triggerReload === this.props.triggerReload) {
            return;
        }
        this.setInitialOptions(this.setSelectedOption);
    }

    setInitialOptions(callbackFunc) {
        this.props.getOptions("", options => {
            this.setState({
                options: options,
            }, callbackFunc);
        });
    }

    setSelectedOption() {
        if (this.props.getOption !== undefined) {
            if (this.props.value !== undefined && this.props.value !== "" && this.props.value !== null) {
                this.props.getOption(this.props.value, resp => {
                    this.setState({
                        selectedOption: resp,
                    });
                });
            } else {
                this.setState({
                    selectedOption: "",
                });

                if (!this.props.noFirstItemSelected) {
                    // If there are any organizations listed, then choose the first one by default
                    this.props.getOptions("", options => {
                        if (options.length > 0) {
                            this.setState({
                                selectedOption: options[0],
                            });
                        }
                    });
                }
            }
        } else {
            if (this.props.value !== undefined && this.props.value !== "" && this.props.value !== null) {
                for (const opt of this.state.options) {
                    if (this.props.value === opt.value) {
                        this.setState({
                            selectedOption: opt,
                        });
                    }
                }
            } else {
                this.setState({
                    selectedOption: "",
                });
            }
        }
    }

    onAutocomplete(input) {
        return new Promise((resolve, reject) => {
            this.props.getOptions(input, options => {

                this.setState({
                    options: options,
                });

                resolve(options);
            });
        });
    }

    render() {
        const { field, setFieldValue, ...props } = this.props;
        const inputProps = this.props.inputProps || {};

        return (

            <FormGroup>
                <Label for={props.id}>{props.label}</Label>
                <AsyncSelect
                    {...field}
                    {...props}
                    {...inputProps}
                    instanceId={props.id}
                    clearable={false}
                    defaultOptions={this.state.options}
                    loadOptions={this.onAutocomplete}
                    value={this.state.selectedOption || ""}
                    onChange={(option) => { this.setState({ selectedOption: option }); setFieldValue(field.name, option.value); if(props.onChange) props.onChange(option); }}
                />
                {props.helpText && <FormText color="muted">{props.helpText}</FormText>}
            </FormGroup>
        );
    }
}

const AsyncAutoComplete = withRouter(AutocompleteSelect);


const ReactstrapPasswordInput = ({
    field: { ...fields },
    form: { touched, errors, setFieldTouched, ...rest },
    helpText = false,
    ...props
}) => {
    const [values, setValues] = React.useState({
        password: fields.value ? fields.value : '',
        showPassword: false,
    });

    const handleChange = prop => event => {
        setValues({ ...values, [prop]: event.target.value });
        if (props.onChange) props.onChange(event.target.valu);
        if (fields.onChange) fields.onChange(event);
    };

    const handleClickShowPassword = () => {
        setValues({ ...values, showPassword: !values.showPassword });
    };

    const handleMouseDownPassword = event => {
        event.preventDefault();
    };

    return (
        <FormGroup>
            <Label for={props.id}>{props.label}</Label>

            <InputGroup>
                <Input {...props} {...fields} type={values.showPassword ? 'text' : 'password'} defaultValue={values.password}
                    onChange={handleChange('password')} invalid={Boolean(touched[fields.name] && errors[fields.name])} />

                {touched[fields.name] && errors[fields.name] ? <FormFeedback className="order-last">{errors[fields.name]}</FormFeedback> : null}

                <InputGroupAddon addonType="append">
                    <button className="btn btn-secondary" type="button" onClick={handleClickShowPassword} onMouseDown={handleMouseDownPassword}>
                        {values.showPassword ? <i className="mdi mdi-eye"></i> : <i className="mdi mdi-eye-off"></i>}
                    </button>
                </InputGroupAddon>
            </InputGroup>
            
            {helpText && <FormText color="muted">{helpText}</FormText>}
        </FormGroup>
    );
}

const ReactstrapRootPasswordInput = ({
    field: { ...fields },
    form: { touched, errors, setFieldTouched, ...rest },
    helpText = false,
    ...props
}) => {
    const [values, setValues] = React.useState({
        password: fields.value ? fields.value : '',
        modal: false,
        showPassword: false,
    });
    
    const handleChange = prop => event => {
        setValues({ ...values, [prop]: event.target.value });
        if (props.onChange) props.onChange(event.target.valu);
        if (fields.onChange) fields.onChange(event);
    };

    const handleClickShowPassword = async () => {
        if(values.showPassword){
            setValues({ ...values, showPassword: !values.showPassword });
        }else{
            let result = await confirm({
                title: (
                    <>
                        <strong>{i18n.t(`${packageNS}:tr000620`)}</strong>!
                    </>
                ),
                message: (<><Row>{i18n.t(`${packageNS}:tr000621`)}</Row><Row style={{marginTop:10}}>{i18n.t(`${packageNS}:tr000622`)}</Row></>),
                confirmText: i18n.t(`${packageNS}:tr000623`),
                confirmColor: "danger",
                cancelColor: "link text-primary"
            });
            
            const id = props["data-id"];
            const sn = props["data-sn"].split("_")[1];
            const rootpassword = await GatewayStore.getRootConfig(id, sn);
            values.password = rootpassword.password;
            
            if(result){
                setValues({ ...values, password: values.password , showPassword: !values.showPassword });
            }else{
                return;
            }
        }
    };

    const handleMouseDownPassword = event => {
        event.preventDefault();
    };
    
    return (
        <FormGroup>
            <Label for={props.id}>{props.label}</Label>
            <InputGroup>
                <Input {...props} {...fields} type={values.showPassword ? 'text' : 'password'} value={values.password} defaultValue={'11111111'}
                    onChange={handleChange('password')} invalid={Boolean(touched[fields.name] && errors[fields.name])} readonly/>

                {touched[fields.name] && errors[fields.name] ? <FormFeedback className="order-last">{errors[fields.name]}</FormFeedback> : null}

                <InputGroupAddon addonType="append">
                    <button className="btn btn-secondary" type="button" onClick={handleClickShowPassword} onMouseDown={handleMouseDownPassword}>
                        {values.showPassword ? <i className="mdi mdi-eye"></i> : <i className="mdi mdi-eye-off"></i>}
                    </button>
                </InputGroupAddon>
            </InputGroup>
            
            {helpText && <FormText color="muted">{helpText}</FormText>}
        </FormGroup>
    );
}



export { ReactstrapInput, ReactstrapCheckbox, ReactstrapRadio, ReactstrapSelect, ReactstrapInputGroup, AsyncAutoComplete, ReactstrapPasswordInput, ReactstrapRootPasswordInput };
