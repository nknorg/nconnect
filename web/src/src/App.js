import React from 'react';
import QRCode from 'qrcode';
import { Button, Collapse, Container, TextField } from '@material-ui/core';
import { ExpandLess, ExpandMore } from '@material-ui/icons';

import * as rpc from './rpc';

import './App.css';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      adminTokenStr: '',
      adminTokenQRCode: '',
      acceptAddrs: '',
      adminAddrs: '',
      addr: '',
      localIP: [],
      showAdvanced: false,
    };
    this.handleAcceptAddrsChange = this.handleAcceptAddrsChange.bind(this);
    this.handleAdminAddrsChange = this.handleAdminAddrsChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.updateAdminToken = this.updateAdminToken.bind(this);
    this.handleAdvancedChange = this.handleAdvancedChange.bind(this);
  }

  handleAcceptAddrsChange(event) {
    this.setState({ acceptAddrs: event.target.value });
  }

  handleAdminAddrsChange(event) {
    this.setState({ adminAddrs: event.target.value });
  }

  handleAdvancedChange(event) {
    this.setState({ showAdvanced: !this.state.showAdvanced });
    console.log(this.state.showAdvanced);
  }

  async handleSubmit(event) {
    event.preventDefault();
    try {
      let addrs = await rpc.setAddrs(strToAddrs(this.state.acceptAddrs), strToAddrs(this.state.adminAddrs));
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
      alert('Save success!');
    } catch (e) {
      console.error(e);
      alert(e);
    }
  }

  async updateAdminToken() {
    try {
      let adminToken = await rpc.getAdminToken();
      if (adminToken) {
        let adminTokenStr = JSON.stringify(adminToken);
        let adminTokenQRCode = await QRCode.toDataURL(adminTokenStr);
        this.setState({
          adminTokenStr,
          adminTokenQRCode,
        });
      }
    } catch (e) {
      console.error(e);
    }
  }

  componentDidMount() {
    this.updateAdminToken();
    setInterval(this.updateAdminToken, 5 * 60 * 1000);

    rpc.getAddrs().then((addrs) => {
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
    }).catch((e) => {
      console.error(e);
      alert(e);
    });

    rpc.getInfo().then((info) => {
      this.setState({
        addr: info.addr,
        localIP: info.localIP.ipv4,
      });
    }).catch((e) => {
      console.error(e);
    });
  }

  render() {
    return (
      <div className="App">
        <Container>
          <div className="row">
            <img src="/static/media/nkn_logo.png" />
          </div>
          <div className="row">
            <img src={this.state.adminTokenQRCode} />
          </div>
          <div className="row">
            Scan the QR code on nMobile Pro to connect and manage device.
          </div>
          <div className="row">
            Purchase data plan on nMobile Pro or <a target="_blank" href={"https://nconnect-payment.nkncdn.com/payment/?addr=" + addrToPubKey(this.state.addr)}>web payment portal</a>
          </div>
          <div className="row">
            <Button
              variant="outlined"
              color="primary"
              onClick={this.handleAdvancedChange}
              style={{width: '100%'}}
              >
              {this.state.showAdvanced ? <ExpandLess /> : <ExpandMore /> }
              {this.state.showAdvanced ? "Hide Advanced" : "Show Advanced"}
            </Button>
          </div>
          <Collapse in={this.state.showAdvanced}>
            <div className="advanced-row">
              <TextField
                disabled
                multiline
                label="Local IP address"
                value={this.state.localIP.join('\n')}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <TextField
                disabled
                multiline
                label="Access key (expires in 5 minutes)"
                value={this.state.adminTokenStr}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <TextField
                multiline
                variant="filled"
                label="Accept addresses"
                value={this.state.acceptAddrs}
                onChange={this.handleAcceptAddrsChange}
                style={{width: '100%'}}
                />
              <TextField
                multiline
                variant="filled"
                label="Admins"
                value={this.state.adminAddrs}
                onChange={this.handleAdminAddrsChange}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <Button
                variant="contained"
                color="primary"
                onClick={this.handleSubmit}
                style={{width: '100%'}}
                >
                Save
              </Button>
            </div>
          </Collapse>
        </Container>
      </div>
    );
  }
}

function addrsToStr(addrs) {
  if (!addrs) {
    return '';
  }
  return addrs.join('\n');
}

function strToAddrs(str) {
  if (!str) {
    return [];
  }
  return str.split('\n').filter(s => s.length > 0);
}

function addrToPubKey(addr) {
  let s = addr.split('.');
  return s[s.length-1];
}

export default App;
