import theme from "../../../theme";

const HistoryStyle = {
    tabs: {
        borderBottom: "1px solid " + theme.palette.divider,
        height: "49px",
        
      },
      tabsBlock:{
        marginTop:'25px',
        marginBottom:'15px'
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
        width: 250,
        backgroundColor: "#0C027060",
        color:'#ffffff',
        padding:'15px'
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
  
export default HistoryStyle;
