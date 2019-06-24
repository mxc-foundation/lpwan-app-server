import { createMuiTheme } from "@material-ui/core/styles";
import teal from "@material-ui/core/colors/teal";


const theme = createMuiTheme({
    palette: {
      primary: { main: teal['A200'] },
      secondary: { main: '#11cb5f' }, 
    },
});
  
export default theme;
