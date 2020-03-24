import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import i18n, { packageNS } from '../../i18n';

export default function ConfirmDialog(props) {
  const { open, onClose, title, description } = props

  const agree = () => {
    const { data, onConfirm } = props;

    onConfirm(data);

    if (onClose) onClose();
  }

  return (
      <Dialog
        open={open}
        onClose={onClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{title}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {description}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} color="primary25" autoFocus>
            {i18n.t(`${packageNS}:menu.staking.cancel`)}
          </Button>
          <Button onClick={agree} color="primary25">
            {i18n.t(`${packageNS}:menu.settings.proceed`)}
          </Button>
        </DialogActions>
      </Dialog>
  );
}
