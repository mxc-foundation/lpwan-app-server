import React from 'react';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';

export default function SwitchLabels(props) {
  const [state, setState] = React.useState({
    checked: props.on
  });

  React.useEffect(() => {
    setState({checked: props.on})
  }, [props.on])

  const handleChange = name => event => {
    setState({ ...state, [name]: event.target.checked });
    
    props.onSwitchChange(props.dvId, event.target.checked, event);
  };

  return (
    <FormGroup row>
      <FormControlLabel
        control={
          <Switch
            checked={state.checked}
            onChange={handleChange('checked')}
            value="checked"
            color="primary"
          />
        }
      />
    </FormGroup>
  );
}
