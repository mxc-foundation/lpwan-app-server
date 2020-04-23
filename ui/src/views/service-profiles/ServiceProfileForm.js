import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Button, Col, Row } from 'reactstrap';
import * as Yup from 'yup';
import Admin from '../../components/Admin';
import { AsyncAutoComplete, ReactstrapCheckbox, ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";




class ServiceProfileForm extends Component {
  constructor(props) {
    super(props);

    this.state = {};

    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
  }

  componentDidMount() {
    this.setState({
      ...this.props.object,
    });
  }

  getNetworkServerOptions = async (search, callbackFunc) => {
    const res = await NetworkServerStore.list(0, 10, 0);
    const options = res.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
    callbackFunc(options);
  }


  render() {
    const object = this.state;

    if (object === undefined) {
      return (<div></div>);
    }

    let fieldsSchema = {
      name: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
      networkServerID: Yup.string(),
      id: Yup.string(),
      addGWMetaData: Yup.bool(),
      nwkGeoLoc: Yup.bool(),
      devStatusReqFreq: Yup.number().moreThan(-1, i18n.t(`${packageNS}:menu.messages.min`)),
      drMin: Yup.number().moreThan(-1, i18n.t(`${packageNS}:menu.messages.min`)),
      drMax: Yup.number().moreThan(-1, i18n.t(`${packageNS}:menu.messages.min`))
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
            enableReinitialize
            initialValues={{
              name: object.name || '',
              networkServerID: object.networkServerID || '',
              id: object.id,
              addGWMetaData: object.addGWMetaData || false,
              nwkGeoLoc: object.nwkGeoLoc || false,
              devStatusReqFreq: object.devStatusReqFreq || '',

              /* reportDevStatusBattery: object.reportDevStatusBattery,
              reportDevStatusMargin: object.reportDevStatusMargin, */

              drMin: object.drMin || '',
              drMax: object.drMax || ''
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
                    type="text"
                    label={i18n.t(`${packageNS}:tr000149`) + "*"}
                    name="name"
                    id="name"
                    value={values.name}
                    onChange={handleChange}
                    helpText={i18n.t(`${packageNS}:tr000150`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                    required
                  />

                  {!this.props.update && <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000047`) + "*"}
                    name="networkServerID"
                    id="networkServerID"
                    getOptions={this.getNetworkServerOptions}
                    setFieldValue={setFieldValue}
                    helpText={i18n.t(`${packageNS}:tr000223`)}
                    onBlur={handleBlur}
                    inputProps={{
                      clearable: true,
                      cache: false,
                    }}
                    component={AsyncAutoComplete}
                    required
                  />}

                  <Field
                    type="checkbox"
                    name="addGWMetaData"
                    id="addGWMetaData"
                    label={i18n.t(`${packageNS}:tr000151`)}
                    component={ReactstrapCheckbox}
                    onChange={handleChange}
                    helpText={i18n.t(`${packageNS}:tr000152`)}
                    onBlur={handleBlur}
                  />

                  <Field
                    type="checkbox"
                    name="nwkGeoLoc"
                    id="nwkGeoLoc"
                    label={i18n.t(`${packageNS}:tr000153`)}
                    component={ReactstrapCheckbox}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    helpText={i18n.t(`${packageNS}:tr000154`)}
                  />

                  <Field
                    type="number"
                    label={i18n.t(`${packageNS}:tr000155`)}
                    name="devStatusReqFreq"
                    id="devStatusReqFreq"
                    helpText={i18n.t(`${packageNS}:tr000156`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                  />

                  {/* <FormGroup>
                        <Field
                            type="checkbox"
                            label={i18n.t(`${packageNS}:tr000157`)}
                            name="reportDevStatusBattery"
                            id="reportDevStatusBattery"
                            component={ReactstrapCheckbox}
                            onChange={handleChange}
                        />

                        <Field
                            type="checkbox"
                            label={i18n.t(`${packageNS}:tr000158`)}
                            name="reportDevStatusMargin"
                            id="reportDevStatusMargin"
                            component={ReactstrapCheckbox}
                            onChange={handleChange}
                        />
                    </FormGroup> */}

                  <Field
                    type="number"
                    label={i18n.t(`${packageNS}:tr000159`) + "*"}
                    name="drMin"
                    id="drMin"
                    helpText={i18n.t(`${packageNS}:tr000160`)}
                    component={ReactstrapInput}
                    required
                  />

                  <Field
                    type="number"
                    label={i18n.t(`${packageNS}:tr000161`) + "*"}
                    name="drMax"
                    id="drMax"
                    helpText={i18n.t(`${packageNS}:tr000162`)}
                    component={ReactstrapInput}
                    required
                  />
                  <Admin>
                    <Button type="submit" color="primary">{this.props.submitLabel}</Button>
                  </Admin>
                </Form>
              )}
          </Formik>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ServiceProfileForm;
