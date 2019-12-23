import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import WithdrawStore from "../../stores/WithdrawStore";

export default function AlertDialog(props) {
  const [open, setOpen] = React.useState(true);

  function handleClose() {
    setOpen(false);
  }

  const agree = () => {
    const data = props;

    /* WithdrawStore.update(data, resp => {
        props.history.push(`/withdraw/${props.match.params.organizationID}`);
    }); */

    handleClose();
  }

  return (
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Use Google's location service?"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Let Google help apps determine location. This means sending anonymous location data to
            Google, even when no apps are running.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={agree} color="primary" autoFocus>
            Agree
          </Button>
        </DialogActions>
      </Dialog>
  );
}
