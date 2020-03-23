import theme from "../../../theme";

const dashboardStyle = {
    root:{
        color:'#ffffff',
        
       

    },
    TextField:{
      '& input':{
         color:'#FFFFFF'
      }
     
  },
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
        flexDirection: 'column'
      },
      card: {
      
        width: '100%',
        backgroundColor: "#0C027060",
        color:"#ffffff",
      },
      cardTable:{
          '& td':{
           
            borderBottom:'none',
            '& span':{
                color:'#00FFD9',
                fontSize:'18px',
                fontWeight:'bold'
            }
          }
        
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
  
export default dashboardStyle;
