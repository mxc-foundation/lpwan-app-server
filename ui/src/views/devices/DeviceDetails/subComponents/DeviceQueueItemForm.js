import FormHelperText from "@material-ui/core/FormHelperText";
import { withStyles } from "@material-ui/core/styles";
import classnames from 'classnames';
import "codemirror/mode/javascript/javascript";
import { Field, Form, Formik } from 'formik';
import React from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { Button, Card, Nav, NavItem, NavLink, TabContent, TabPane } from 'reactstrap';
import * as Yup from 'yup';
import FormComponent from "../../../../classes/FormComponent";
import { ReactstrapCheckbox, ReactstrapInput } from '../../../../components/FormInputs';
import i18n, { packageNS } from '../../../../i18n';





const styles = {
  codeMirror: {
    border: "1px solid #ebeff2",
    height: "90%",
    marginRight: "0px",
    zIndex: 1,
  },
};

class DeviceQueueItemForm extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {
      object: {},
      // Important: We need to have tab 2 selected by default, otherwise
      // the code in the CodeMirror code snippet doesn't appear until you
      // click the code snippet input field.
      // See https://github.com/scniro/react-codemirror2/issues/83
      activeTab: "2",
    };
  }

  setActiveTab = (tab) => {
    this.setState({
      activeTab: tab
    })
  }

  toggle = (tab) => {
    const { activeTab } = this.state;
  
    if (activeTab !== tab) {
      this.setActiveTab(tab);
    }
  }

  onCodeChange = (field, editor, data, newCode) => {
    let object = this.state.object;
    object[field] = newCode;
    this.setState({
      object: object,
    });
  }

  render() {
    const { activeTab } = this.state;

    if (this.state.object === undefined) {
      return null;
    }

    const codeMirrorOptions = {
      lineNumbers: true,
      mode: "javascript",
      theme: "material"
    };

    let objectCode = this.state.object.jsonObject;
    if (objectCode === "" || objectCode === undefined) {
      objectCode = `{}`
    }

    let fieldsSchema = {
      fPort: Yup.number()
        .required(i18n.t(`${packageNS}:tr000431`)),
    }

    if (activeTab === "1") {
      fieldsSchema['data'] = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={(values) => {
            const castValues = formSchema.cast(values);
            this.props.onSubmit({ ...castValues })
          }}>
          {({
            handleSubmit,
            setFieldValue,
            values,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="number"
                  label={i18n.t(`${packageNS}:tr000285`)}
                  name="fPort"
                  id="fPort"
                  helpText={i18n.t(`${packageNS}:tr000286`)}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                />

                <Field
                  type="checkbox"
                  label={i18n.t(`${packageNS}:tr000287`)}
                  name="confirmed"
                  id="confirmed"
                  component={ReactstrapCheckbox}
                  onBlur={handleBlur}
                />

                <Card body style={{ backgroundColor: "#ebeff2" }}>
                  <Nav tabs>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '1' })}
                        onClick={() => { this.toggle('1'); }}
                      >
                        {i18n.t(`${packageNS}:tr000288`)}
                      </NavLink>
                    </NavItem>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '2' })}
                        onClick={() => { this.toggle('2'); }}
                      >
                        {i18n.t(`${packageNS}:tr000290`)}
                      </NavLink>
                    </NavItem>
                  </Nav>

                  <TabContent
                    activeTab={activeTab}
                    style={{
                      backgroundColor: "#fff",
                      borderRadius: "2px",
                      borderStyle: "solid",
                      borderWidth: "0 1px 1px 1px",
                      borderColor: "#ddd"
                    }}>
                    <TabPane tabId="1">
                      <Field
                        type="text"
                        label={i18n.t(`${packageNS}:tr000289`)}
                        name="data"
                        id="data"
                        component={ReactstrapInput}
                        onBlur={handleBlur}
                      />
                    </TabPane>
                    <TabPane tabId="2">
                      <h5>JSON Payload</h5>
                      <CodeMirror
                        autoScroll
                        cursor={{
                          line: 1,
                          ch: 2
                        }}
                        value={objectCode}
                        className={this.props.classes.codeMirror}
                        options={codeMirrorOptions}
                        onBeforeChange={this.onCodeChange.bind(this, 'jsonObject')}
                      />
                      <FormHelperText>
                        {i18n.t(`${packageNS}:tr000291`)}
                      </FormHelperText>
                    </TabPane>
                  </TabContent>
                </Card>
            
                <Button type="submit" color="primary">
                  {this.props.submitLabel || i18n.t(`${packageNS}:tr000292`)}
                </Button>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(DeviceQueueItemForm);
