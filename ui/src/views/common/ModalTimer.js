import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import LinearDeterminate from '../../components/LinearDeterminate';
import i18n, { packageNS } from '../../i18n';

export default function ConfirmDialog(props) {
  const { open, onCancelProgress, title, description, onProgress } = props

  const agree = () => {
    const { data, onConfirm } = props;

    onConfirm(data);

    if (onCancelProgress) onCancelProgress();
  }

  return (
      <Dialog
        open={open}
        onCancelProgress={onCancelProgress}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{title}</DialogTitle>
        <DialogContent>
        <LinearDeterminate onProgress={onProgress}/>
          <DialogContentText id="alert-dialog-description">
            {description}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancelProgress} color="primary.main" autoFocus>
            {i18n.t(`${packageNS}:menu.withdraw.cancel`)}
          </Button>
        </DialogActions>
      </Dialog>
  );
}
