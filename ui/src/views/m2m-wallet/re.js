import React, { Component } from 'react';
import Chart from './components/overview/SampleChart.jsx';
import Head from './components/overview/Head'
import StatsTable from './components/overview/StatsTable.jsx';
import './components/overview/overview.css';
import Sidebar from './components/overview/Sidebar'
import coin from './coindata.json';

class App extends Component {

  constructor() {
    super();

    this.state = {
        coin
    }

    setInterval(this.updateCoinData.bind(this), 1000);
  }

  updateCoinData() {

  }

  formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }

  render() {
    
    const { coin } = this.state;

    //console.log(coin)

    Object.keys(coin).forEach(attr => {
    const value = coin[attr];

    if (typeof value === 'number') {
        coin[attr] = this.formatNumber(value);
    }
    });

    return (
      <div className="app">
        <div className="overview-container">
          <div className="row">
            {/* Coin data */}
            <div className="col s12 m3 no-pading-right">
              <Sidebar />
            </div>
            <div className="col s12 m9 no-pading-left">
              <Head coin={coin}/>
              <StatsTable coin={coin} />
              <Chart />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

// export default withHighcharts(App, Highcharts);
export default App;