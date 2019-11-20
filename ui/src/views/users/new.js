import i18n, { packageNS } from '../../i18n';

<div className={this.props.classes.root}>
        <AppBar position="static" className={this.props.classes.appBar}>
          <Toolbar>
            <div className={this.props.logoSection}>
              <img src="/logo/logo_mx.png" className={this.props.classes.logo} alt={i18n.t(`${packageNS}:tr000051`)} />
            </div>
            <IconButton edge="start" className={this.props.classes.menuButton} color="inherit" aria-label="menu">
              {/* <MenuIcon /> */}
            </IconButton>
            <Typography variant="h6" className={this.props.classes.title}></Typography>
            <Button variant="outlined"
                    color="inherit"
                    onClick={this.onClick}
            >{i18n.t(`${packageNS}:tr000002`)}</Button>
          </Toolbar>
        </AppBar>
      </div>