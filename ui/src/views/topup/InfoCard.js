import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';

import i18n, { packageNS } from '../../i18n';
import TitleBarTitle from "../../components/TitleBarTitle";
import { Link  } from "react-router-dom";

import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles({
  card: {
    maxWidth: '100%',
    //margin: '16px 0px 16px 0px',
    margin: 16,
    backgroundColor: 'white',
  },
  media: {
    height: 140,
  },
});

export default function MediaCard(props) {
  const classes = useStyles();

  return (
    <Card className={classes.card}>
      <CardActionArea>
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            {i18n.t(`${packageNS}:menu.topup.synchronize_your_eth_account`)}
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
            {i18n.t(`${packageNS}:menu.topup.note`)}
          </Typography>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <TitleBarTitle component={Link} to={`${props.path}`} title={i18n.t(`${packageNS}:menu.topup.change_eth_account`)} />
      </CardActions>
    </Card>
  );
}