import theme from "../../theme";

const WithdrawStyles = {
    card: {
      minWidth: 180,
      width: 220,
      backgroundColor: "#0C0270",
    },
    title: {
      color: '#FFFFFF',
      fontSize: 14,
      padding: 6,
    },
    balance: {
      fontSize: 24,
      color: '#FFFFFF',
      textAlign: 'center',
    },
    newBalance: {
      fontSize: 24,
      textAlign: 'center',
      color: theme.palette.primary.main,
    },
    pos: {
      marginBottom: 12,
      color: '#FFFFFF',
      textAlign: 'right',
    },
    between: {
      display: 'flex',
      justifyContent:'spaceBetween'
    },
    flex: {
      display: 'flex',
      flexDirection: 'column'
    },
    navText: {
      fontSize: 14,
    },
    TitleBar: {
      height: 115,
      width: '50%',
      light: true,
      display: 'flex',
      flexDirection: 'column'
    },
    divider: {
      padding: 0,
      color: '#FFFFFF',
      width: '100%',
    },
    padding: {
      padding: 0,
    },
    link: {
      textDecoration: "none",
      fontWeight: "bold",
      fontSize: 12,
      color: theme.palette.textSecondary.main,
      opacity: 0.7,
        "&:hover": {
          opacity: 1,
        }
    },
  };
  
export default WithdrawStyles;
