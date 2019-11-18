<div className={this.props.classes.root}>
        <AppBar position="static" className={this.props.classes.appBar}>
          <Toolbar>
            <div className={this.props.logoSection}>
              <img src="/logo/logo_mx.png" className={this.props.classes.logo} alt="LPWAN Server" />
            </div>
            <IconButton edge="start" className={this.props.classes.menuButton} color="inherit" aria-label="menu">
              {/* <MenuIcon /> */}
            </IconButton>
            <Typography variant="h6" className={this.props.classes.title}></Typography>
            <Button variant="outlined"
                    color="inherit"
                    onClick={this.onClick}
            >ACCESS</Button>
          </Toolbar>
        </AppBar>
      </div>