import React, { Component } from 'react'
import QrReader from 'react-qr-reader'
 
class QRCodeReader extends Component {
  state = {
    result: 'No result'
  }
 
  handleScan = data => {
    if (data) {
      this.setState({
        result: data
      })
      this.props.toggle(data);
    }
  }
  handleError = err => {
    console.error(err)
  }
  render() {
    return (
      <div style={{width:500}}>
        <QrReader
          delay={300}
          onError={this.handleError}
          onScan={this.handleScan}
          style={{ width: '100%' }}
        />
        
      </div>
    )
  }
}

export default QRCodeReader;
