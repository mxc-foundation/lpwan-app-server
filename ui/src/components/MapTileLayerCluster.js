import React, { Component } from 'react';

import GatewayStore from '../stores/GatewayStore';
import Paper from '../components/Paper';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import styled from 'styled-components';
import 'leaflet.markercluster/dist/MarkerCluster.css'
import 'leaflet.markercluster/dist/MarkerCluster.Default.css'
import 'leaflet.markercluster/dist/leaflet.markercluster.js'
import { withRouter } from "react-router-dom";

const Wrapper = styled.div`
    width: ${props => props.width};
    height: ${props => props.height};
`
function loadGatewayData() {
    return new Promise((resolve, reject) => {
        GatewayStore.listLocations(
        resp => {
          resolve(resp);
        })
    });
} 
class Context extends Component {
  render() {
    return(
      <Paper>
        dfsdf
      </Paper>
    );
  }
}
class MapTileLayerCluster extends Component {
    componentDidMount(){
        this.map = L.map('map', {
            //preferCanvas: true,
            center: [51.1657, 10.4515],
            zoom: 6,
            zoomContorl: false
        });

        L.tileLayer('https://cartodb-basemaps-{s}.global.ssl.fastly.net/rastertiles/voyager/{z}/{x}/{y}{r}.png',{
            detectRetina: true,
            maxZoom: 20,
            maxNativeZoom: 17
        }).addTo(this.map);
        
        let greenIcon = L.icon({
          iconUrl: 'data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4KPCEtLSBHZW5lcmF0b3I6IEFkb2JlIElsbHVzdHJhdG9yIDIzLjAuNCwgU1ZHIEV4cG9ydCBQbHVnLUluIC4gU1ZHIFZlcnNpb246IDYuMDAgQnVpbGQgMCkgIC0tPgo8c3ZnIHZlcnNpb249IjEuMSIgaWQ9IkNhcGFfMSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgeD0iMHB4IiB5PSIwcHgiCgkgdmlld0JveD0iMCAwIDQ2NSA2OTYuMSIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNDY1IDY5Ni4xOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+CjxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+Cgkuc3Qwe29wYWNpdHk6MC4yNTt9Cjwvc3R5bGU+CjxnPgoJPHBhdGggZD0iTTMxNCw2MzEuMnYtMzcwYzAtMjYuMi0yMS4zLTQ3LjUtNDcuNS00Ny41aC00OC43Yy0yNi4yLDAtNDcuNSwyMS4zLTQ3LjUsNDcuNXYzNzBjMCwyNi4yLDIxLjMsNDcuNSw0Ny41LDQ3LjVoNDguNwoJCUMyOTIuNiw2NzguNywzMTQsNjU3LjQsMzE0LDYzMS4yeiBNMjE3LjgsNjYzLjdjLTE3LjksMC0zMi41LTE0LjYtMzIuNS0zMi41di0zNzBjMC0xNy45LDE0LjYtMzIuNSwzMi41LTMyLjVoNDguNwoJCWMxNy45LDAsMzIuNSwxNC42LDMyLjUsMzIuNXYzNzBjMCwxNy45LTE0LjYsMzIuNS0zMi41LDMyLjVIMjE3Ljh6Ii8+Cgk8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMjQxLjQsNzYuOGMxLjYsMCwzLjEtMC42LDQuMy0xLjhjMi40LTIuNCwyLjQtNi4yLDAtOC42Yy0yNS43LTI1LjctNjcuNC0yNS43LTkzLjEsMGMtMi40LDIuNC0yLjQsNi4yLDAsOC42CgkJYzIuNCwyLjQsNi4yLDIuNCw4LjYsMGMyMC45LTIwLjksNTUtMjAuOSw3NS45LDBDMjM4LjMsNzYuMiwyMzkuOSw3Ni44LDI0MS40LDc2Ljh6Ii8+Cgk8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMTQwLjEsNTMuOWMzMi42LTMyLjYsODUuNi0zMi42LDExOC4yLDBjMS4yLDEuMiwyLjcsMS44LDQuMywxLjhjMS42LDAsMy4xLTAuNiw0LjMtMS44CgkJYzIuNC0yLjQsMi40LTYuMiwwLTguNmMtMTgtMTgtNDIuMS0yOC02Ny43LTI4cy00OS42LDkuOS02Ny43LDI4Yy0yLjQsMi40LTIuNCw2LjIsMCw4LjZDMTMzLjksNTYuMywxMzcuOCw1Ni4zLDE0MC4xLDUzLjl6Ii8+Cgk8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMjI0LjYsOTYuMmMyLjQtMi40LDIuNC02LjIsMC04LjZjLTE0LTE0LTM2LjgtMTQtNTAuOCwwYy0yLjQsMi40LTIuNCw2LjIsMCw4LjZjMi40LDIuNCw2LjIsMi40LDguNiwwCgkJYzkuMy05LjMsMjQuNC05LjMsMzMuNywwYzEuMiwxLjIsMi43LDEuOCw0LjMsMS44UzIyMy40LDk3LjQsMjI0LjYsOTYuMnoiLz4KCTxwYXRoIGQ9Ik0yNjQuMiw1NTguN2MwLTQuMS0zLjQtNy41LTcuNS03LjVoLTI3LjNjLTQuMSwwLTcuNSwzLjQtNy41LDcuNXMzLjQsNy41LDcuNSw3LjVoMjcuMwoJCUMyNjAuOCw1NjYuMiwyNjQuMiw1NjIuOCwyNjQuMiw1NTguN3oiLz4KCTxwYXRoIGQ9Ik0yMTEuMywyMTAuNWMzLjQsMCw2LjEtMi43LDYuMS02LjF2LTc5LjNjMC0xMC04LjItMTguMi0xOC4yLTE4LjJzLTE4LjIsOC4yLTE4LjIsMTguMnY3OS4zYzAsMy40LDIuNyw2LjEsNi4xLDYuMQoJCWMzLjQsMCw2LjEtMi43LDYuMS02LjF2LTc5LjNjMC0zLjMsMi43LTYuMSw2LjEtNi4xYzMuMywwLDYuMSwyLjcsNi4xLDYuMXY3OS4zQzIwNS4zLDIwNy44LDIwOCwyMTAuNSwyMTEuMywyMTAuNXoiLz4KCTxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik0zMjEuOSw4OC4zYzEuNCwwLDIuNy0wLjUsMy44LTEuNmMyLjEtMi4xLDIuMS01LjUsMC03LjZjLTIyLjYtMjIuNi01OS40LTIyLjYtODIsMGMtMi4xLDIuMS0yLjEsNS41LDAsNy42CgkJYzIuMSwyLjEsNS41LDIuMSw3LjYsMGMxOC40LTE4LjQsNDguNC0xOC40LDY2LjksMEMzMTkuMiw4Ny44LDMyMC41LDg4LjMsMzIxLjksODguM3oiLz4KCTxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik0zMDcuMSwxMDUuM2MyLjEtMi4xLDIuMS01LjUsMC03LjZjLTEyLjMtMTIuMy0zMi40LTEyLjMtNDQuOCwwYy0yLjEsMi4xLTIuMSw1LjUsMCw3LjYKCQljMi4xLDIuMSw1LjUsMi4xLDcuNiwwYzguMi04LjIsMjEuNS04LjIsMjkuNywwYzEsMSwyLjQsMS42LDMuOCwxLjZTMzA2LDEwNi40LDMwNy4xLDEwNS4zeiIvPgoJPHBhdGggZD0iTTI5Ny4xLDIxMC41YzMuNCwwLDYuMS0yLjcsNi4xLTYuMXYtNzkuM2MwLTEwLTguMi0xOC4yLTE4LjItMTguMmMtMTAsMC0xOC4yLDguMi0xOC4yLDE4LjJ2NzkuM2MwLDMuNCwyLjcsNi4xLDYuMSw2LjEKCQlzNi4xLTIuNyw2LjEtNi4xdi03OS4zYzAtMy4zLDIuNy02LjEsNi4xLTYuMWMzLjMsMCw2LjEsMi43LDYuMSw2LjF2NzkuM0MyOTEsMjA3LjgsMjkzLjcsMjEwLjUsMjk3LjEsMjEwLjV6Ii8+CjwvZz4KPC9zdmc+Cg==',
          //shadowUrl: '',
        
          iconSize:     [30, 56], // size of the icon
          shadowSize:   [50, 64], // size of the shadow
          iconAnchor:   [22, 94], // point of the icon which will correspond to marker's location
          shadowAnchor: [4, 62],  // the same for the shadow
          popupAnchor:  [-3, -76] // point from which the popup should open relative to the iconAnchor
        });

        var markers = L.markerClusterGroup({
          showCoverageOnHover: true,
        });
        
        const res = this.loadData();
        res.then(resp => {
            for(let i=0;i<resp.result.length; i++){
                markers.addLayer(L.marker([resp.result[i].location.latitude, resp.result[i].location.longitude],{icon: greenIcon}));
            }
        })
        
        this.map.addLayer(markers);
    }

    loadData = async () => {
        try {
            var result = await loadGatewayData();
            return result;
        } catch (error) {
            this.setState({loading: false})
            console.error(error);
            this.setState({ error });
        }
    }
    
  render() {
    return(
      <Wrapper
        width="100%" height='100%' id='map'
      />
    )
  }
}

export default withRouter(MapTileLayerCluster);
