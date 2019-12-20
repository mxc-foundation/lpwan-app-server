import theme from "../../theme";

const TopupStyle = {
    tabs: {
        borderBottom: "1px solid " + theme.palette.divider,
        height: "49px",
      },
      navText: {
        fontSize: 14,
      },
      TitleBar: {
        height: 115,
        width: '50%',
        light: true,
        display: 'flex',
        flexDirection: 'column',
        padding: '0px 0px 50px 0px' 
      },
      card: {
        minWidth: 180,
        width: '100%',
        backgroundColor: 'white',
      },
      divider: {
        padding: 0,
        color: '#FFFFFF',
        width: '100%',
      },
      padding: {
        padding: 0,
      },
      column: {
        display: 'flex',
        flexDirection: 'column',
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
  
export default TopupStyle;
