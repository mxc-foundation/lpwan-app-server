import { withStyles } from "@material-ui/core/styles";
import Typography from '@material-ui/core/Typography';
import "codemirror/mode/javascript/javascript";
import { Field, Form, Formik } from 'formik';
import React from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { Button, FormText, Label } from 'reactstrap';
import * as Yup from 'yup';
import FormComponent from "../../classes/FormComponent";
import { AsyncAutoComplete, ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import ServiceProfileStore from "../../stores/ServiceProfileStore";




const styles = {
  codeMirror: {
    zIndex: 1,
  }
};

class ApplicationForm extends FormComponent {
  getServiceProfileOption = (id, callbackFunc) => {
    ServiceProfileStore.get(id, resp => {
      callbackFunc({label: resp.serviceProfile.name, value: resp.serviceProfile.id});
    });
  }

  getServiceProfileOptions = (search, callbackFunc) => {
    ServiceProfileStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((sp, i) => {return {label: sp.name, value: sp.id}});
      callbackFunc(options);
    });
  }

  getPayloadCodecOptions = (search, callbackFunc) => {
    const payloadCodecOptions = [
      {value: "", label: i18n.t(`${packageNS}:tr000211`)},
      {value: "CAYENNE_LPP", label: i18n.t(`${packageNS}:tr000212`)},
      {value: "CUSTOM_JS", label: i18n.t(`${packageNS}:tr000212`)},
    ];

    callbackFunc(payloadCodecOptions);
  }

  onCodeChange = (field, editor, data, newCode) => {
    let object = this.state.object;
    object[field] = newCode;
    this.setState({
      object: object,
    });
  }

  render() {
    const { submitLabel } = this.props;

    if (this.state.object === undefined) {
      return(<div></div>);
    }

    const codeMirrorOptions = {
      lineNumbers: true,
      mode: "javascript",
      theme: "default",
    };
    
    let payloadEncoderScript = this.state.object.payloadEncoderScript;
    let payloadDecoderScript = this.state.object.payloadDecoderScript;

    if (payloadEncoderScript === "" || payloadEncoderScript === undefined) {
      payloadEncoderScript = `// Encode encodes the given object into an array of bytes.
//  - fPort contains the LoRaWAN fPort number
//  - obj is an object, e.g. {"temperature": 22.5}
// The function must return an array of bytes, e.g. [225, 230, 255, 0]
function Encode(fPort, obj) {
  return [];
}`;
    }

    if (payloadDecoderScript === "" || payloadDecoderScript === undefined) {
      payloadDecoderScript = `// Decode decodes an array of bytes into an object.
//  - fPort contains the LoRaWAN fPort number
//  - bytes is an array of bytes, e.g. [225, 230, 255, 0]
// The function must return an object, e.g. {"temperature": 22.5}
function Decode(fPort, bytes) {
  return {};
}`;
    }

    let fieldsSchema = {
      name: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
      //.matches(/^[a-zA-Z0-9]+$/, i18n.t(`${packageNS}:tr000429`)).required(i18n.t(`${packageNS}:tr000431`)),
      description: Yup.string().trim()
        .required(i18n.t(`${packageNS}:tr000431`)),
      serviceProfileID: Yup.string().trim()
        .required(i18n.t(`${packageNS}:tr000431`)),
    }
    
    if (!this.props.update) {
      fieldsSchema['name'] = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema['serviceProfileID'] = Yup.string().trim();
    }
    const formSchema = Yup.object().shape(fieldsSchema);

    return(
      <React.Fragment>
        <Formik
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={(values) => {
            const castValues = formSchema.cast(values);
            this.props.onSubmit({ ...castValues })
          }}>
          {
            ({
              handleSubmit,
              setFieldValue,
              values
            }) => (
            <Form onSubmit={handleSubmit} noValidate>
              <Field
                type="text"
                label={i18n.t(`${packageNS}:tr000254`)}
                name="name"
                id="name"
                helpText={i18n.t(`${packageNS}:tr000062`)}
                component={ReactstrapInput}
              />

              <Field
                type="text"
                label={i18n.t(`${packageNS}:tr000255`)}
                name="description"
                id="description"
                component={ReactstrapInput}
              />

              {!this.props.update &&
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:tr000078`)}
                  name="serviceProfileID"
                  id="serviceProfileID"
                  value={this.state.object.serviceProfileID || ""}
                  getOption={this.getServiceProfileOption}
                  getOptions={this.getServiceProfileOptions}
                  setFieldValue={setFieldValue}
                  helpText={i18n.t(`${packageNS}:tr000257`)}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                  component={AsyncAutoComplete}
                />
              }

              {this.state.object.payloadCodec !== "" &&
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:tr000209`)}
                  name="payloadCodec"
                  id="payloadCodec"
                  value={this.state.object.payloadCodec || ""}
                  getOption={this.getPayloadCodecOption}
                  getOptions={this.getPayloadCodecOptions}
                  setFieldValue={setFieldValue}
                  helpText={i18n.t(`${packageNS}:tr000258`)}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                  component={AsyncAutoComplete}
                />
              }

              {this.state.object.payloadCodec === "CUSTOM_JS" &&
                <>
                  <Label for="payloadEncoderScript">
                    {i18n.t(`${packageNS}:tr000551`)}
                  </Label>
                  <CodeMirror
                    value={payloadDecoderScript}
                    options={codeMirrorOptions}
                    onBeforeChange={this.onCodeChange.bind(this, 'payloadDecoderScript')}
                    className={this.props.classes.codeMirror}
                  />
                  <FormText>
                    {i18n.t(`${packageNS}:tr000215`)}
                  </FormText>
                  <br />
                </>
              }

              {this.state.object.payloadCodec === "CUSTOM_JS" &&
                <>
                  <Label for="payloadEncoderScript">
                    {i18n.t(`${packageNS}:tr000552`)}
                  </Label>
                  <CodeMirror
                    value={payloadEncoderScript}
                    options={codeMirrorOptions}
                    onBeforeChange={this.onCodeChange.bind(this, 'payloadEncoderScript')}
                    className={this.props.classes.codeMirror}
                  />
                  <FormText>
                    {i18n.t(`${packageNS}:tr000216`)}
                  </FormText>
                </>
              }

              {this.state.object.payloadCodec === "" &&
                <Typography variant="body1">
                  <br />
                  {i18n.t(`${packageNS}:tr000259`)}
                </Typography>
              }

              <br />
              <Button
                aria-label={submitLabel}
                block
                color="primary"
                size="md"
                type="submit"
              >
                {submitLabel}
              </Button>
            </Form>
          )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ApplicationForm);
